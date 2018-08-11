package azureMetricClient

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/selection"

	"github.com/Azure/azure-k8s-metrics-adapter/pkg/aiapiclient"
	"github.com/Azure/azure-k8s-metrics-adapter/pkg/aim"
	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/metrics/pkg/apis/external_metrics"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2018-03-01/insights"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

// AzureMetricClient is used to make requests to Azure Monitor
type AzureMetricClient struct {
	monitorClient         insights.MetricsClient
	appinsightsclient     aiapiclient.AiAPIClient
	defaultSubscriptionID string
}

// NewAzureMetricClient creates a client for making requests to Azure Monitor
func NewAzureMetricClient() AzureMetricClient {
	defaultSubscriptionID := getDefaultSubscriptionID()

	// looks for ENV variables then falls back to AIM issue #10
	monitorClient := insights.NewMetricsClient(defaultSubscriptionID)
	appInsightsClient := aiapiclient.NewAiAPIClient()
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		monitorClient.Authorizer = authorizer
	}

	return AzureMetricClient{
		monitorClient:         monitorClient,
		appinsightsclient:     appInsightsClient,
		defaultSubscriptionID: defaultSubscriptionID,
	}
}

// GetAzureMetric calls Azure Monitor endpoint and returns a metric based on label selectors
func (c AzureMetricClient) GetAzureMetric(metricSelector labels.Selector) (external_metrics.ExternalMetricValue, error) {
	azMetricRequest, err := parseAzureMetric(metricSelector, c.defaultSubscriptionID)
	if err != nil {
		return external_metrics.ExternalMetricValue{}, err
	}
	metricResourceURI := azMetricRequest.metricResourceURI()

	glog.V(2).Infof("resource uri: %s", metricResourceURI)
	glog.V(2).Infof("filter: %s", azMetricRequest.filter)
	glog.V(2).Infof("metric name : %s", azMetricRequest.metricName)

	// make call to azure resource provider with subscription id provided issue #9
	c.monitorClient.SubscriptionID = azMetricRequest.subscriptionID
	metricResult, err := c.monitorClient.List(context.Background(), metricResourceURI,
		azMetricRequest.timespan, nil,
		azMetricRequest.metricName, azMetricRequest.aggregation, nil,
		"", azMetricRequest.filter, "", "")
	if err != nil {
		return external_metrics.ExternalMetricValue{}, err
	}

	total := extractValue(metricResult)

	// TODO set Value based on aggregations type
	return external_metrics.ExternalMetricValue{
		MetricName: azMetricRequest.resourceName,
		Value:      *resource.NewQuantity(int64(total), resource.DecimalSI),
		Timestamp:  metav1.Now(),
	}, nil
}

func (c AzureMetricClient) GetCustomMetric(groupResource schema.GroupResource, namespace string, selector labels.Selector, metricName string) (int64, error) {

	result, err := c.appinsightsclient.GetMetric("performanceCounters/requestsPerSecond", "avg")
	if err != nil {
		return 0, err
	}

	return *result, nil
}

func extractValue(metricResult insights.Response) float64 {
	//TODO extract value based on aggregation type
	//TODO check for nils
	metricVals := *metricResult.Value
	Timeseries := *metricVals[0].Timeseries
	data := *Timeseries[0].Data
	total := *data[len(data)-1].Total

	return total
}

func getDefaultSubscriptionID() string {
	// if the user explicitly sets we should use that
	subscriptionID := os.Getenv("SUBSCRIPTION_ID")
	if subscriptionID == "" {
		//fallback to trying azure instance meta data
		azureConfig, err := aim.GetAzureConfig()
		if err != nil {
			glog.Errorf("Unable to get azure config from MSI: %v", err)
		}

		subscriptionID = azureConfig.SubscriptionID
	}

	if subscriptionID == "" {
		glog.V(0).Info("Default Azure Subscription is not set.  You must provide subscription id via HPA lables, set an environment variable, or enable MSI.  See docs for more details")
	}

	return subscriptionID
}

type azureMetricRequest struct {
	metricName                string
	resourceGroup             string
	resourceName              string
	resourceProviderNamespace string
	resourceType              string
	aggregation               string
	timespan                  string
	filter                    string
	subscriptionID            string
}

func parseAzureMetric(metricSelector labels.Selector, defaultSubscriptionID string) (azureMetricRequest, error) {
	glog.V(2).Infof("begin parsing metric")

	// Using selectors to pass required values thorugh
	// to retain camel case as azure provider is case sensitive.
	//
	// There is are restrictions so using some conversion
	// restrictions here
	// note: requirement values are already validated by apiserver
	merticReq := azureMetricRequest{
		timespan:       timeSpan(),
		subscriptionID: defaultSubscriptionID,
	}
	requirements, _ := metricSelector.Requirements()
	for _, request := range requirements {
		if request.Operator() != selection.Equals {
			return azureMetricRequest{}, errors.New("selector type not supported. only equals is supported at this time")
		}

		value := request.Values().List()[0]

		switch request.Key() {
		case "metricName":
			glog.V(2).Infof("metricName: %s", value)
			merticReq.metricName = value
		case "resourceGroup":
			glog.V(2).Infof("resourceGroup: %s", value)
			merticReq.resourceGroup = value
		case "resourceName":
			glog.V(2).Infof("resourceName: %s", value)
			merticReq.resourceName = value
		case "resourceProviderNamespace":
			glog.V(2).Infof("resourceProviderNamespace: %s", value)
			merticReq.resourceProviderNamespace = value
		case "resourceType":
			glog.V(2).Infof("resourceType: %s", value)
			merticReq.resourceType = value
		case "aggregation":
			glog.V(2).Infof("aggregation: %s", value)
			merticReq.aggregation = value
		case "filter":
			// TODO: Should handle filters by converting equality and setbased label selectors
			// to  oData syntax: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#label-selectors
			glog.V(2).Infof("filter: %s", value)
			filterStrings := strings.Split(value, "_")
			merticReq.filter = fmt.Sprintf("%s %s '%s'", filterStrings[0], filterStrings[1], filterStrings[2])
			glog.V(2).Infof("filter formatted: %s", merticReq.filter)
		case "subscriptionID":
			// if sub id is passed via label selectors then it takes precedence
			glog.V(2).Infof("override azure subscription id with : %s", value)
			merticReq.subscriptionID = value
		default:
			return azureMetricRequest{}, fmt.Errorf("selector label '%s' not supported", request.Key())
		}
	}

	err := merticReq.Validate()
	if err != nil {
		return azureMetricRequest{}, err
	}
	return merticReq, nil
}

func (amr azureMetricRequest) Validate() error {
	if amr.metricName == "" {
		return fmt.Errorf("metricName is required")
	}
	if amr.resourceGroup == "" {
		return fmt.Errorf("resourceGroup is required")
	}
	if amr.resourceName == "" {
		return fmt.Errorf("resourceName is required")
	}
	if amr.resourceProviderNamespace == "" {
		return fmt.Errorf("resourceProviderNamespace is required")
	}
	if amr.resourceType == "" {
		return fmt.Errorf("resourceType is required")
	}
	if amr.aggregation == "" {
		return fmt.Errorf("aggregation is required")
	}
	if amr.timespan == "" {
		return fmt.Errorf("timespan is required")
	}
	if amr.filter == "" {
		return fmt.Errorf("filter is required")
	}

	if amr.subscriptionID == "" {
		return fmt.Errorf("subscriptionID is required. set a default or pass via label selectors")
	}

	// if here then valid!
	return nil
}

func (amr azureMetricRequest) metricResourceURI() string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/%s/%s/%s",
		amr.subscriptionID,
		amr.resourceGroup,
		amr.resourceProviderNamespace,
		amr.resourceType,
		amr.resourceName)
}

func timeSpan() string {
	// defaults to last five minutes.
	// TODO support configuration via config
	endtime := time.Now().UTC().Format(time.RFC3339)
	starttime := time.Now().Add(-(5 * time.Minute)).UTC().Format(time.RFC3339)
	return fmt.Sprintf("%s/%s", starttime, endtime)
}

package azureMetricClient

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/selection"

	"github.com/golang/glog"
	"github.com/jsturtevant/azure-k8-metrics-adapter/pkg/aim"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/metrics/pkg/apis/external_metrics"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2018-03-01/insights"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

type AzureMetricClient struct {
	client         insights.MetricsClient
	subscriptionID string
}

func NewAzureMetricClient() AzureMetricClient {
	azureConfig, err := aim.GetAzureConfig()
	if err != nil {
		glog.Errorf("unable to get azure config: %v", err)
	}

	metricsClient := insights.NewMetricsClient(azureConfig.SubscriptionID)
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		metricsClient.Authorizer = authorizer
	}

	return AzureMetricClient{
		client:         metricsClient,
		subscriptionID: azureConfig.SubscriptionID,
	}
}

// GetAzureMetric calls arm rest endpoint
func (c AzureMetricClient) GetAzureMetric(metricSelector labels.Selector) (external_metrics.ExternalMetricValue, error) {
	azMetricRequest, err := parseAzureMetric(metricSelector)
	if err != nil {
		return external_metrics.ExternalMetricValue{}, err
	}
	metricResourceURI := azMetricRequest.metricResourceUri(c.subscriptionID)

	glog.V(2).Infof("resource uri: %s", metricResourceURI)
	glog.V(2).Infof("filter: %s", azMetricRequest.filter)
	glog.V(2).Infof("metric name : %s", azMetricRequest.metricName)

	// make call to azure resource provider
	metricResult, err := c.client.List(context.Background(), metricResourceURI,
		azMetricRequest.timespan, nil,
		azMetricRequest.metricName, azMetricRequest.aggregation, nil,
		"", azMetricRequest.filter, "", "")
	if err != nil {
		return external_metrics.ExternalMetricValue{}, err
	}

	//TODO check for nils
	metricVals := *metricResult.Value
	Timeseries := *metricVals[0].Timeseries
	data := *Timeseries[0].Data
	total := *data[len(data)-1].Total

	// TODO set Value based on aggregations type
	return external_metrics.ExternalMetricValue{
		MetricName: azMetricRequest.resourceName,
		Value:      *resource.NewQuantity(int64(total), resource.DecimalSI),
		Timestamp:  metav1.Now(),
	}, nil
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
}

func parseAzureMetric(metricSelector labels.Selector) (azureMetricRequest, error) {
	glog.V(2).Infof("begin parsing metric")

	// using selectors to pass required values thorugh
	// to retain case
	// there is are restrictions so using some conversion
	// restrictions here
	// note: requirement values are already validated by apiserver
	merticReq := azureMetricRequest{
		timespan: timeSpan(),
	}
	requirements, _ := metricSelector.Requirements()
	for _, request := range requirements {
		if request.Operator() != selection.Equals {
			return azureMetricRequest{}, errors.New("selector type not supported. only eqauls is supported at this time")
		}

		switch request.Key() {
		case "metricName":
			glog.V(2).Infof("metricName: %s", request.Values().List()[0])
			merticReq.metricName = request.Values().List()[0]
		case "resourceGroup":
			glog.V(2).Infof("resourceGroup: %s", request.Values().List()[0])
			merticReq.resourceGroup = request.Values().List()[0]
		case "resourceName":
			glog.V(2).Infof("resourceName: %s", request.Values().List()[0])
			merticReq.resourceName = request.Values().List()[0]
		case "resourceProviderNamespace":
			glog.V(2).Infof("resourceProviderNamespace: %s", request.Values().List()[0])
			merticReq.resourceProviderNamespace = request.Values().List()[0]
		case "resourceType":
			glog.V(2).Infof("resourceType: %s", request.Values().List()[0])
			merticReq.resourceType = request.Values().List()[0]
		case "aggregation":
			glog.V(2).Infof("aggregation: %s", request.Values().List()[0])
			merticReq.aggregation = request.Values().List()[0]
		case "filter":
			// TODO: Should handle filters by converting equality and setbased label selectors
			// to  oData syntax: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#label-selectors
			glog.V(2).Infof("filter: %s", request.Values().List()[0])
			filterStrings := strings.Split(request.Values().List()[0], "_")
			merticReq.filter = fmt.Sprintf("%s %s '%s'", filterStrings[0], filterStrings[1], filterStrings[2])
			glog.V(2).Infof("filter formatted: %s", merticReq.filter)
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

	// if here then valid!
	return nil
}

func (azr azureMetricRequest) metricResourceUri(subId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/%s/%s/%s",
		subId,
		azr.resourceGroup,
		azr.resourceProviderNamespace,
		azr.resourceType,
		azr.resourceName)
}

func timeSpan() string {
	// defaults to last five minutes.
	// TODO support configuration via config
	endtime := time.Now().UTC().Format(time.RFC3339)
	starttime := time.Now().Add(-(5 * time.Minute)).UTC().Format(time.RFC3339)
	return fmt.Sprintf("%s/%s", starttime, endtime)
}

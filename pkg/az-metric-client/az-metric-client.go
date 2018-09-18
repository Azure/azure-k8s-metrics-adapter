package azureMetricClient

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Azure/azure-k8s-metrics-adapter/pkg/aiapiclient"
	"github.com/Azure/azure-k8s-metrics-adapter/pkg/aim"
	"github.com/Azure/azure-k8s-metrics-adapter/pkg/azmetricrequest"
	"github.com/Azure/azure-k8s-metrics-adapter/pkg/metriccache"

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
	metriccache           *metriccache.MetricCache
}

// NewAzureMetricClient creates a client for making requests to Azure Monitor
func NewAzureMetricClient(metricCache *metriccache.MetricCache) AzureMetricClient {
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
		metriccache:           metricCache,
	}
}

// GetAzureMetric calls Azure Monitor endpoint and returns a metric based on label selectors
func (c AzureMetricClient) GetAzureMetric(namespace string, metricName string, metricSelector labels.Selector) (external_metrics.ExternalMetricValue, error) {

	azMetricRequest, err := c.getMetricRequest(namespace, metricName, metricSelector)
	if err != nil {
		return external_metrics.ExternalMetricValue{}, err
	}

	metricResourceURI := azMetricRequest.MetricResourceURI()
	glog.V(2).Infof("resource uri: %s", metricResourceURI)

	// make call to azure resource provider with subscription id provided issue #9
	c.monitorClient.SubscriptionID = azMetricRequest.SubscriptionID
	metricResult, err := c.monitorClient.List(context.Background(), metricResourceURI,
		azMetricRequest.Timespan, nil,
		azMetricRequest.MetricName, azMetricRequest.Aggregation, nil,
		"", azMetricRequest.Filter, "", "")
	if err != nil {
		return external_metrics.ExternalMetricValue{}, err
	}

	total := extractValue(metricResult)

	glog.V(2).Infof("found metric value: %f", total)

	// TODO set Value based on aggregations type
	return external_metrics.ExternalMetricValue{
		MetricName: azMetricRequest.ResourceName,
		Value:      *resource.NewQuantity(int64(total), resource.DecimalSI),
		Timestamp:  metav1.Now(),
	}, nil
}

func (c AzureMetricClient) getMetricRequest(namespace string, metricName string, metricSelector labels.Selector) (azmetricrequest.AzureMetricRequest, error) {
	key := metricKey(namespace, metricName)

	azMetricRequest, found := c.metriccache.GetMetric(key)
	if found {
		return azMetricRequest, nil
	}

	azMetricRequest, err := azmetricrequest.ParseAzureMetric(metricSelector, c.defaultSubscriptionID)
	if err != nil {
		return azmetricrequest.AzureMetricRequest{}, err
	}
	return azMetricRequest, nil
}

// GetCustomMetric calls to Application Insights to retrieve the value of the metric requested
func (c AzureMetricClient) GetCustomMetric(groupResource schema.GroupResource, namespace string, selector labels.Selector, metricName string) (float64, error) {
	// because metrics names are multipart in AI and we can not pass an extra /
	// through k8s api we convert - to / to get around that
	convertedMetricName := strings.Replace(metricName, "-", "/", -1)
	glog.V(2).Infof("New call to GetCustomMetric: %s", convertedMetricName)

	// get the last 5 mins and chunking into 30 seconds
	// this seems to be the best way to get the closest average rate at time of request
	// any smaller time intervals and the values come back null
	metricRequestInfo := aiapiclient.NewMetricRequest(convertedMetricName)
	metricRequestInfo.Timespan = "PT5M"
	metricRequestInfo.Interval = "PT30S"

	metricsResult, err := c.appinsightsclient.GetMetric(metricRequestInfo)
	if err != nil {
		return 0, err
	}

	if metricsResult.Value == nil || metricsResult.Value.Segments == nil {
		return 0, errors.New("metrics result is nil")
	}

	segments := *metricsResult.Value.Segments
	if len(segments) <= 0 {
		glog.V(2).Info("segments length = 0")
		return 0, nil
	}

	// grab just the last value which will be the latest value of the metric
	metric := segments[len(segments)-1].AdditionalProperties[convertedMetricName]
	metricMap := metric.(map[string]interface{})
	value := metricMap["avg"]
	normalizedValue := normalizeValue(value)

	glog.V(2).Infof("found metric value: %f", normalizedValue)
	return normalizedValue, nil
}

func normalizeValue(value interface{}) float64 {
	switch t := value.(type) {
	case int32:
		return float64(value.(int32))
	case float32:
		return float64(value.(float32))
	case float64:
		return value.(float64)
	case int64:
		return float64(value.(int64))
	default:
		glog.V(0).Infof("unexpected type: %T", t)
		return 0
	}
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

func metricKey(namespace string, name string) string {
	return fmt.Sprintf("%s/%s", namespace, name)
}

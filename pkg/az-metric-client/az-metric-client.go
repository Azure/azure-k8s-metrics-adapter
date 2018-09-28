package azureMetricClient

import (
	"errors"
	"strings"

	"github.com/Azure/azure-k8s-metrics-adapter/pkg/azure/appinsights"

	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// AzureMetricClient is used to make requests to Azure Monitor
type AzureMetricClient struct {
	appinsightsclient     appinsights.AiAPIClient
	defaultSubscriptionID string
}

// NewAzureMetricClient creates a client for making requests to Azure Monitor
func NewAzureMetricClient(defaultSubscriptionID string) AzureMetricClient {
	appInsightsClient := appinsights.NewAiAPIClient()

	return AzureMetricClient{
		appinsightsclient:     appInsightsClient,
		defaultSubscriptionID: defaultSubscriptionID,
	}
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
	metricRequestInfo := appinsights.NewMetricRequest(convertedMetricName)
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

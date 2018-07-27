package azureMetricClient

import (
	"context"
	"fmt"
	"time"

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
	glog.V(2).Infof("requirement")

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

func (c AzureMetricClient) Do(namespace string, metricName string, metricSelector labels.Selector) (external_metrics.ExternalMetricValue, error) {
	metricName = "Messages"
	metricResourceUri := metricResourceUri(c.subscriptionID, "k8metrics", "k8custom")

	endtime := time.Now().UTC().Format(time.RFC3339)
	starttime := time.Now().Add(-(5 * time.Minute)).UTC().Format(time.RFC3339)
	timespan := fmt.Sprintf("%s/%s", starttime, endtime)

	metricResult, err := c.client.List(context.Background(), metricResourceUri, timespan, nil, metricName, "Total", nil, "", "", "", "")
	if err != nil {
		return external_metrics.ExternalMetricValue{}, err
	}

	metricVals := *metricResult.Value
	Timeseries := *metricVals[0].Timeseries
	data := *Timeseries[0].Data
	total := *data[len(data)-1].Total

	return external_metrics.ExternalMetricValue{
		MetricName: metricName,
		Value:      *resource.NewQuantity(int64(total), resource.DecimalSI),
		Timestamp:  metav1.Now(),
	}, nil

}

func metricResourceUri(subId string, resourceGroup string, sbNameSpace string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ServiceBus/namespaces/%s", subId, resourceGroup, sbNameSpace)
}

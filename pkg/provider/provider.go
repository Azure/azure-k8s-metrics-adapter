package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/jsturtevant/azure-k8-metrics-adapter/pkg/aim"
	"k8s.io/apimachinery/pkg/api/errors"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/metrics/pkg/apis/custom_metrics"

	"github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/provider"
	"k8s.io/metrics/pkg/apis/external_metrics"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2018-03-01/insights"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

type externalMetric struct {
	info  provider.ExternalMetricInfo
	value external_metrics.ExternalMetricValue
}

type AzureProvider struct {
	client      dynamic.Interface
	mapper      apimeta.RESTMapper
	azureConfig *aim.AzureConfig

	values          map[provider.CustomMetricInfo]int64
	externalMetrics []externalMetric
}

func NewAzureProvider(client dynamic.Interface, mapper apimeta.RESTMapper) provider.MetricsProvider {
	azureConfig, err := aim.GetAzureConfig()
	if err != nil {
		glog.Errorf("unable to get azure config: %v", err)
	}

	return &AzureProvider{
		client:      client,
		mapper:      mapper,
		azureConfig: azureConfig,
		values:      make(map[provider.CustomMetricInfo]int64),
	}
}

/* Custom metric interface methods */
// not implemented
func (p *AzureProvider) GetRootScopedMetricByName(groupResource schema.GroupResource, name string, metricName string) (*custom_metrics.MetricValue, error) {
	//not implemented yet
	return nil, errors.NewServiceUnavailable("not implemented yet")
}

// not implemented
func (p *AzureProvider) GetRootScopedMetricBySelector(groupResource schema.GroupResource, selector labels.Selector, metricName string) (*custom_metrics.MetricValueList, error) {
	// not implemented yet
	return nil, errors.NewServiceUnavailable("not implemented yet")
}

// not implemented
func (p *AzureProvider) GetNamespacedMetricByName(groupResource schema.GroupResource, namespace string, name string, metricName string) (*custom_metrics.MetricValue, error) {
	// not implemented yet
	return nil, errors.NewServiceUnavailable("not implemented yet")
}

// not implemented
func (p *AzureProvider) GetNamespacedMetricBySelector(groupResource schema.GroupResource, namespace string, selector labels.Selector, metricName string) (*custom_metrics.MetricValueList, error) {
	// not implemented yet
	return nil, errors.NewServiceUnavailable("not implemented yet")
}

func (p *AzureProvider) ListAllMetrics() []provider.CustomMetricInfo {
	// not implemented yet
	return []provider.CustomMetricInfo{}
}

func (p *AzureProvider) GetExternalMetric(namespace string, metricName string, metricSelector labels.Selector) (*external_metrics.ExternalMetricValueList, error) {
	glog.V(2).Infof("Recieved request for namespace: %s, metric name: %s, metric selectors: %s", namespace, metricName, metricSelector.String())

	requirements, selectable := metricSelector.Requirements()
	if !selectable {
		return nil, errors.NewBadRequest("label is set to not selectable. this should not happen")
	}
	for _, req := range requirements {
		glog.V(2).Infof("requirement: %s: %s", req.Key(), req.Values())
	}

	// create an authorizer from env vars or Azure Managed Service Idenity
	metricsClient := insights.NewMetricsClient(p.azureConfig.SubscriptionID)
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		metricsClient.Authorizer = authorizer
	}

	metricName = "Messages"
	metricResourceUri := metricResourceUri(p.azureConfig.SubscriptionID, "k8metrics", "k8custom")

	endtime := time.Now().UTC().Format(time.RFC3339)
	starttime := time.Now().Add(-(5 * time.Minute)).UTC().Format(time.RFC3339)
	timespan := fmt.Sprintf("%s/%s", starttime, endtime)

	metricResult, err := metricsClient.List(context.Background(), metricResourceUri, timespan, nil, metricName, "Total", nil, "", "", "", "")
	if err != nil {
		return nil, err
	}

	metricVals := *metricResult.Value
	Timeseries := *metricVals[0].Timeseries
	data := *Timeseries[0].Data
	total := *data[len(data)-1].Total

	metricValue := external_metrics.ExternalMetricValue{
		MetricName: metricName,
		Value:      *resource.NewQuantity(int64(total), resource.DecimalSI),
		Timestamp:  metav1.Now(),
	}

	matchingMetrics := []external_metrics.ExternalMetricValue{}
	matchingMetrics = append(matchingMetrics, metricValue)

	return &external_metrics.ExternalMetricValueList{
		Items: matchingMetrics,
	}, nil
}

func metricResourceUri(subId string, resourceGroup string, sbNameSpace string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ServiceBus/namespaces/%s", subId, resourceGroup, sbNameSpace)
}

func (p *AzureProvider) ListAllExternalMetrics() []provider.ExternalMetricInfo {
	externalMetricsInfo := []provider.ExternalMetricInfo{}

	// not implemented yet
	// TODO
	// iterate over all of the resources we have access to
	// build metric info from that
	// important to remember to cache this and only get it at given interval

	for _, metric := range p.externalMetrics {
		externalMetricsInfo = append(externalMetricsInfo, metric.info)
	}
	return externalMetricsInfo
}

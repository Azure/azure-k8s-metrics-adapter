package provider

import (
	"github.com/golang/glog"
	"github.com/jsturtevant/azure-k8-metrics-adapter/pkg/aim"
	"github.com/jsturtevant/azure-k8-metrics-adapter/pkg/az-metric-client"
	"k8s.io/apimachinery/pkg/api/errors"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/metrics/pkg/apis/custom_metrics"

	"github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/provider"
	"k8s.io/metrics/pkg/apis/external_metrics"
)

type AzureProvider struct {
	client         dynamic.Interface
	mapper         apimeta.RESTMapper
	azureConfig    aim.AzureConfig
	azMetricClient azureMetricClient.AzureMetricClient
}

func NewAzureProvider(client dynamic.Interface, mapper apimeta.RESTMapper, azMetricClient azureMetricClient.AzureMetricClient) provider.MetricsProvider {
	return &AzureProvider{
		client:         client,
		mapper:         mapper,
		azMetricClient: azMetricClient,
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
	// Note:
	//		namespace is Kubernetes namespace when using hpa.
	// 		This doesn't have affect on azure resources so is ignored.
	//
	//		metric name is also ignored as azure metric name is case sensitve
	//		and this metric name is passed via url which removes case
	glog.V(0).Infof("Received request for namespace: %s, metric name: %s, metric selectors: %s", namespace, metricName, metricSelector.String())

	_, selectable := metricSelector.Requirements()
	if !selectable {
		return nil, errors.NewBadRequest("label is set to not selectable. this should not happen")
	}

	metricValue, err := p.azMetricClient.GetAzureMetric(metricSelector)
	if err != nil {
		glog.Errorf("bad request: %v", err)
		return nil, errors.NewBadRequest(err.Error())
	}

	matchingMetrics := []external_metrics.ExternalMetricValue{}
	matchingMetrics = append(matchingMetrics, metricValue)

	return &external_metrics.ExternalMetricValueList{
		Items: matchingMetrics,
	}, nil
}

func (p *AzureProvider) ListAllExternalMetrics() []provider.ExternalMetricInfo {
	externalMetricsInfo := []provider.ExternalMetricInfo{}

	// not implemented yet

	// TODO
	// iterate over all of the resources we have access
	// build metric info from https://docs.microsoft.com/en-us/azure/monitoring-and-diagnostics/monitoring-rest-api-walkthrough#retrieve-metric-definitions-multi-dimensional-api
	// important to remember to cache this and only get it at given interval

	return externalMetricsInfo
}

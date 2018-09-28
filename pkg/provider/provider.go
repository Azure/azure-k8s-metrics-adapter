// Package provider is the implementation of custom metric and external metric apis
// see https://github.com/kubernetes/community/blob/master/contributors/design-proposals/instrumentation/custom-metrics-api.md#api-paths
package provider

import (
	"fmt"
	"time"

	"github.com/Azure/azure-k8s-metrics-adapter/pkg/metriccache"

	"github.com/Azure/azure-k8s-metrics-adapter/pkg/azure/monitor"

	"github.com/Azure/azure-k8s-metrics-adapter/pkg/az-metric-client"
	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/api/errors"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/metrics/pkg/apis/custom_metrics"

	"github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/provider"
	"k8s.io/metrics/pkg/apis/external_metrics"
)

type AzureProvider struct {
	azMetricClient azureMetricClient.AzureMetricClient
	mapper         apimeta.RESTMapper
	monitorClient  monitor.AzureMetricClient
	metricCache    *metriccache.MetricCache
}

func NewAzureProvider(mapper apimeta.RESTMapper, azMetricClient azureMetricClient.AzureMetricClient, monitorClient monitor.AzureMetricClient, metricCache *metriccache.MetricCache) provider.MetricsProvider {
	return &AzureProvider{
		azMetricClient: azMetricClient,
		monitorClient:  monitorClient,
		metricCache:    metricCache,
	}
}

// GetMetricByName fetches a particular metric for a particular object.
// The namespace will be empty if the metric is root-scoped.
func (p *AzureProvider) GetMetricByName(name types.NamespacedName, info provider.CustomMetricInfo) (*custom_metrics.MetricValue, error) {
	// not implemented yet
	return nil, errors.NewServiceUnavailable("not implemented yet")
}

// GetMetricBySelector fetches a particular metric for a set of objects matching
// the given label selector.  The namespace will be empty if the metric is root-scoped.
func (p *AzureProvider) GetMetricBySelector(namespace string, selector labels.Selector, info provider.CustomMetricInfo) (*custom_metrics.MetricValueList, error) {
	glog.V(0).Infof("Received request for custom metric: groupresource: %s, namespace: %s, metric name: %s, selectors: %s", info.GroupResource.String(), namespace, info.Metric, selector.String())

	_, selectable := selector.Requirements()
	if !selectable {
		return nil, errors.NewBadRequest("label is set to not selectable. this should not happen")
	}

	val, err := p.azMetricClient.GetCustomMetric(info.GroupResource, namespace, selector, info.Metric)
	if err != nil {
		glog.Errorf("bad request: %v", err)
		return nil, errors.NewBadRequest(err.Error())
	}

	// TODO what does version do?
	kind, err := p.mapper.KindFor(info.GroupResource.WithVersion(""))
	if err != nil {
		return nil, errors.NewBadRequest(err.Error())
	}

	metricValue := custom_metrics.MetricValue{
		DescribedObject: custom_metrics.ObjectReference{
			APIVersion: info.GroupResource.Group + "/" + runtime.APIVersionInternal,
			Kind:       kind.Kind,
			Name:       info.Metric,
			Namespace:  namespace,
		},
		MetricName: info.Metric,
		Timestamp:  metav1.Time{time.Now()},
		Value:      *resource.NewMilliQuantity(int64(val*1000), resource.DecimalSI),
	}

	metricList := make([]custom_metrics.MetricValue, 0)
	metricList = append(metricList, metricValue)

	return &custom_metrics.MetricValueList{
		Items: metricList,
	}, nil
}

// ListAllMetrics provides a list of all available metrics at
// the current time.  Note that this is not allowed to return
// an error, so it is reccomended that implementors cache and
// periodically update this list, instead of querying every time.
func (p *AzureProvider) ListAllMetrics() []provider.CustomMetricInfo {
	// not implemented yet
	return []provider.CustomMetricInfo{}
}

// GetExternalMetric retrieves metrics from Azure Monitor Endpoint
// Metric is normally identified by a name and a set of labels/tags. It is up to a specific
// implementation how to translate metricSelector to a filter for metric values.
// Namespace can be used by the implementation for metric identification, access control or ignored.
func (p *AzureProvider) GetExternalMetric(namespace string, metricSelector labels.Selector, info provider.ExternalMetricInfo) (*external_metrics.ExternalMetricValueList, error) {
	// Note:
	//		namespace is Kubernetes namespace when using hpa.
	// 		This doesn't have affect on azure resources so is ignored.
	//
	//		metric name is also ignored as azure metric name is case sensitve
	//		and this metric name is passed via url which removes case
	glog.V(0).Infof("Received request for namespace: %s, metric name: %s, metric selectors: %s", namespace, info.Metric, metricSelector.String())

	_, selectable := metricSelector.Requirements()
	if !selectable {
		return nil, errors.NewBadRequest("label is set to not selectable. this should not happen")
	}

	azMetricRequest, err := p.getMetricRequest(namespace, info.Metric, metricSelector)
	if err != nil {
		return nil, errors.NewBadRequest(err.Error())
	}

	metricValue, err := p.monitorClient.GetAzureMetric(azMetricRequest)
	if err != nil {
		glog.Errorf("bad request: %v", err)
		return nil, errors.NewBadRequest(err.Error())
	}

	externalmetric := external_metrics.ExternalMetricValue{
		MetricName: azMetricRequest.MetricName,
		Value:      *resource.NewQuantity(int64(metricValue.Total), resource.DecimalSI),
		Timestamp:  metav1.Now(),
	}

	matchingMetrics := []external_metrics.ExternalMetricValue{}
	matchingMetrics = append(matchingMetrics, externalmetric)

	return &external_metrics.ExternalMetricValueList{
		Items: matchingMetrics,
	}, nil
}

// ListAllExternalMetrics calls out to azure and builds a list of metrics that can be queried against
func (p *AzureProvider) ListAllExternalMetrics() []provider.ExternalMetricInfo {
	externalMetricsInfo := []provider.ExternalMetricInfo{}

	// not implemented yet

	// TODO
	// iterate over all of the resources we have access
	// build metric info from https://docs.microsoft.com/en-us/azure/monitoring-and-diagnostics/monitoring-rest-api-walkthrough#retrieve-metric-definitions-multi-dimensional-api
	// important to remember to cache this and only get it at given interval

	return externalMetricsInfo
}

func (p *AzureProvider) getMetricRequest(namespace string, metricName string, metricSelector labels.Selector) (monitor.AzureMetricRequest, error) {
	key := metricKey(namespace, metricName)

	azMetricRequest, found := p.metricCache.Get(key)
	if found {
		azMetricRequest.Timespan = monitor.TimeSpan()
		if azMetricRequest.SubscriptionID == "" {
			azMetricRequest.SubscriptionID = p.monitorClient.DefaultSubscriptionID
		}
		return azMetricRequest, nil
	}

	azMetricRequest, err := monitor.ParseAzureMetric(metricSelector, p.monitorClient.DefaultSubscriptionID)
	if err != nil {
		return monitor.AzureMetricRequest{}, err
	}

	return azMetricRequest, nil
}

func metricKey(namespace string, name string) string {
	return fmt.Sprintf("%s/%s", namespace, name)
}

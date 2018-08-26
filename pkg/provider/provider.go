package provider

import (
	"time"

	"github.com/Azure/azure-k8s-metrics-adapter/pkg/aim"
	"github.com/Azure/azure-k8s-metrics-adapter/pkg/az-metric-client"
	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/api/errors"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
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

/* Custom metric interface methods
see https://github.com/kubernetes/community/blob/master/contributors/design-proposals/instrumentation/custom-metrics-api.md#api-paths
*/

// GetRootScopedMetricByName fetches a particular metric for a particular root-scoped object (such as Node, PersistentVolume)
func (p *AzureProvider) GetRootScopedMetricByName(groupResource schema.GroupResource, name string, metricName string) (*custom_metrics.MetricValue, error) {
	//not implemented yet
	return nil, errors.NewServiceUnavailable("not implemented yet")
}

// GetRootScopedMetricBySelector fetches a particular metric for a set of root-scoped objects (such as Node, PersistentVolume)
// matching the given label selector.
func (p *AzureProvider) GetRootScopedMetricBySelector(groupResource schema.GroupResource, selector labels.Selector, metricName string) (*custom_metrics.MetricValueList, error) {
	// not implemented yet
	return nil, errors.NewServiceUnavailable("not implemented yet")
}

// GetNamespacedMetricByName fetches a particular metric for a particular namespaced object (such as pod, deployment)
func (p *AzureProvider) GetNamespacedMetricByName(groupResource schema.GroupResource, namespace string, name string, metricName string) (*custom_metrics.MetricValue, error) {
	glog.V(0).Infof("Received request for custom metric: groupresource: %s, namespace: %s, name: %s, metric name: %s", groupResource.String(), namespace, name, metricName)

	return nil, errors.NewServiceUnavailable("not implemented yet")
}

// GetNamespacedMetricBySelector fetches a particular metric for a set of namespaced objects (such as pod, deployment)
// matching the given label selector.
func (p *AzureProvider) GetNamespacedMetricBySelector(groupResource schema.GroupResource, namespace string, selector labels.Selector, metricName string) (*custom_metrics.MetricValueList, error) {
	glog.V(0).Infof("Received request for custom metric: groupresource: %s, namespace: %s, metric name: %s, selectors: %s", groupResource.String(), namespace, metricName, selector.String())

	_, selectable := selector.Requirements()
	if !selectable {
		return nil, errors.NewBadRequest("label is set to not selectable. this should not happen")
	}

	val, err := p.azMetricClient.GetCustomMetric(groupResource, namespace, selector, metricName)
	if err != nil {
		glog.Errorf("bad request: %v", err)
		return nil, errors.NewBadRequest(err.Error())
	}

	// TODO what does version do?
	kind, err := p.mapper.KindFor(groupResource.WithVersion(""))
	if err != nil {
		return nil, errors.NewBadRequest(err.Error())
	}

	metricValue := custom_metrics.MetricValue{
		DescribedObject: custom_metrics.ObjectReference{
			APIVersion: groupResource.Group + "/" + runtime.APIVersionInternal,
			Kind:       kind.Kind,
			Name:       metricName,
			Namespace:  namespace,
		},
		MetricName: metricName,
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

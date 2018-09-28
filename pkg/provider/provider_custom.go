// Package provider is the implementation of custom metric and external metric apis
// see https://github.com/kubernetes/community/blob/master/contributors/design-proposals/instrumentation/custom-metrics-api.md#api-paths
package provider

import (
	"time"

	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/metrics/pkg/apis/custom_metrics"

	"github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/provider"
)

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

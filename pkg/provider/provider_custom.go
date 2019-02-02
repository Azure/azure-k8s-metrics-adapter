// Package provider is the implementation of custom metric and external metric apis
// see https://github.com/kubernetes/community/blob/master/contributors/design-proposals/instrumentation/custom-metrics-api.md#api-paths
package provider

import (
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-k8s-metrics-adapter/pkg/azure/custommetrics"
	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/metrics/pkg/apis/custom_metrics"

	"github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/provider"
	"github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/provider/helpers"
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

	metricRequestInfo := p.getCustomMetricRequest(namespace, selector, info)

	// TODO use selector info to restrict metric query to specific app.
	val, err := p.appinsightsClient.GetCustomMetric(metricRequestInfo)
	if err != nil {
		glog.Errorf("bad request: %v", err)
		return nil, errors.NewBadRequest(err.Error())
	}

	resourceNames, err := helpers.ListObjectNames(p.mapper, p.kubeClient, namespace, selector, info)
	if err != nil {
		glog.Errorf("not able to list objects from api server: %v", err)
		return nil, errors.NewInternalError(fmt.Errorf("not able to list objects from api server for this resource"))
	}

	// TODO: Add support for app insights where pods are mapped 1 to 1.
	// Currently App insights does not out of the box support kubernetes pod information
	// so we are using the value from AI and passing to all instances of the pods.
	// We should be passing pod level metric info to App insights but there is currently on the developer to wire that up and
	// maping it here based on pod name.
	metricList := make([]custom_metrics.MetricValue, 0)
	for _, name := range resourceNames {
		ref, err := helpers.ReferenceFor(p.mapper, types.NamespacedName{Namespace: namespace, Name: name}, info)
		if err != nil {
			return nil, err
		}

		metricValue := custom_metrics.MetricValue{
			DescribedObject: ref,
			Metric: custom_metrics.MetricIdentifier{
				Name: info.Metric,
			},
			Timestamp: metav1.Time{time.Now()},
			Value:     *resource.NewMilliQuantity(int64(val*1000), resource.DecimalSI),
		}

		// add back the meta data about the request selectors
		if len(selector.String()) > 0 {
			labelSelector, err := metav1.ParseToLabelSelector(selector.String())
			if err != nil {
				return nil, err
			}
			metricValue.Metric.Selector = labelSelector
		}

		metricList = append(metricList, metricValue)
	}

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

func (p *AzureProvider) getCustomMetricRequest(namespace string, selector labels.Selector, info provider.CustomMetricInfo) custommetrics.MetricRequest {

	cachedRequest, found := p.metricCache.GetAppInsightsRequest(namespace, info.Metric)
	if found {
		return cachedRequest
	}

	// because metrics names are multipart in AI and we can not pass an extra /
	// through k8s api we convert - to / to get around that
	convertedMetricName := strings.Replace(info.Metric, "-", "/", -1)
	glog.V(2).Infof("New call to GetCustomMetric: %s", convertedMetricName)
	metricRequestInfo := custommetrics.NewMetricRequest(convertedMetricName)

	return metricRequestInfo
}

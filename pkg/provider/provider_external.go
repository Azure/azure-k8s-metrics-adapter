// Package provider is the implementation of custom metric and external metric apis
// see https://github.com/kubernetes/community/blob/master/contributors/design-proposals/instrumentation/custom-metrics-api.md#api-paths
package provider

import (
	"github.com/Azure/azure-k8s-metrics-adapter/pkg/azure/externalmetrics"
	"github.com/golang/glog"
	"github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/provider"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/metrics/pkg/apis/external_metrics"
)

// GetExternalMetric retrieves metrics from Azure Monitor Endpoint
// Metric is normally identified by a name and a set of labels/tags. It is up to a specific
// implementation how to translate metricSelector to a filter for metric values.
// Namespace can be used by the implementation for metric identification, access control or ignored.
func (p *AzureProvider) GetExternalMetric(namespace string, metricSelector labels.Selector, info provider.ExternalMetricInfo) (*external_metrics.ExternalMetricValueList, error) {
	// Note:
	//		metric name and namespace is used to lookup for the CRD which contains configuration to call azure
	// 		if not found then ignored and label selector is parsed for all the metrics
	glog.V(0).Infof("Received request for namespace: %s, metric name: %s, metric selectors: %s", namespace, info.Metric, metricSelector.String())

	_, selectable := metricSelector.Requirements()
	if !selectable {
		return nil, errors.NewBadRequest("label is set to not selectable. this should not happen")
	}

	azMetricRequest, err := p.getMetricRequest(namespace, info.Metric, metricSelector)
	if err != nil {
		return nil, errors.NewBadRequest(err.Error())
	}

	externalMetricClient, err := p.azureClientFactory.GetAzureExternalMetricClient(azMetricRequest.Type)
	if err != nil {
		return nil, errors.NewBadRequest(err.Error())
	}

	metricValue, err := externalMetricClient.GetAzureMetric(azMetricRequest)
	if err != nil {
		glog.Errorf("bad request: %v", err)
		return nil, errors.NewBadRequest(err.Error())
	}

	externalmetric := external_metrics.ExternalMetricValue{
		MetricName: info.Metric,
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

func (p *AzureProvider) getMetricRequest(namespace string, metricName string, metricSelector labels.Selector) (externalmetrics.AzureExternalMetricRequest, error) {

	azMetricRequest, found := p.metricCache.GetAzureExternalMetricRequest(namespace, metricName)
	if found {
		azMetricRequest.Timespan = externalmetrics.TimeSpan()
		if azMetricRequest.SubscriptionID == "" {
			azMetricRequest.SubscriptionID = p.defaultSubscriptionID
		}
		return azMetricRequest, nil
	}

	azMetricRequest, err := externalmetrics.ParseAzureMetric(metricSelector, p.defaultSubscriptionID)
	if err != nil {
		return externalmetrics.AzureExternalMetricRequest{}, err
	}

	return azMetricRequest, nil
}

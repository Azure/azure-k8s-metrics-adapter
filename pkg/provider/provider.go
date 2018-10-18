// Package provider is the implementation of custom metric and external metric apis
// see https://github.com/kubernetes/community/blob/master/contributors/design-proposals/instrumentation/custom-metrics-api.md#api-paths
package provider

import (
	"github.com/Azure/azure-k8s-metrics-adapter/pkg/azure/appinsights"
	"github.com/Azure/azure-k8s-metrics-adapter/pkg/metriccache"

	"github.com/Azure/azure-k8s-metrics-adapter/pkg/azure/monitor"

	apimeta "k8s.io/apimachinery/pkg/api/meta"

	"github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/provider"
)

type AzureProvider struct {
	appinsightsClient     appinsights.AzureAppInsightsClient
	mapper                apimeta.RESTMapper
	monitorClient         monitor.AzureMonitorClient
	metricCache           *metriccache.MetricCache
	defaultSubscriptionID string
}

func NewAzureProvider(defaultSubscriptionID string, mapper apimeta.RESTMapper, appinsightsClient appinsights.AzureAppInsightsClient, monitorClient monitor.AzureMonitorClient, metricCache *metriccache.MetricCache) provider.MetricsProvider {
	return &AzureProvider{
		defaultSubscriptionID: defaultSubscriptionID,
		mapper:                mapper,
		appinsightsClient:     appinsightsClient,
		monitorClient:         monitorClient,
		metricCache:           metricCache,
	}
}

// Package provider is the implementation of custom metric and external metric apis
// see https://github.com/kubernetes/community/blob/master/contributors/design-proposals/instrumentation/custom-metrics-api.md#api-paths
package provider

import (
	"github.com/Azure/azure-k8s-metrics-adapter/pkg/azure/custommetrics"
	"github.com/Azure/azure-k8s-metrics-adapter/pkg/azure/externalmetrics"
	"github.com/Azure/azure-k8s-metrics-adapter/pkg/metriccache"
	"github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/provider"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/dynamic"
)

type AzureProvider struct {
	appinsightsClient     custommetrics.AzureAppInsightsClient
	mapper                apimeta.RESTMapper
	kubeClient            dynamic.Interface
	metricCache           *metriccache.MetricCache
	azureClientFactory    externalmetrics.AzureClientFactory
	defaultSubscriptionID string
}

func NewAzureProvider(defaultSubscriptionID string, mapper apimeta.RESTMapper, kubeClient dynamic.Interface, appinsightsClient custommetrics.AzureAppInsightsClient, azureClientFactory externalmetrics.AzureClientFactory, metricCache *metriccache.MetricCache) provider.MetricsProvider {
	return &AzureProvider{
		defaultSubscriptionID: defaultSubscriptionID,
		mapper:                mapper,
		kubeClient:            kubeClient,
		appinsightsClient:     appinsightsClient,
		metricCache:           metricCache,
		azureClientFactory:    azureClientFactory,
	}
}

package controller

import (
	"fmt"

	"github.com/Azure/azure-k8s-metrics-adapter/pkg/azure/custommetrics"
	"github.com/Azure/azure-k8s-metrics-adapter/pkg/azure/externalmetrics"
	listers "github.com/Azure/azure-k8s-metrics-adapter/pkg/client/listers/metrics/v1alpha2"
	"github.com/Azure/azure-k8s-metrics-adapter/pkg/metriccache"
	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"
)

// Handler processes the events from the controler for external metrics
type Handler struct {
	externalmetricLister listers.ExternalMetricLister
	metriccache          *metriccache.MetricCache
	customMetricLister   listers.CustomMetricLister
}

// NewHandler created a new handler
func NewHandler(externalmetricLister listers.ExternalMetricLister, customMetricLister listers.CustomMetricLister, metricCache *metriccache.MetricCache) Handler {
	return Handler{
		externalmetricLister: externalmetricLister,
		customMetricLister:   customMetricLister,
		metriccache:          metricCache,
	}
}

type ControllerHandler interface {
	Process(queueItem namespacedQueueItem) error
}

// Process validates the item exists then stores updates the metric cached used to make requests to azure
func (h *Handler) Process(queueItem namespacedQueueItem) error {
	ns, name, err := cache.SplitMetaNamespaceKey(queueItem.namespaceKey)
	if err != nil {
		// not a valid key do not put back on queue
		runtime.HandleError(fmt.Errorf("expected namespace/name key in workqueue but got %s", queueItem.namespaceKey))
		return err
	}

	switch queueItem.kind {
	case "CustomMetric":
		return h.handleCustomMetric(ns, name, queueItem)
	case "ExternalMetric":
		return h.handleExternalMetric(ns, name, queueItem)
	}

	return nil
}

func (h *Handler) handleCustomMetric(ns, name string, queueItem namespacedQueueItem) error {
	// check if item exists
	glog.V(2).Infof("processing item '%s' in namespace '%s'", name, ns)
	customMetricInfo, err := h.customMetricLister.CustomMetrics(ns).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			// Then this we should remove
			glog.V(2).Infof("removing item from cache '%s' in namespace '%s'", name, ns)
			h.metriccache.Remove(queueItem.Key())
			return nil
		}

		return err
	}

	metric := custommetrics.MetricRequest{
		MetricName: customMetricInfo.Spec.MetricConfig.MetricName,
	}

	glog.V(2).Infof("adding to cache item '%s' in namespace '%s'", name, ns)
	h.metriccache.Update(queueItem.Key(), metric)

	return nil
}

func (h *Handler) handleExternalMetric(ns, name string, queueItem namespacedQueueItem) error {
	// check if item exists
	glog.V(2).Infof("processing item '%s' in namespace '%s'", name, ns)
	externalMetricInfo, err := h.externalmetricLister.ExternalMetrics(ns).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			// Then this we should remove
			glog.V(2).Infof("removing item from cache '%s' in namespace '%s'", name, ns)
			h.metriccache.Remove(queueItem.Key())
			return nil
		}

		return err
	}

	// TODO: Map the new fields here for Service Bus
	azureMetricRequest := externalmetrics.AzureExternalMetricRequest{
		ResourceGroup:             externalMetricInfo.Spec.AzureConfig.ResourceGroup,
		ResourceName:              externalMetricInfo.Spec.AzureConfig.ResourceName,
		ResourceProviderNamespace: externalMetricInfo.Spec.AzureConfig.ResourceProviderNamespace,
		ResourceType:              externalMetricInfo.Spec.AzureConfig.ResourceType,
		SubscriptionID:            externalMetricInfo.Spec.AzureConfig.SubscriptionID,
		MetricName:                externalMetricInfo.Spec.MetricConfig.MetricName,
		Filter:                    externalMetricInfo.Spec.MetricConfig.Filter,
		Aggregation:               externalMetricInfo.Spec.MetricConfig.Aggregation,
		Topic:                     externalMetricInfo.Spec.AzureConfig.ServiceBusTopic,
		Type:                      externalMetricInfo.Spec.Type,
		Namespace:                 externalMetricInfo.Spec.AzureConfig.ServiceBusNamespace,
		Subscription:              externalMetricInfo.Spec.AzureConfig.ServiceBusSubscription,
	}

	glog.V(2).Infof("adding to cache item '%s' in namespace '%s'", name, ns)
	h.metriccache.Update(queueItem.Key(), azureMetricRequest)

	return nil
}

package metriccache

import (
	"fmt"
	"sync"

	"github.com/Azure/azure-k8s-metrics-adapter/pkg/azure/custommetrics"
	"github.com/Azure/azure-k8s-metrics-adapter/pkg/azure/externalmetrics"
	"github.com/golang/glog"
)

// MetricCache holds the loaded metric request info in the system
type MetricCache struct {
	metricMutext   sync.RWMutex
	metricRequests map[string]interface{}
}

// NewMetricCache creates the cache
func NewMetricCache() *MetricCache {
	return &MetricCache{
		metricRequests: make(map[string]interface{}),
	}
}

// Update sets a metric request in the cache
func (mc *MetricCache) Update(key string, metricRequest interface{}) {
	mc.metricMutext.Lock()
	defer mc.metricMutext.Unlock()

	mc.metricRequests[key] = metricRequest
}

// GetAzureExternalMetricRequest retrieves a metric request from the cache
func (mc *MetricCache) GetAzureExternalMetricRequest(namepace, name string) (externalmetrics.AzureExternalMetricRequest, bool) {
	mc.metricMutext.RLock()
	defer mc.metricMutext.RUnlock()

	key := externalMetricKey(namepace, name)
	metricRequest, exists := mc.metricRequests[key]
	if !exists {
		glog.V(2).Infof("metric not found %s", key)
		return externalmetrics.AzureExternalMetricRequest{}, false
	}

	return metricRequest.(externalmetrics.AzureExternalMetricRequest), true
}

// GetAppInsightsRequest retrieves a metric request from the cache
func (mc *MetricCache) GetAppInsightsRequest(namespace, name string) (custommetrics.MetricRequest, bool) {
	mc.metricMutext.RLock()
	defer mc.metricMutext.RUnlock()

	key := customMetricKey(namespace, name)
	metricRequest, exists := mc.metricRequests[key]
	if !exists {
		glog.V(2).Infof("metric not found %s", key)
		return custommetrics.MetricRequest{}, false
	}

	return metricRequest.(custommetrics.MetricRequest), true
}

// Remove retrieves a metric request from the cache
func (mc *MetricCache) Remove(key string) {
	mc.metricMutext.Lock()
	defer mc.metricMutext.Unlock()

	delete(mc.metricRequests, key)
}

func externalMetricKey(namespace string, name string) string {
	return fmt.Sprintf("ExternalMetric/%s/%s", namespace, name)
}

func customMetricKey(namespace string, name string) string {
	return fmt.Sprintf("CustomMetric/%s/%s", namespace, name)
}

package metriccache

import (
	"sync"

	"github.com/Azure/azure-k8s-metrics-adapter/pkg/azure/monitor"
	"github.com/golang/glog"
)

// MetricCache holds the loaded metric request info in the system
type MetricCache struct {
	metricMutext sync.RWMutex
	metrics      map[string]monitor.AzureMetricRequest
}

// NewMetricCache creates the cache
func NewMetricCache() *MetricCache {
	return &MetricCache{
		metrics: make(map[string]monitor.AzureMetricRequest),
	}
}

// Update sets a metric request in the cache
func (mc *MetricCache) Update(key string, metricRequest monitor.AzureMetricRequest) {
	mc.metricMutext.Lock()
	defer mc.metricMutext.Unlock()

	mc.metrics[key] = metricRequest
}

// Get retrieves a metric request from the cache
func (mc *MetricCache) Get(key string) (monitor.AzureMetricRequest, bool) {
	mc.metricMutext.RLock()
	defer mc.metricMutext.RUnlock()

	metricRequest, exists := mc.metrics[key]
	if !exists {
		glog.V(2).Infof("metric not found %s", key)
		return monitor.AzureMetricRequest{}, false
	}

	return metricRequest, true
}

// Remove retrieves a metric request from the cache
func (mc *MetricCache) Remove(key string) {
	mc.metricMutext.Lock()
	defer mc.metricMutext.Unlock()

	delete(mc.metrics, key)
}

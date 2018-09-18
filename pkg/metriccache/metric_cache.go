package metriccache

import (
	"sync"

	"github.com/Azure/azure-k8s-metrics-adapter/pkg/azmetricrequest"
	"github.com/golang/glog"
)

// MetricCache holds the loaded metric request info in the system
type MetricCache struct {
	metricMutext sync.RWMutex
	metrics      map[string]azmetricrequest.AzureMetricRequest
}

// NewMetricCache creates the cache
func NewMetricCache() *MetricCache {
	return &MetricCache{
		metrics: make(map[string]azmetricrequest.AzureMetricRequest),
	}
}

// UpdateMetric sets a metric request in the cache
func (mc *MetricCache) UpdateMetric(key string, metricRequest azmetricrequest.AzureMetricRequest) {
	mc.metricMutext.Lock()
	defer mc.metricMutext.Unlock()

	mc.metrics[key] = metricRequest
}

// GetMetric retrieves a metric request from the cache
func (mc *MetricCache) GetMetric(key string) (azmetricrequest.AzureMetricRequest, bool) {
	mc.metricMutext.RLock()
	defer mc.metricMutext.RUnlock()

	metricRequest, exists := mc.metrics[key]
	if !exists {
		glog.V(2).Infof("metric not found %s", key)
		return azmetricrequest.AzureMetricRequest{}, false
	}

	return metricRequest, true
}

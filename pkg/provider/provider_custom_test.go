package provider

import (
	"errors"
	"testing"
	"time"

	"github.com/Azure/azure-k8s-metrics-adapter/pkg/azure/appinsights"

	"github.com/Azure/azure-k8s-metrics-adapter/pkg/metriccache"
	"github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/dynamicmapper"
	k8sprovider "github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/provider"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"

	core "k8s.io/client-go/testing"
)

func TestReturnsCustomMetricConverted(t *testing.T) {

	fakeClient := fakeAppInsightsClient{
		result: 15,
		err:    nil,
	}

	selector, _ := labels.Parse("")
	info := k8sprovider.CustomMetricInfo{
		Metric: "Metric-Name",
		GroupResource: schema.GroupResource{
			Resource: "pods",
		},
	}

	provider, _ := newFakeCustomProvider(fakeClient)
	returnList, err := provider.GetMetricBySelector("default", selector, info)

	if err != nil {
		t.Errorf("error after processing got: %v, want nil", err)
	}

	if len(returnList.Items) != 1 {
		t.Errorf("returnList.Items length = %v, want there 1", len(returnList.Items))
	}

	customMetric := returnList.Items[0]
	if customMetric.MetricName != "Metric/Name" {
		t.Errorf("customMetric.MetricName = %v, want there %v", customMetric.MetricName, "Metric/Name")
	}

	if customMetric.Value.MilliValue() != int64(15000) {
		t.Errorf("customMetric.Value.MilliValue() = %v, want there %v", customMetric.Value.MilliValue(), int64(15000))
	}
}

func TestReturnsCustomMetricWhenInCache(t *testing.T) {

	fakeClient := fakeAppInsightsClient{
		result: 15,
		err:    nil,
	}

	selector, _ := labels.Parse("")
	info := k8sprovider.CustomMetricInfo{
		Metric: "MetricName",
		GroupResource: schema.GroupResource{
			Resource: "pods",
		},
	}

	provider, cache := newFakeCustomProvider(fakeClient)

	request := appinsights.MetricRequest{
		MetricName: "cachedName",
	}

	cache.Update("CustomMetric/default/MetricName", request)

	returnList, err := provider.GetMetricBySelector("default", selector, info)

	if err != nil {
		t.Errorf("error after processing got: %v, want nil", err)
	}

	if len(returnList.Items) != 1 {
		t.Errorf("returnList.Items length = %v, want there 1", len(returnList.Items))
	}

	customMetric := returnList.Items[0]
	if customMetric.MetricName != request.MetricName {
		t.Errorf("customMetric.MetricName = %v, want there %v", customMetric.MetricName, request.MetricName)
	}

	if customMetric.Value.MilliValue() != int64(15000) {
		t.Errorf("customMetric.Value.MilliValue() = %v, want there %v", customMetric.Value.MilliValue(), int64(15000))
	}
}

func TestReturnsErrorIfAppInsightsFails(t *testing.T) {

	fakeClient := fakeAppInsightsClient{
		err: errors.New("force error for test"),
	}

	selector, _ := labels.Parse("")
	info := k8sprovider.CustomMetricInfo{
		Metric: "MetricName",
		GroupResource: schema.GroupResource{
			Resource: "pods",
		},
	}

	provider, _ := newFakeCustomProvider(fakeClient)
	_, err := provider.GetMetricBySelector("default", selector, info)

	if !k8serrors.IsBadRequest(err) {
		t.Errorf("error after processing got: %v, want an bad request error", err)
	}
}

func newFakeCustomProvider(fakeclient fakeAppInsightsClient) (AzureProvider, *metriccache.MetricCache) {
	metricCache := metriccache.NewMetricCache()

	// set up a fake mapper
	fakeDiscovery := &dynamicmapper.FakeDiscovery{Fake: &core.Fake{}}
	mapper, _ := dynamicmapper.NewRESTMapper(fakeDiscovery, 1*time.Second)

	fakeDiscovery.Resources = []*metav1.APIResourceList{
		{
			GroupVersion: "v1",
			APIResources: []metav1.APIResource{
				{Name: "pods", Namespaced: true, Kind: "Pod"},
			},
		},
	}

	mapper.RegenerateMappings()

	provider := AzureProvider{
		metricCache:       metricCache,
		appinsightsClient: fakeclient,
		mapper:            mapper,
	}

	return provider, metricCache
}

type fakeAppInsightsClient struct {
	result float64
	err    error
}

func (f fakeAppInsightsClient) GetCustomMetric(request appinsights.MetricRequest) (float64, error) {
	return f.result, f.err
}

package provider

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/Azure/azure-k8s-metrics-adapter/pkg/azure/custommetrics"
	"github.com/Azure/azure-k8s-metrics-adapter/pkg/client/clientset/versioned/scheme"

	"github.com/Azure/azure-k8s-metrics-adapter/pkg/metriccache"
	"github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/dynamicmapper"
	k8sprovider "github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/provider"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	k8sclient "k8s.io/client-go/dynamic/fake"

	core "k8s.io/client-go/testing"
)

func TestReturnsCustomMetricConverted(t *testing.T) {

	fakeClient := fakeAppInsightsClient{
		result: 15,
		err:    nil,
	}

	selector, _ := labels.Parse("")
	info := k8sprovider.CustomMetricInfo{
		Namespaced: true,
		Metric:     "Metric-Name",
		GroupResource: schema.GroupResource{
			Resource: "pods",
		},
	}

	var storeObjects []runtime.Object
	pod := newUnstructured("v1", "Pod", "default", "pod1")
	storeObjects = append(storeObjects, pod)

	provider, _ := newFakeCustomProvider(fakeClient, storeObjects)
	returnList, err := provider.GetMetricBySelector("default", selector, info)

	if err != nil {
		t.Errorf("error after processing got: %v, want nil", err)
	}

	if len(returnList.Items) != 1 {
		t.Errorf("returnList.Items length = %v, want there 1", len(returnList.Items))
	}

	customMetric := returnList.Items[0]
	if customMetric.Metric.Name != "Metric-Name" {
		t.Errorf("customMetric.Metric.Name = %v, want there %v", customMetric.Metric.Name, "Metric/Name")
	}

	if customMetric.Value.MilliValue() != int64(15000) {
		t.Errorf("customMetric.Value.MilliValue() = %v, want there %v", customMetric.Value.MilliValue(), int64(15000))
	}
}

func TestReturnsCustomMetricConvertedWithMultiplePods(t *testing.T) {
	fakeClient := fakeAppInsightsClient{
		result: 15,
		err:    nil,
	}

	selector, _ := labels.Parse("")
	info := k8sprovider.CustomMetricInfo{
		Namespaced: true,
		Metric:     "Metric-Name",
		GroupResource: schema.GroupResource{
			Resource: "pods",
		},
	}

	var storeObjects []runtime.Object
	pod := newUnstructured("v1", "Pod", "default", "pod0")
	pod2 := newUnstructured("v1", "Pod", "default", "pod1")
	pod3 := newUnstructured("v1", "Pod", "default", "pod2")
	storeObjects = append(storeObjects, pod, pod2, pod3)

	provider, _ := newFakeCustomProvider(fakeClient, storeObjects)
	returnList, err := provider.GetMetricBySelector("default", selector, info)

	if err != nil {
		t.Errorf("error after processing got: %v, want nil", err)
	}

	if len(returnList.Items) != 3 {
		t.Errorf("returnList.Items length = %v, want there 3", len(returnList.Items))
	}

	for i, customMetric := range returnList.Items {
		if customMetric.Metric.Name != "Metric-Name" {
			t.Errorf("customMetric.Metric.Name = %v, want there %v", customMetric.Metric.Name, "Metric/Name")
		}

		if customMetric.Value.MilliValue() != int64(15000) {
			t.Errorf("customMetric.Value.MilliValue() = %v, want there %v", customMetric.Value.MilliValue(), int64(15000))
		}

		if customMetric.DescribedObject.Name != fmt.Sprintf("pod%d", i) {
			t.Errorf("customMetric.Value.MilliValue() = %v, want there %v", customMetric.Value.MilliValue(), int64(15000))
		}
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

	var storeObjects []runtime.Object
	pod := newUnstructured("v1", "Pod", "default", "pod1")
	storeObjects = append(storeObjects, pod)

	provider, cache := newFakeCustomProvider(fakeClient, storeObjects)

	request := custommetrics.MetricRequest{
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
	if customMetric.Metric.Name != "MetricName" {
		t.Errorf("customMetric.Metric.Name = %v, want there %v", customMetric.Metric.Name, request.MetricName)
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

	var storeObjects []runtime.Object
	pod := newUnstructured("v1", "Pod", "default", "pod1")
	storeObjects = append(storeObjects, pod)

	provider, _ := newFakeCustomProvider(fakeClient, storeObjects)
	_, err := provider.GetMetricBySelector("default", selector, info)

	if !k8serrors.IsBadRequest(err) {
		t.Errorf("error after processing got: %v, want an bad request error", err)
	}
}

func newFakeCustomProvider(fakeclient fakeAppInsightsClient, store []runtime.Object) (AzureProvider, *metriccache.MetricCache) {
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

	// set up fake dynamic client
	s := scheme.Scheme
	corev1.SchemeBuilder.AddToScheme(s)

	fakeK8sClient := k8sclient.NewSimpleDynamicClient(s, store...)

	provider := AzureProvider{
		metricCache:       metricCache,
		appinsightsClient: fakeclient,
		mapper:            mapper,
		kubeClient:        fakeK8sClient,
	}

	return provider, metricCache
}

type fakeAppInsightsClient struct {
	result float64
	err    error
}

func (f fakeAppInsightsClient) GetCustomMetric(request custommetrics.MetricRequest) (float64, error) {
	return f.result, f.err
}

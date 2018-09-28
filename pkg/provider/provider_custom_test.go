package provider

import (
	"errors"
	"testing"
	"time"

	"github.com/Azure/azure-k8s-metrics-adapter/pkg/metriccache"
	"github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/dynamicmapper"
	k8sprovider "github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/provider"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"

	core "k8s.io/client-go/testing"
)

func TestReturnsCustomMetric(t *testing.T) {

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

	provider := newFakeCustomProvider(fakeClient)
	returnList, err := provider.GetMetricBySelector("default", selector, info)

	if err != nil {
		t.Errorf("error after processing got: %v, want nil", err)
	}

	if len(returnList.Items) != 1 {
		t.Errorf("returnList.Items length = %v, want there 1", len(returnList.Items))
	}

	customMetric := returnList.Items[0]
	if customMetric.MetricName != info.Metric {
		t.Errorf("externalMetric.MetricName = %v, want there %v", customMetric.MetricName, info.Metric)
	}

	if customMetric.Value.MilliValue() != int64(15000) {
		t.Errorf("externalMetric.Value.MilliValue() = %v, want there %v", customMetric.Value.MilliValue(), int64(15000))
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

	provider := newFakeCustomProvider(fakeClient)
	_, err := provider.GetMetricBySelector("default", selector, info)

	if !k8serrors.IsBadRequest(err) {
		t.Errorf("error after processing got: %v, want an bad request error", err)
	}
}

func newFakeCustomProvider(fakeclient fakeAppInsightsClient) AzureProvider {
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

	return provider
}

type fakeAppInsightsClient struct {
	result float64
	err    error
}

func (f fakeAppInsightsClient) GetCustomMetric(namespace string, metricName string) (float64, error) {
	return f.result, f.err
}

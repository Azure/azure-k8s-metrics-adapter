package controller

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-k8s-metrics-adapter/pkg/azmetricrequest"

	api "github.com/Azure/azure-k8s-metrics-adapter/pkg/apis/metrics/v1alpha1"
	"github.com/Azure/azure-k8s-metrics-adapter/pkg/metriccache"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/Azure/azure-k8s-metrics-adapter/pkg/client/clientset/versioned/fake"
	informers "github.com/Azure/azure-k8s-metrics-adapter/pkg/client/informers/externalversions"
)

func getKey(externalMetric *api.ExternalMetric) string {
	return fmt.Sprintf("%s/%s", externalMetric.Namespace, externalMetric.Name)
}

func TestMetricValueIsStored(t *testing.T) {
	var storeObjects []runtime.Object
	var externalMetricsListerCache []*api.ExternalMetric

	externalMetric := newFullExternalMetric("test")
	storeObjects = append(storeObjects, externalMetric)
	externalMetricsListerCache = append(externalMetricsListerCache, externalMetric)

	handler, metriccache := newHandler(storeObjects, externalMetricsListerCache)

	key := getKey(externalMetric)
	err := handler.Process(key)

	if err != nil {
		t.Errorf("error after processing = %v, want %v", err, nil)
	}

	metricRequest, exists := metriccache.Get(key)

	if exists == false {
		t.Errorf("exist = %v, want %v", exists, true)
	}

	validateMetricResult(metricRequest, externalMetric, t)
}

func TestShouldFailOnInvalidCacheKey(t *testing.T) {
	var storeObjects []runtime.Object
	var externalMetricsListerCache []*api.ExternalMetric

	externalMetric := newFullExternalMetric("test")
	storeObjects = append(storeObjects, externalMetric)
	externalMetricsListerCache = append(externalMetricsListerCache, externalMetric)

	handler, metriccache := newHandler(storeObjects, externalMetricsListerCache)

	key := "invalidkey/with/extrainfo"
	err := handler.Process(key)

	if err == nil {
		t.Errorf("error after processing nil, want non nil")
	}

	_, exists := metriccache.Get(key)

	if exists == true {
		t.Errorf("exist = %v, want %v", exists, false)
	}
}

func TestWhenItemHasBeenDeleted(t *testing.T) {
	var storeObjects []runtime.Object
	var externalMetricsListerCache []*api.ExternalMetric

	externalMetric := newFullExternalMetric("test")

	// don't put anything in the stores
	handler, metriccache := newHandler(storeObjects, externalMetricsListerCache)

	// add the item to the cache then test if it gets deleted
	key := getKey(externalMetric)
	metriccache.Update(key, azmetricrequest.AzureMetricRequest{})

	err := handler.Process(key)

	if err != nil {
		t.Errorf("error == %v, want nil", err)
	}

	_, exists := metriccache.Get(key)

	if exists == true {
		t.Errorf("exist = %v, want %v", exists, false)
	}
}

func newHandler(storeObjects []runtime.Object, externalMetricsListerCache []*api.ExternalMetric) (Handler, *metriccache.MetricCache) {
	fakeClient := fake.NewSimpleClientset(storeObjects...)
	i := informers.NewSharedInformerFactory(fakeClient, 0)

	lister := i.Azure().V1alpha1().ExternalMetrics().Lister()

	for _, em := range externalMetricsListerCache {
		i.Azure().V1alpha1().ExternalMetrics().Informer().GetIndexer().Add(em)
	}

	metriccache := metriccache.NewMetricCache()
	handler := NewHandler(lister, metriccache)

	return handler, metriccache
}

func validateMetricResult(metricRequest azmetricrequest.AzureMetricRequest, externalMetricInfo *api.ExternalMetric, t *testing.T) {

	// Metric Config
	if metricRequest.MetricName != externalMetricInfo.Spec.MetricConfig.MetricName {
		t.Errorf("metricRequest MetricName = %v, want %v", metricRequest.MetricName, externalMetricInfo.Spec.MetricConfig.MetricName)
	}

	if metricRequest.Filter != externalMetricInfo.Spec.MetricConfig.Filter {
		t.Errorf("metricRequest Filter = %v, want %v", metricRequest.Filter, externalMetricInfo.Spec.MetricConfig.Filter)
	}

	if metricRequest.Aggregation != externalMetricInfo.Spec.MetricConfig.Aggregation {
		t.Errorf("metricRequest Aggregation = %v, want %v", metricRequest.Aggregation, externalMetricInfo.Spec.MetricConfig.Aggregation)
	}

	// Azure Config
	if metricRequest.ResourceGroup != externalMetricInfo.Spec.AzureConfig.ResourceGroup {
		t.Errorf("metricRequest ResourceGroup = %v, want %v", metricRequest.ResourceGroup, externalMetricInfo.Spec.AzureConfig.ResourceGroup)
	}

	if metricRequest.ResourceName != externalMetricInfo.Spec.AzureConfig.ResourceName {
		t.Errorf("metricRequest ResourceName = %v, want %v", metricRequest.ResourceName, externalMetricInfo.Spec.AzureConfig.ResourceName)
	}

	if metricRequest.ResourceProviderNamespace != externalMetricInfo.Spec.AzureConfig.ResourceProviderNamespace {
		t.Errorf("metricRequest ResourceProviderNamespace = %v, want %v", metricRequest.ResourceProviderNamespace, externalMetricInfo.Spec.AzureConfig.ResourceProviderNamespace)
	}

	if metricRequest.ResourceType != externalMetricInfo.Spec.AzureConfig.ResourceType {
		t.Errorf("metricRequest ResourceType = %v, want %v", metricRequest.ResourceType, externalMetricInfo.Spec.AzureConfig.ResourceType)
	}

	if metricRequest.SubscriptionID != externalMetricInfo.Spec.AzureConfig.SubscriptionID {
		t.Errorf("metricRequest SubscriptionID = %v, want %v", metricRequest.SubscriptionID, externalMetricInfo.Spec.AzureConfig.SubscriptionID)
	}

}

func newFullExternalMetric(name string) *api.ExternalMetric {
	// must preserve upper casing for azure api
	return &api.ExternalMetric{
		TypeMeta: metav1.TypeMeta{APIVersion: api.SchemeGroupVersion.String()},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: metav1.NamespaceDefault,
		},
		Spec: api.ExternalMetricSpec{
			AzureConfig: api.AzureConfig{
				ResourceGroup:             "rg",
				ResourceName:              "rn",
				ResourceProviderNamespace: "Resource.NameSpace",
				ResourceType:              "rt",
			},
			MetricConfig: api.MetricConfig{
				Aggregation: "Total",
				MetricName:  "Name",
				Filter:      "EntityName eq 'externalq'",
			},
		},
	}
}

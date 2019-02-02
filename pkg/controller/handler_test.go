package controller

import (
	"fmt"
	"testing"

	api "github.com/Azure/azure-k8s-metrics-adapter/pkg/apis/metrics/v1alpha2"
	"github.com/Azure/azure-k8s-metrics-adapter/pkg/azure/custommetrics"
	"github.com/Azure/azure-k8s-metrics-adapter/pkg/azure/externalmetrics"
	"github.com/Azure/azure-k8s-metrics-adapter/pkg/client/clientset/versioned/fake"
	informers "github.com/Azure/azure-k8s-metrics-adapter/pkg/client/informers/externalversions"
	"github.com/Azure/azure-k8s-metrics-adapter/pkg/metriccache"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func getExternalKey(externalMetric *api.ExternalMetric) namespacedQueueItem {
	return namespacedQueueItem{
		namespaceKey: fmt.Sprintf("%s/%s", externalMetric.Namespace, externalMetric.Name),
		kind:         externalMetric.TypeMeta.Kind,
	}
}

func getCustomKey(customMetric *api.CustomMetric) namespacedQueueItem {
	return namespacedQueueItem{
		namespaceKey: fmt.Sprintf("%s/%s", customMetric.Namespace, customMetric.Name),
		kind:         customMetric.TypeMeta.Kind,
	}
}

func TestExternalMetricValueIsStored(t *testing.T) {
	var storeObjects []runtime.Object
	var externalMetricsListerCache []*api.ExternalMetric
	var customMetricsListerCache []*api.CustomMetric

	externalMetric := newFullExternalMetric("test")
	storeObjects = append(storeObjects, externalMetric)
	externalMetricsListerCache = append(externalMetricsListerCache, externalMetric)

	handler, metriccache := newHandler(storeObjects, externalMetricsListerCache, customMetricsListerCache)

	queueItem := getExternalKey(externalMetric)
	err := handler.Process(queueItem)

	if err != nil {
		t.Errorf("error after processing = %v, want %v", err, nil)
	}

	metricRequest, exists := metriccache.GetAzureExternalMetricRequest(externalMetric.Namespace, externalMetric.Name)

	if exists == false {
		t.Errorf("exist = %v, want %v", exists, true)
	}

	validateExternalMetricResult(metricRequest, externalMetric, t)
}

func TestShouldBeAbleToStoreCustomAndExternalWithSameNameAndNamespace(t *testing.T) {
	var storeObjects []runtime.Object
	var externalMetricsListerCache []*api.ExternalMetric
	var customMetricsListerCache []*api.CustomMetric

	externalMetric := newFullExternalMetric("test")
	customMetric := newFullCustomMetric("test")
	storeObjects = append(storeObjects, externalMetric, customMetric)
	externalMetricsListerCache = append(externalMetricsListerCache, externalMetric)
	customMetricsListerCache = append(customMetricsListerCache, customMetric)

	handler, metriccache := newHandler(storeObjects, externalMetricsListerCache, customMetricsListerCache)

	externalItem := getExternalKey(externalMetric)
	err := handler.Process(externalItem)

	if err != nil {
		t.Errorf("error after processing = %v, want %v", err, nil)
	}

	customItem := getCustomKey(customMetric)
	err = handler.Process(customItem)

	if err != nil {
		t.Errorf("error after processing = %v, want %v", err, nil)
	}

	externalRequest, exists := metriccache.GetAzureExternalMetricRequest(externalMetric.Namespace, externalMetric.Name)

	if exists == false {
		t.Errorf("exist = %v, want %v", exists, true)
	}

	validateExternalMetricResult(externalRequest, externalMetric, t)

	metricRequest, exists := metriccache.GetAppInsightsRequest(customMetric.Namespace, customMetric.Name)

	if exists == false {
		t.Errorf("exist = %v, want %v", exists, true)
	}

	validateCustomMetricResult(metricRequest, customMetric, t)
}

func TestCustomMetricValueIsStored(t *testing.T) {
	var storeObjects []runtime.Object
	var externalMetricsListerCache []*api.ExternalMetric
	var customMetricsListerCache []*api.CustomMetric

	customMetric := newFullCustomMetric("test")
	storeObjects = append(storeObjects, customMetric)
	customMetricsListerCache = append(customMetricsListerCache, customMetric)

	handler, metriccache := newHandler(storeObjects, externalMetricsListerCache, customMetricsListerCache)

	queueItem := getCustomKey(customMetric)
	err := handler.Process(queueItem)

	if err != nil {
		t.Errorf("error after processing = %v, want %v", err, nil)
	}

	metricRequest, exists := metriccache.GetAppInsightsRequest(customMetric.Namespace, customMetric.Name)

	if exists == false {
		t.Errorf("exist = %v, want %v", exists, true)
	}

	validateCustomMetricResult(metricRequest, customMetric, t)
}

func TestShouldFailOnInvalidCacheKey(t *testing.T) {
	var storeObjects []runtime.Object
	var externalMetricsListerCache []*api.ExternalMetric
	var customMetricsListerCache []*api.CustomMetric

	externalMetric := newFullExternalMetric("test")
	storeObjects = append(storeObjects, externalMetric)
	externalMetricsListerCache = append(externalMetricsListerCache, externalMetric)

	handler, metriccache := newHandler(storeObjects, externalMetricsListerCache, customMetricsListerCache)

	queueItem := namespacedQueueItem{
		namespaceKey: "invalidkey/with/extrainfo",
		kind:         "somethingwrong",
	}
	err := handler.Process(queueItem)

	if err == nil {
		t.Errorf("error after processing nil, want non nil")
	}

	_, exists := metriccache.GetAzureExternalMetricRequest(externalMetric.Namespace, externalMetric.Name)

	if exists == true {
		t.Errorf("exist = %v, want %v", exists, false)
	}
}

func TestWhenExternalItemHasBeenDeleted(t *testing.T) {
	var storeObjects []runtime.Object
	var externalMetricsListerCache []*api.ExternalMetric
	var customMetricsListerCache []*api.CustomMetric

	externalMetric := newFullExternalMetric("test")

	// don't put anything in the stores
	handler, metriccache := newHandler(storeObjects, externalMetricsListerCache, customMetricsListerCache)

	// add the item to the cache then test if it gets deleted
	queueItem := getExternalKey(externalMetric)
	metriccache.Update(queueItem.Key(), externalmetrics.AzureExternalMetricRequest{})

	err := handler.Process(queueItem)

	if err != nil {
		t.Errorf("error == %v, want nil", err)
	}

	_, exists := metriccache.GetAzureExternalMetricRequest(externalMetric.Namespace, externalMetric.Name)

	if exists == true {
		t.Errorf("exist = %v, want %v", exists, false)
	}
}

func TestWhenCustomItemHasBeenDeleted(t *testing.T) {
	var storeObjects []runtime.Object
	var externalMetricsListerCache []*api.ExternalMetric
	var customMetricsListerCache []*api.CustomMetric

	customMetric := newFullCustomMetric("test")

	// don't put anything in the stores
	handler, metriccache := newHandler(storeObjects, externalMetricsListerCache, customMetricsListerCache)

	// add the item to the cache then test if it gets deleted
	queueItem := getCustomKey(customMetric)
	metriccache.Update(queueItem.Key(), custommetrics.MetricRequest{})

	err := handler.Process(queueItem)

	if err != nil {
		t.Errorf("error == %v, want nil", err)
	}

	_, exists := metriccache.GetAppInsightsRequest(customMetric.Namespace, customMetric.Name)

	if exists == true {
		t.Errorf("exist = %v, want %v", exists, false)
	}
}

func TestWhenItemKindIsUnknown(t *testing.T) {
	var storeObjects []runtime.Object
	var externalMetricsListerCache []*api.ExternalMetric
	var customMetricsListerCache []*api.CustomMetric

	// don't put anything in the stores, as we are not looking anything up
	handler, metriccache := newHandler(storeObjects, externalMetricsListerCache, customMetricsListerCache)

	// add the item to the cache then test if it gets deleted
	queueItem := namespacedQueueItem{
		namespaceKey: "default/unknown",
		kind:         "Unknown",
	}

	err := handler.Process(queueItem)

	if err != nil {
		t.Errorf("error == %v, want nil", err)
	}

	_, exists := metriccache.GetAppInsightsRequest("default", "unkown")

	if exists == true {
		t.Errorf("exist = %v, want %v", exists, false)
	}
}

func newHandler(storeObjects []runtime.Object, externalMetricsListerCache []*api.ExternalMetric, customMetricsListerCache []*api.CustomMetric) (Handler, *metriccache.MetricCache) {
	fakeClient := fake.NewSimpleClientset(storeObjects...)
	i := informers.NewSharedInformerFactory(fakeClient, 0)

	externalMetricLister := i.Azure().V1alpha2().ExternalMetrics().Lister()
	customMetricLister := i.Azure().V1alpha2().CustomMetrics().Lister()

	for _, em := range externalMetricsListerCache {
		i.Azure().V1alpha2().ExternalMetrics().Informer().GetIndexer().Add(em)
	}

	for _, cm := range customMetricsListerCache {
		i.Azure().V1alpha2().CustomMetrics().Informer().GetIndexer().Add(cm)
	}

	metriccache := metriccache.NewMetricCache()
	handler := NewHandler(externalMetricLister, customMetricLister, metriccache)

	return handler, metriccache
}

func validateExternalMetricResult(metricRequest externalmetrics.AzureExternalMetricRequest, externalMetricInfo *api.ExternalMetric, t *testing.T) {

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

func validateCustomMetricResult(metricRequest custommetrics.MetricRequest, customMetricInfo *api.CustomMetric, t *testing.T) {
	// Metric Config
	if metricRequest.MetricName != customMetricInfo.Spec.MetricConfig.MetricName {
		t.Errorf("metricRequest MetricName = %v, want %v", metricRequest.MetricName, customMetricInfo.Spec.MetricConfig.MetricName)
	}

}

func newFullExternalMetric(name string) *api.ExternalMetric {
	// must preserve upper casing for azure api
	return &api.ExternalMetric{
		TypeMeta: metav1.TypeMeta{APIVersion: api.SchemeGroupVersion.String(), Kind: "ExternalMetric"},
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
			MetricConfig: api.ExternalMetricConfig{
				Aggregation: "Total",
				MetricName:  "Name",
				Filter:      "EntityName eq 'externalq'",
			},
		},
	}
}

func newFullCustomMetric(name string) *api.CustomMetric {
	// must preserve upper casing for azure api
	return &api.CustomMetric{
		TypeMeta: metav1.TypeMeta{APIVersion: api.SchemeGroupVersion.String(), Kind: "CustomMetric"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: metav1.NamespaceDefault,
		},
		Spec: api.CustomMetricSpec{
			MetricConfig: api.CustomMetricConfig{
				MetricName: "performance/requestpersecond",
			},
		},
	}
}

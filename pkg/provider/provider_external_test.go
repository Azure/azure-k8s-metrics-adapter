package provider

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-k8s-metrics-adapter/pkg/azure/externalmetrics"
	"github.com/Azure/azure-k8s-metrics-adapter/pkg/metriccache"
	k8sprovider "github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/provider"
	"k8s.io/apimachinery/pkg/labels"
)

const validLabelSelector = "resourceProviderNamespace=Microsoft.Servicebus,resourceType=namespaces,aggregation=Total,filter=EntityName_eq_externalq,resourceGroup=sb-external-example,resourceName=sb-external-ns,metricName=%s"

func createLabelSelector(metricName, subscriptionID string) labels.Selector {
	labelSelector := validLabelSelector

	if subscriptionID != "" {
		labelSelector = fmt.Sprintf("%s,subscriptionID=%s", labelSelector, subscriptionID)
	}

	label := fmt.Sprintf(labelSelector, metricName)
	selector, err := labels.Parse(label)

	if err != nil {
		panic("was not able to make label selector")
	}

	return selector
}

func TestFindMetricInCache(t *testing.T) {
	metricCache := metriccache.NewMetricCache()

	request := externalmetrics.AzureExternalMetricRequest{
		MetricName: "MessageCount",
	}
	metricCache.Update("ExternalMetric/default/metricname", request)

	provider := AzureProvider{
		metricCache:           metricCache,
		defaultSubscriptionID: "1234",
	}

	selector, _ := labels.Parse("")
	foundRequest, err := provider.getMetricRequest("default", "metricname", selector)

	if err != nil {
		t.Errorf("error after processing got: %v, want nil", err)
	}

	if foundRequest.MetricName != request.MetricName {
		t.Errorf("foundRequest.MetricName = %v, want %s", foundRequest.MetricName, request.MetricName)
	}

	if foundRequest.Timespan == "" {
		t.Errorf("foundRequest.TimeSpan = %v, want there to be value", foundRequest.Timespan)
	}

	if foundRequest.SubscriptionID != provider.defaultSubscriptionID {
		t.Errorf("foundRequest.SubscriptionID = %v, want %s", foundRequest.SubscriptionID, provider.defaultSubscriptionID)
	}
}

func TestFindMetricInCacheUsesOverrideSubscriptionId(t *testing.T) {
	metricCache := metriccache.NewMetricCache()

	request := externalmetrics.AzureExternalMetricRequest{
		MetricName:     "MessageCount",
		SubscriptionID: "9876",
	}
	metricCache.Update("ExternalMetric/default/metricname", request)

	provider := AzureProvider{
		metricCache:           metricCache,
		defaultSubscriptionID: "1234",
	}

	selector, _ := labels.Parse("")
	foundRequest, err := provider.getMetricRequest("default", "metricname", selector)

	if err != nil {
		t.Errorf("error after processing got: %v, want nil", err)
	}

	if foundRequest.MetricName != request.MetricName {
		t.Errorf("foundRequest.MetricName = %v, want %s", foundRequest.MetricName, request.MetricName)
	}

	if foundRequest.SubscriptionID != request.SubscriptionID {
		t.Errorf("foundRequest.SubscriptionID = %v, want %s", foundRequest.SubscriptionID, request.SubscriptionID)
	}
}

func TestNoMetricInCache(t *testing.T) {
	metricCache := metriccache.NewMetricCache()

	provider := AzureProvider{
		metricCache:           metricCache,
		defaultSubscriptionID: "1234",
	}

	metricName := "Messages"
	selector := createLabelSelector(metricName, "")
	foundRequest, err := provider.getMetricRequest("default", "metricname", selector)

	if err != nil {
		t.Errorf("error after processing got: %v, want nil", err)
	}

	if foundRequest.MetricName != metricName {
		t.Errorf("foundRequest = %v, want %s", foundRequest.MetricName, metricName)
	}

	if foundRequest.Timespan == "" {
		t.Errorf("foundRequest.TimeSpan = %v, want there to be value", foundRequest.MetricName)
	}

	if foundRequest.SubscriptionID != provider.defaultSubscriptionID {
		t.Errorf("foundRequest.SubscriptionID = %v, want %s", foundRequest.SubscriptionID, provider.defaultSubscriptionID)
	}
}

func TestNoMetricInCacheUsesOverrideSubscriptionID(t *testing.T) {
	metricCache := metriccache.NewMetricCache()

	provider := AzureProvider{
		metricCache:           metricCache,
		defaultSubscriptionID: "1234",
	}

	metricName := "Messages"
	overrideSubID := "9876"
	selector := createLabelSelector(metricName, overrideSubID)
	foundRequest, err := provider.getMetricRequest("default", "metricname", selector)

	if err != nil {
		t.Errorf("error after processing got: %v, want nil", err)
	}

	if foundRequest.MetricName != metricName {
		t.Errorf("foundRequest = %v, want %s", foundRequest.MetricName, metricName)
	}

	if foundRequest.Timespan == "" {
		t.Errorf("foundRequest.TimeSpan = %v, want there to be value", foundRequest.MetricName)
	}

	if foundRequest.SubscriptionID != overrideSubID {
		t.Errorf("foundRequest.SubscriptionID = %v, want %s", foundRequest.SubscriptionID, overrideSubID)
	}
}

func TestInvalidLabelSelector(t *testing.T) {
	metricCache := metriccache.NewMetricCache()

	provider := AzureProvider{
		metricCache: metricCache,
	}

	_, err := provider.getMetricRequest("default", "metricname", nil)

	if err == nil {
		t.Errorf("no error after processing got: %v, want error", nil)
	}
}

func TestReturnsExeternalMetric(t *testing.T) {
	fakeFactory := fakeAzureExternalClientFactory{}

	selector, _ := labels.Parse("")
	info := k8sprovider.ExternalMetricInfo{
		Metric: "MetricName",
	}

	provider := newProvider(fakeFactory)
	returnList, err := provider.GetExternalMetric("default", selector, info)

	if err != nil {
		t.Errorf("error after processing got: %v, want nil", err)
	}

	if len(returnList.Items) != 1 {
		t.Errorf("returnList.Items length = %v, want there 1", len(returnList.Items))
	}

	externalMetric := returnList.Items[0]
	if externalMetric.MetricName != info.Metric {
		t.Errorf("externalMetric.MetricName = %v, want there %v", externalMetric.MetricName, info.Metric)
	}

	if externalMetric.Value.MilliValue() != int64(15000) {
		t.Errorf("externalMetric.Value.MilliValue() = %v, want there %v", externalMetric.Value.MilliValue(), int64(15000))
	}
}

func newProvider(fakeFactory fakeAzureExternalClientFactory) AzureProvider {
	// func newProvider(fakeclient fakeAzureMonitorClient) AzureProvider {
	metricCache := metriccache.NewMetricCache()

	provider := AzureProvider{
		metricCache:        metricCache,
		azureClientFactory: fakeFactory,
	}

	return provider
}

// externalMetricClient, err := p.azureExternalClientFactory.GetAzureExternalMetricClient(azMetricRequest.Type)

type fakeAzureExternalClientFactory struct {
}

func (f fakeAzureExternalClientFactory) GetAzureExternalMetricClient(clientType string) (client externalmetrics.AzureExternalMetricClient, err error) {
	fakeClient := fakeAzureMonitorClient{
		err:    nil,
		result: externalmetrics.AzureExternalMetricResponse{Total: 15},
	}

	return fakeClient, nil
}

type fakeAzureMonitorClient struct {
	result externalmetrics.AzureExternalMetricResponse
	err    error
}

func (f fakeAzureMonitorClient) GetAzureMetric(azMetricRequest externalmetrics.AzureExternalMetricRequest) (externalmetrics.AzureExternalMetricResponse, error) {
	return f.result, f.err
}

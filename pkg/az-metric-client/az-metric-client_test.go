package azureMetricClient

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2018-03-01/insights"

	"k8s.io/apimachinery/pkg/labels"

	"github.com/Azure/azure-k8s-metrics-adapter/pkg/azmetricrequest"

	"github.com/Azure/azure-k8s-metrics-adapter/pkg/metriccache"
)

func Test_normalizeValue(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "int64 to float64",
			args: args{
				value: int64(42),
			},
			want: float64(42),
		},
		{
			name: "float64 to float64",
			args: args{
				value: float64(42.0),
			},
			want: float64(42),
		},
		{
			name: "int32 to float64",
			args: args{
				value: int32(42),
			},
			want: float64(42),
		},
		{
			name: "float32 to float64",
			args: args{
				value: float32(42.0),
			},
			want: float64(42),
		},
		{
			name: "if something random like a string, return 0",
			args: args{
				value: "this is not the answer",
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeValue(tt.args.value); got != tt.want {
				t.Errorf("normalizeValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func newFakeMonitorClient(result insights.Response, err error) MonitorClient {
	return fakeMonitorClient{
		err:    err,
		result: result,
	}
}

type fakeMonitorClient struct {
	result insights.Response
	err    error
}

func (f fakeMonitorClient) List(ctx context.Context, resourceURI string, timespan string, interval *string, metricnames string, aggregation string, top *int32, orderby string, filter string, resultType insights.ResultType, metricnamespace string) (result insights.Response, err error) {
	result = f.result
	err = f.err
	return
}

func makeResponse(value float64) insights.Response {
	metricValues := []insights.MetricValue{}
	metrics := []insights.Metric{}
	timeseries := []insights.TimeSeriesElement{}

	te := insights.TimeSeriesElement{
		Data: &metricValues,
	}

	v := value
	mv := insights.MetricValue{
		Total: &v,
	}

	metric := insights.Metric{
		Timeseries: &timeseries,
	}

	metrics = append(metrics, metric)
	timeseries = append(timeseries, te)
	metricValues = append(metricValues, mv)

	response := insights.Response{
		Value: &metrics,
	}
	return response
}

func Test_IfNotValidCallGetError(t *testing.T) {
	monitorClient := newFakeMonitorClient(insights.Response{}, nil)
	metricCache := metriccache.NewMetricCache()

	client := NewAzureMetricClient("", metricCache, monitorClient)

	_, err := client.GetAzureMetric("ns", "name", nil)

	if err == nil {
		t.Errorf("error after processing got nil, want non nil")
	}
}

func Test_IfInsufficientDataGetError(t *testing.T) {
	monitorClient := newFakeMonitorClient(insights.Response{}, nil)
	metricCache := metriccache.NewMetricCache()

	client := NewAzureMetricClient("", metricCache, monitorClient)

	// This doesn't have the all requeired selectors so should report that it is missing
	selector, _ := labels.Parse("resourceProviderNamespace=Microsoft.Servicebus")
	_, err := client.GetAzureMetric("ns", "name", selector)

	if !azmetricrequest.IsInvalidMetricRequestError(err) {
		t.Errorf("should be InvalidMetricRequest error got %v, want InvalidMetricRequestError", err)
	}
}

func Test_IfPassedViaLabelSelectorsItReturns(t *testing.T) {
	response := makeResponse(15)

	monitorClient := newFakeMonitorClient(response, nil)
	metricCache := metriccache.NewMetricCache()

	client := NewAzureMetricClient("1234", metricCache, monitorClient)

	selector, _ := labels.Parse("resourceProviderNamespace=Microsoft.Servicebus,resourceType=namespaces,aggregation=Total,filter=EntityName_eq_externalq,resourceGroup=sb-external-example,resourceName=sb-external-ns,metricName=Messages")

	metricValue, err := client.GetAzureMetric("ns", "name", selector)

	if err != nil {
		t.Errorf("error after processing got %v, want nil", err)
	}

	if metricValue.MetricName != "Messages" {
		t.Errorf("error after processing got %v, want nil", err)
	}

	valueReturned := metricValue.Value.MilliValue()
	if valueReturned != int64(15000) {
		t.Errorf("MilliValue() got %v, want 15000", valueReturned)
	}
}

func Test_IfCacheHasItReturn(t *testing.T) {
	response := makeResponse(15)

	monitorClient := newFakeMonitorClient(response, nil)
	metricCache := metriccache.NewMetricCache()

	metricRequest := newMetricRequest()
	metricCache.UpdateMetric("default/name", metricRequest)
	client := NewAzureMetricClient("", metricCache, monitorClient)

	metricValue, err := client.GetAzureMetric("default", "name", nil)

	if err != nil {
		t.Errorf("error after processing got %v, want nil", err)
	}

	if metricValue.MetricName != metricRequest.MetricName {
		t.Errorf("error after processing got %v, want nil", err)
	}

	valueReturned := metricValue.Value.MilliValue()
	if valueReturned != int64(15000) {
		t.Errorf("MilliValue() got %v, want 15000", valueReturned)
	}
}

func newMetricRequest() azmetricrequest.AzureMetricRequest {
	return azmetricrequest.AzureMetricRequest{
		ResourceGroup:             "ResourceGroup",
		ResourceName:              "ResourceName",
		ResourceProviderNamespace: "ResourceProviderNamespace",
		ResourceType:              "ResourceType",
		SubscriptionID:            "SubscriptionID",
		MetricName:                "MetricName",
		Filter:                    "Filter",
		Aggregation:               "Aggregation",
	}
}

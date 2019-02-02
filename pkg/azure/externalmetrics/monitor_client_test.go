package externalmetrics

import (
	"context"
	"errors"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2018-03-01/insights"
)

func TestAzureMonitorIfEmptyRequestGetError(t *testing.T) {
	monitorClient := newFakeMonitorClient(insights.Response{}, nil)

	client := newMonitorClient("", monitorClient)

	request := AzureExternalMetricRequest{}
	_, err := client.GetAzureMetric(request)

	if err == nil {
		t.Errorf("no error after processing got: %v, want error", nil)
	}

	if !IsInvalidMetricRequestError(err) {
		t.Errorf("should be InvalidMetricRequest error got %v, want InvalidMetricRequestError", err)
	}
}

func TestAzureMonitorIfFailedResponseGetError(t *testing.T) {
	fakeError := errors.New("fake monitor failed")
	monitorClient := newFakeMonitorClient(insights.Response{}, fakeError)

	client := newMonitorClient("", monitorClient)

	request := newAzureMonitorMetricRequest()
	_, err := client.GetAzureMetric(request)

	if err == nil {
		t.Errorf("no error after processing got: %v, want error", nil)
	}

	if err.Error() != fakeError.Error() {
		t.Errorf("should be InvalidMetricRequest error got: %v, want: %v", err.Error(), fakeError.Error())
	}
}

func TestAzureMonitorIfValidRequestGetResult(t *testing.T) {
	response := makeAzureMonitorResponse(15)
	monitorClient := newFakeMonitorClient(response, nil)

	client := newMonitorClient("", monitorClient)

	request := newAzureMonitorMetricRequest()
	metricResponse, err := client.GetAzureMetric(request)

	if err != nil {
		t.Errorf("error after processing got: %v, want nil", err)
	}

	if metricResponse.Total != 15 {
		t.Errorf("metricResponse.Total = %v, want = %v", metricResponse.Total, 15)
	}
}

func makeAzureMonitorResponse(value float64) insights.Response {
	// create metric value
	mv := insights.MetricValue{
		Total: &value,
	}
	metricValues := []insights.MetricValue{}
	metricValues = append(metricValues, mv)

	// create timeseries
	te := insights.TimeSeriesElement{
		Data: &metricValues,
	}
	timeseries := []insights.TimeSeriesElement{}
	timeseries = append(timeseries, te)

	// create metric
	metric := insights.Metric{
		Timeseries: &timeseries,
	}
	metrics := []insights.Metric{}
	metrics = append(metrics, metric)

	// finish with response
	response := insights.Response{
		Value: &metrics,
	}
	return response
}

func newAzureMonitorMetricRequest() AzureExternalMetricRequest {
	return AzureExternalMetricRequest{
		ResourceGroup:             "ResourceGroup",
		ResourceName:              "ResourceName",
		ResourceProviderNamespace: "ResourceProviderNamespace",
		ResourceType:              "ResourceType",
		SubscriptionID:            "SubscriptionID",
		MetricName:                "MetricName",
		Filter:                    "Filter",
		Aggregation:               "Aggregation",
		Timespan:                  "PT10",
	}
}

func newFakeMonitorClient(result insights.Response, err error) insightsmonitorClient {
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

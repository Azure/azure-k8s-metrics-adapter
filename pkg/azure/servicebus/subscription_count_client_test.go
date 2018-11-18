package servicebus

import (
	"context"
	"errors"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
)

func TestIfEmptyRequestGetError(t *testing.T) {
	monitorClient := newFakeServicebusClient(servicebus.SBSubscription{}, nil)

	client := newClient("", monitorClient)

	request := AzureMetricRequest{}
	_, err := client.GetAzureMetric(request)

	if err == nil {
		t.Errorf("no error after processing got: %v, want error", nil)
	}

	if !IsInvalidMetricRequestError(err) {
		t.Errorf("should be InvalidMetricRequest error got %v, want InvalidMetricRequestError", err)
	}
}

func TestIfFailedResponseGetError(t *testing.T) {
	fakeError := errors.New("fake servicebus failed")
	serviceBusClient := newFakeServicebusClient(servicebus.SBSubscription{}, fakeError)

	client := newClient("", serviceBusClient)

	request := newMetricRequest()
	_, err := client.GetAzureMetric(request)

	if err == nil {
		t.Errorf("no error after processing got: %v, want error", nil)
	}

	if err.Error() != fakeError.Error() {
		t.Errorf("should be InvalidMetricRequest error got: %v, want: %v", err.Error(), fakeError.Error())
	}
}

func TestIfValidRequestGetResult(t *testing.T) {
	response := makeResponse(15)
	serviceBusClient := newFakeServicebusClient(response, nil)

	client := newClient("", serviceBusClient)

	request := newMetricRequest()
	metricResponse, err := client.GetAzureMetric(request)

	if err != nil {
		t.Errorf("error after processing got: %v, want nil", err)
	}

	if metricResponse.Total != 15 {
		t.Errorf("metricResponse.Total = %v, want = %v", metricResponse.Total, 15)
	}
}

func makeResponse(value int64) servicebus.SBSubscription {
	messageCountDetails := servicebus.MessageCountDetails{
		ActiveMessageCount: &value,
	}

	subscriptionProperties := servicebus.SBSubscriptionProperties{
		CountDetails: &messageCountDetails,
	}

	response := servicebus.SBSubscription{
		SBSubscriptionProperties: &subscriptionProperties,
	}

	return response
}

func newMetricRequest() AzureMetricRequest {
	return AzureMetricRequest{
		ResourceGroup:  "ResourceGroup",
		SubscriptionID: "SubscriptionID",
		MetricName:     "MetricName",
		Topic:          "Topic",
		Subscription:   "Subscription",
		Namespace:      "Namespace",
	}
}

func newFakeServicebusClient(result servicebus.SBSubscription, err error) fakeServicebusClient {
	return fakeServicebusClient{
		err:    err,
		result: result,
	}
}

type fakeServicebusClient struct {
	result servicebus.SBSubscription
	err    error
}

func (f fakeServicebusClient) Get(ctx context.Context, resourceGroupName string, namespaceName string, topicName string, subscriptionName string) (result servicebus.SBSubscription, err error) {
	result = f.result
	err = f.err
	return
}

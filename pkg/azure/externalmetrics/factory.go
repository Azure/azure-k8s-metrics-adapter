package externalmetrics

import "fmt"

type AzureClientFactory interface {
	GetAzureExternalMetricClient(clientType string) (AzureExternalMetricClient, error)
}

type AzureExternalMetricClientFactory struct {
	DefaultSubscriptionID string
}

func (f AzureExternalMetricClientFactory) GetAzureExternalMetricClient(clientType string) (client AzureExternalMetricClient, err error) {
	switch clientType {
	case Monitor:
		client = NewMonitorClient(f.DefaultSubscriptionID)
	case ServiceBusSubscription:
		client = NewServiceBusSubscriptionClient(f.DefaultSubscriptionID)
	default:
		err = fmt.Errorf("Unknown Azure external metric client type provided: %s", clientType)
	}

	return client, err
}

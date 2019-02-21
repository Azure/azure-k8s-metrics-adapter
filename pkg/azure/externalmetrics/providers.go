package externalmetrics

type AzureExternalMetricClientProvider interface {
	NewClient(defaultSubscriptionID string)
}

const (
	Monitor                string = "azuremonitor"
	ServiceBusSubscription string = "servicebussubscription"
)

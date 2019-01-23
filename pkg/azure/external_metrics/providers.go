package azureexternalmetrics

type AzureExternalMetricClientProvider interface {
	NewClient(defaultSubscriptionID string)
}

const (
	Monitor                string = "monitor"
	ServiceBusSubscription string = "servicebussubscription"
)

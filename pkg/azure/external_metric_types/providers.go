package externalmetrictypes

type AzureExternalMetricProvider string

const (
	Monitor                AzureExternalMetricProvider = "monitor"
	ServiceBusSubscription AzureExternalMetricProvider = "servicebussubscription"
)

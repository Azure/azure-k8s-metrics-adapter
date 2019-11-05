package externalmetrics

type AzureExternalMetricResponse struct {
	Value float64
}

type AzureExternalMetricClient interface {
	GetAzureMetric(azMetricRequest AzureExternalMetricRequest) (AzureExternalMetricResponse, error)
}

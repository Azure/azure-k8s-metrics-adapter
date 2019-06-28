package externalmetrics

type AzureExternalMetricResponse struct {
	Total float64
}

type AzureExternalMetricClient interface {
	GetAzureMetric(azMetricRequest AzureExternalMetricRequest) (AzureExternalMetricResponse, error)
}

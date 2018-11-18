package externalmetrictypes

type AzureExternalMetricResponse struct {
	Total int64
}

type AzureExternalMetricClient interface {
	GetAzureMetric(azMetricRequest AzureExternalMetricRequest) (AzureExternalMetricResponse, error)
}

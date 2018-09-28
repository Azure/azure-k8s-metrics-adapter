package monitor

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2018-03-01/insights"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/golang/glog"
)

// AzureMonitorClient provides an interface to make requests to Azure Monitor
type AzureMonitorClient interface {
	GetAzureMetric(azMetricRequest AzureMetricRequest) (AzureMetricResponse, error)
}

type AzureMetricResponse struct {
	Total float64
}

type insightsmonitorClient interface {
	List(ctx context.Context, resourceURI string, timespan string, interval *string, metricnames string, aggregation string, top *int32, orderby string, filter string, resultType insights.ResultType, metricnamespace string) (result insights.Response, err error)
}

type monitorClient struct {
	client                insightsmonitorClient
	DefaultSubscriptionID string
}

func NewClient(defaultsubscriptionID string) AzureMonitorClient {
	client := insights.NewMetricsClient(defaultsubscriptionID)
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		client.Authorizer = authorizer
	}

	return &monitorClient{
		client:                client,
		DefaultSubscriptionID: defaultsubscriptionID,
	}
}

func newClient(defaultsubscriptionID string, client insightsmonitorClient) monitorClient {
	return monitorClient{
		client:                client,
		DefaultSubscriptionID: defaultsubscriptionID,
	}
}

// GetAzureMetric calls Azure Monitor endpoint and returns a metric
func (c *monitorClient) GetAzureMetric(azMetricRequest AzureMetricRequest) (AzureMetricResponse, error) {
	err := azMetricRequest.Validate()
	if err != nil {
		return AzureMetricResponse{}, err
	}

	metricResourceURI := azMetricRequest.MetricResourceURI()
	glog.V(2).Infof("resource uri: %s", metricResourceURI)

	metricResult, err := c.client.List(context.Background(), metricResourceURI,
		azMetricRequest.Timespan, nil,
		azMetricRequest.MetricName, azMetricRequest.Aggregation, nil,
		"", azMetricRequest.Filter, "", "")
	if err != nil {
		return AzureMetricResponse{}, err
	}

	total := extractValue(metricResult)

	glog.V(2).Infof("found metric value: %f", total)

	// TODO set Value based on aggregations type
	return AzureMetricResponse{
		Total: total,
	}, nil
}

func extractValue(metricResult insights.Response) float64 {
	//TODO extract value based on aggregation type
	//TODO check for nils
	metricVals := *metricResult.Value
	Timeseries := *metricVals[0].Timeseries
	data := *Timeseries[0].Data
	total := *data[len(data)-1].Total

	return total
}

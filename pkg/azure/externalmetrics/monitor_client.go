package externalmetrics

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2018-03-01/insights"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/golang/glog"
)

type insightsmonitorClient interface {
	List(ctx context.Context, resourceURI string, timespan string, interval *string, metricnames string, aggregation string, top *int32, orderby string, filter string, resultType insights.ResultType, metricnamespace string) (result insights.Response, err error)
}

type monitorClient struct {
	client                insightsmonitorClient
	DefaultSubscriptionID string
}

func NewMonitorClient(defaultsubscriptionID string) AzureExternalMetricClient {
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

func newMonitorClient(defaultsubscriptionID string, client insightsmonitorClient) monitorClient {
	return monitorClient{
		client:                client,
		DefaultSubscriptionID: defaultsubscriptionID,
	}
}

// GetAzureMetric calls Azure Monitor endpoint and returns a metric
func (c *monitorClient) GetAzureMetric(azMetricRequest AzureExternalMetricRequest) (AzureExternalMetricResponse, error) {
	err := azMetricRequest.Validate()
	if err != nil {
		return AzureExternalMetricResponse{}, err
	}

	metricResourceURI := azMetricRequest.MetricResourceURI()
	glog.V(2).Infof("resource uri: %s", metricResourceURI)

	metricResult, err := c.client.List(context.Background(), metricResourceURI,
		azMetricRequest.Timespan, nil,
		azMetricRequest.MetricName, azMetricRequest.Aggregation, nil,
		"", azMetricRequest.Filter, "", "")
	if err != nil {
		return AzureExternalMetricResponse{}, err
	}

	total := extractValue(metricResult)

	glog.V(2).Infof("found metric value: %f", total)

	// TODO set Value based on aggregations type
	return AzureExternalMetricResponse{
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

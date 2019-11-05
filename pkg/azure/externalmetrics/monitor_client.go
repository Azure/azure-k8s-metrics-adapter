package externalmetrics

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2018-03-01/insights"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"k8s.io/klog"
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
	klog.V(2).Infof("resource uri: %s", metricResourceURI)

	metricResult, err := c.client.List(context.Background(), metricResourceURI,
		azMetricRequest.Timespan, nil,
		azMetricRequest.MetricName, azMetricRequest.Aggregation, nil,
		"", azMetricRequest.Filter, "", "")
	if err != nil {
		return AzureExternalMetricResponse{}, err
	}

	value, err := extractValue(azMetricRequest, metricResult)

	// TODO set Value based on aggregations type
	return AzureExternalMetricResponse{
		Value: value,
	}, err
}

func extractValue(azMetricRequest AzureExternalMetricRequest, metricResult insights.Response) (float64, error) {
	metricVals := *metricResult.Value
	Timeseries := *metricVals[0].Timeseries
	data := *Timeseries[0].Data

	var valuePtr *float64
	switch insights.AggregationType(azMetricRequest.Aggregation) {
	case insights.Average:
		if data[len(data)-1].Average != nil {
			valuePtr = data[len(data)-1].Average
		}
	case insights.Total:
		if data[len(data)-1].Total != nil {
			valuePtr = data[len(data)-1].Total
		}
	case insights.Maximum:
		if data[len(data)-1].Maximum != nil {
			valuePtr = data[len(data)-1].Maximum
		}
	case insights.Minimum:
		if data[len(data)-1].Minimum != nil {
			valuePtr = data[len(data)-1].Minimum
		}
	default:
		err := fmt.Errorf("Unsupported aggregation type %s specified in metric %s/%s", azMetricRequest.Aggregation, azMetricRequest.Namespace, azMetricRequest.MetricName)
		return 0, err
	}

	if valuePtr == nil {
		err := fmt.Errorf("Unable to get value for metric %s/%s with aggregation %s. No value returned by the Azure Monitor", azMetricRequest.Namespace, azMetricRequest.MetricName, azMetricRequest.Aggregation)
		return 0, err
	}

	klog.V(2).Infof("metric type: %s %f", azMetricRequest.Aggregation, *valuePtr)

	return *valuePtr, nil
}

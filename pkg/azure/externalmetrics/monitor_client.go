package externalmetrics

import (
	"context"
	"fmt"
	"strings"

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

	return AzureExternalMetricResponse{
		Value: value,
	}, err
}

func extractValue(azMetricRequest AzureExternalMetricRequest, metricResult insights.Response) (float64, error) {
	metricVals := *metricResult.Value

	if len(metricVals) == 0 {
		err := fmt.Errorf("Got an empty response for metric %s/%s and aggregate type %s", azMetricRequest.Namespace, azMetricRequest.MetricName, insights.AggregationType(strings.ToTitle(azMetricRequest.Aggregation)))
		return 0, err
	}

	timeseries := *metricVals[0].Timeseries
	if timeseries == nil {
		err := fmt.Errorf("Got metric result for %s/%s and aggregate type %s without timeseries", azMetricRequest.Namespace, azMetricRequest.MetricName, insights.AggregationType(strings.ToTitle(azMetricRequest.Aggregation)))
		return 0, err
	}

	data := *timeseries[0].Data
	if data == nil {
		err := fmt.Errorf("Got metric result for %s/%s and aggregate type %s without any metric values", azMetricRequest.Namespace, azMetricRequest.MetricName, insights.AggregationType(strings.ToTitle(azMetricRequest.Aggregation)))
		return 0, err
	}

	var valuePtr *float64
	if strings.EqualFold(string(insights.Average), azMetricRequest.Aggregation) && data[len(data)-1].Average != nil {
		valuePtr = data[len(data)-1].Average
	} else if strings.EqualFold(string(insights.Total), azMetricRequest.Aggregation) && data[len(data)-1].Total != nil {
		valuePtr = data[len(data)-1].Total
	} else if strings.EqualFold(string(insights.Maximum), azMetricRequest.Aggregation) && data[len(data)-1].Maximum != nil {
		valuePtr = data[len(data)-1].Maximum
	} else if strings.EqualFold(string(insights.Minimum), azMetricRequest.Aggregation) && data[len(data)-1].Minimum != nil {
		valuePtr = data[len(data)-1].Minimum
	} else if strings.EqualFold(string(insights.Count), azMetricRequest.Aggregation) && data[len(data)-1].Count != nil {
		fValue := float64(*data[len(data)-1].Count)
		valuePtr = &fValue
	} else {
		err := fmt.Errorf("Unsupported aggregation type %s specified in metric %s/%s", insights.AggregationType(strings.ToTitle(azMetricRequest.Aggregation)), azMetricRequest.Namespace, azMetricRequest.MetricName)
		return 0, err
	}

	if valuePtr == nil {
		err := fmt.Errorf("Unable to get value for metric %s/%s with aggregation %s. No value returned by the Azure Monitor", azMetricRequest.Namespace, azMetricRequest.MetricName, azMetricRequest.Aggregation)
		return 0, err
	}

	klog.V(2).Infof("metric type: %s %f", azMetricRequest.Aggregation, *valuePtr)

	return *valuePtr, nil
}

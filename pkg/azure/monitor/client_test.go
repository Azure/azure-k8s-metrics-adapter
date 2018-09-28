package monitor

// func TestIfNotValidCallGetError(t *testing.T) {
// 	monitorClient := newFakeMonitorClient(insights.Response{}, nil)
// 	metricCache := metriccache.NewMetricCache()

// 	client := NewAzureMetricClient("", metricCache, monitorClient)

// 	_, err := client.GetAzureMetric("ns", "name", nil)

// 	if err == nil {
// 		t.Errorf("error after processing got nil, want non nil")
// 	}
// }

// func TestIfInsufficientDataGetError(t *testing.T) {
// 	monitorClient := newFakeMonitorClient(insights.Response{}, nil)

// 	client := NewAzureMetricClient("", metricCache, monitorClient)

// 	// This doesn't have the all requeired selectors so should report that it is missing
// 	selector, _ := labels.Parse("resourceProviderNamespace=Microsoft.Servicebus")
// 	_, err := client.GetAzureMetric("ns", "name", selector)

// 	if !IsInvalidMetricRequestError(err) {
// 		t.Errorf("should be InvalidMetricRequest error got %v, want InvalidMetricRequestError", err)
// 	}
// }

// func TestIfPassedViaLabelSelectorsItReturns(t *testing.T) {
// 	response := makeResponse(15)

// 	monitorClient := newFakeMonitorClient(response, nil)

// 	client := NewAzureMetricClient("1234", metricCache, monitorClient)

// 	selector, _ := labels.Parse("resourceProviderNamespace=Microsoft.Servicebus,resourceType=namespaces,aggregation=Total,filter=EntityName_eq_externalq,resourceGroup=sb-external-example,resourceName=sb-external-ns,metricName=Messages")

// 	metricValue, err := client.GetAzureMetric("ns", "name", selector)

// 	if err != nil {
// 		t.Errorf("error after processing got %v, want nil", err)
// 	}

// 	if metricValue.MetricName != "Messages" {
// 		t.Errorf("error after processing got %v, want nil", err)
// 	}

// 	valueReturned := metricValue.Value.MilliValue()
// 	if valueReturned != int64(15000) {
// 		t.Errorf("MilliValue() got %v, want 15000", valueReturned)
// 	}
// }

// func TestIfCacheHasItReturn(t *testing.T) {
// 	monitorClient := newFakeMonitorClient(response, nil)

// 	metricRequest := newMetricRequest()
// 	metricCache.Update("default/name", metricRequest)
// 	client := NewAzureMetricClient("", metricCache, monitorClient)

// 	metricValue, err := client.GetAzureMetric("default", "name", nil)

// 	if err != nil {
// 		t.Errorf("error after processing got %v, want nil", err)
// 	}

// 	if metricValue.MetricName != metricRequest.MetricName {
// 		t.Errorf("error after processing got %v, want nil", err)
// 	}

// 	valueReturned := metricValue.Value.MilliValue()
// 	if valueReturned != int64(15000) {
// 		t.Errorf("MilliValue() got %v, want 15000", valueReturned)
// 	}
// }

// func makeResponse(value float64) insights.Response {
// 	// create metric value
// 	mv := insights.MetricValue{
// 		Total: &value,
// 	}
// 	metricValues := []insights.MetricValue{}
// 	metricValues = append(metricValues, mv)

// 	// create timeseries
// 	te := insights.TimeSeriesElement{
// 		Data: &metricValues,
// 	}
// 	timeseries := []insights.TimeSeriesElement{}
// 	timeseries = append(timeseries, te)

// 	// create metric
// 	metric := insights.Metric{
// 		Timeseries: &timeseries,
// 	}
// 	metrics := []insights.Metric{}
// 	metrics = append(metrics, metric)

// 	// finish with response
// 	response := insights.Response{
// 		Value: &metrics,
// 	}
// 	return response
// }

// func newMetricRequest() AzureMetricRequest {
// 	return AzureMetricRequest{
// 		ResourceGroup:             "ResourceGroup",
// 		ResourceName:              "ResourceName",
// 		ResourceProviderNamespace: "ResourceProviderNamespace",
// 		ResourceType:              "ResourceType",
// 		SubscriptionID:            "SubscriptionID",
// 		MetricName:                "MetricName",
// 		Filter:                    "Filter",
// 		Aggregation:               "Aggregation",
// 	}
// }

// func newFakeMonitorClient(result insights.Response, err error) MonitorClient {
// 	return fakeMonitorClient{
// 		err:    err,
// 		result: result,
// 	}
// }

// type fakeMonitorClient struct {
// 	result insights.Response
// 	err    error
// }

// func (f fakeMonitorClient) List(ctx context.Context, resourceURI string, timespan string, interval *string, metricnames string, aggregation string, top *int32, orderby string, filter string, resultType insights.ResultType, metricnamespace string) (result insights.Response, err error) {
// 	result = f.result
// 	err = f.err
// 	return
// }

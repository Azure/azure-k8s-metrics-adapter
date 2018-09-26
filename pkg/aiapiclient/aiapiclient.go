package aiapiclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Azure/azure-sdk-for-go/services/appinsights/v1/insights"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/golang/glog"
)

const (
	defaultAPIUrl   = "api.applicationinsights.io"
	apiVersion      = "v1"
	azureAdResource = "https://api.applicationinsights.io"
)

// AiAPIClient is used to call Application Insights Api
type AiAPIClient struct {
	appID           string
	appKey          string
	useADAuthorizer bool
}

// NewAiAPIClient creates a client for calling Application
// insights api
func NewAiAPIClient() AiAPIClient {
	defaultAppInsightsAppID := os.Getenv("APP_INSIGHTS_APP_ID")
	appInsightsKey := os.Getenv("APP_INSIGHTS_KEY")

	// if no application insights key has been specified, then we will use AD authentication
	return AiAPIClient{
		appID:           defaultAppInsightsAppID,
		appKey:          appInsightsKey,
		useADAuthorizer: appInsightsKey == "",
	}
}

// GetMetric calls to API to retrieve a specific metric
func (ai AiAPIClient) GetMetric(metricInfo MetricRequest) (*insights.ListMetricsResultsItem, error) {
	if ai.useADAuthorizer {
		return getMetricUsingADAuthorizer(ai, metricInfo)
	}

	return getMetricUsingAPIKey(ai, metricInfo)
}

func getMetricUsingADAuthorizer(ai AiAPIClient, metricInfo MetricRequest) (*insights.ListMetricsResultsItem, error) {

	authorizer, err := auth.NewAuthorizerFromEnvironmentWithResource(azureAdResource)
	if err != nil {
		glog.Errorf("unable to retrieve an authorizer from environment: %v", err)
		return nil, err
	}

	metricsClient := insights.NewMetricsClient()
	metricsClient.Authorizer = authorizer

	aggregation := []insights.MetricsAggregation{insights.MetricsAggregation(metricInfo.Aggregation)}
	metricsBodyParameter := insights.MetricsPostBodySchemaParameters{
		Interval:    &metricInfo.Interval,
		Timespan:    &metricInfo.Timespan,
		MetricID:    insights.MetricID(metricInfo.MetricName),
		Aggregation: &aggregation,
		Filter:      &metricInfo.Filter,
		Orderby:     &metricInfo.OrderBy,
	}
	metricsBody := []insights.MetricsPostBodySchema{
		insights.MetricsPostBodySchema{
			Parameters: &metricsBodyParameter,
		},
	}

	metricsResult, err := metricsClient.GetMultiple(context.Background(), ai.appID, metricsBody)
	if err != nil {
		glog.Errorf("unable to get retrive metric: %v", err)
		return nil, err
	}

	return &metricsResult, nil
}

func getMetricUsingAPIKey(ai AiAPIClient, metricInfo MetricRequest) (*insights.ListMetricsResultsItem, error) {
	client := &http.Client{}

	request := fmt.Sprintf("/%s/apps/%s/metrics/%s", apiVersion, ai.appID, metricInfo.MetricName)

	req, _ := http.NewRequest("GET", fmt.Sprintf("https://%s%s", defaultAPIUrl, request), nil)
	req.Header.Add("x-api-key", ai.appKey)

	q := req.URL.Query()
	q.Add("timespan", metricInfo.Timespan)
	q.Add("interval", metricInfo.Interval)
	req.URL.RawQuery = q.Encode()

	glog.V(2).Infoln("request to: ", req.URL)
	resp, err := client.Do(req)
	if err != nil {
		glog.Errorf("unable to get retrive metric: %v", err)
		return nil, err
	}

	response := insights.ListMetricsResultsItem{}

	return unmarshalResponse(resp.Body, &response)
}

func unmarshalResponse(body io.ReadCloser, response *insights.ListMetricsResultsItem) (*insights.ListMetricsResultsItem, error) {
	defer body.Close()
	respBody, err := ioutil.ReadAll(body)
	if err != nil {
		glog.Errorf("unable to get read metric response body: %v", err)
		return nil, err
	}

	err = json.Unmarshal(respBody, response)
	if err != nil {
		return nil, errors.New("unknown response format")
	}

	return response, nil
}

// // MetricsResponse is the response from the api that holds metric values and segments
// type MetricsResponse struct {
// 	StatusCode int
// 	Value      struct {
// 		Start        time.Time `json:"start"`
// 		End          time.Time `json:"end"`
// 		Interval     string    `json:"interval"`
// 		Segments     []Segment `json:"segments"`
// 		MetricValues Segment
// 	} `json:"value"`
// }

// // Segment holds the metric values for a given segment
// type Segment struct {
// 	Start        time.Time `json:"start"`
// 	End          time.Time `json:"end"`
// 	MetricValues map[string]map[string]interface{}
// }

// // UnmarshalJSON is a custom UnMarshaler that parses the Segment information
// func (s *Segment) UnmarshalJSON(b []byte) error {
// 	var segments map[string]interface{}
// 	if err := json.Unmarshal(b, &segments); err != nil {
// 		return err
// 	}

// 	for key, value := range segments {
// 		switch key {
// 		case "start":
// 			t, err := time.Parse(time.RFC3339, value.(string))
// 			if err != nil {
// 				return err
// 			}
// 			s.Start = t
// 		case "end":
// 			t, err := time.Parse(time.RFC3339, value.(string))
// 			if err != nil {
// 				return err
// 			}
// 			s.End = t
// 		default:
// 			if s.MetricValues == nil {
// 				s.MetricValues = make(map[string]map[string]interface{})
// 			}
// 			s.MetricValues[key] = value.(map[string]interface{})
// 		}
// 	}

// 	return nil
// }

// MetricRequest represents options for the AI endpoint
type MetricRequest struct {
	MetricName  string
	Aggregation string
	Timespan    string
	Interval    string
	Segment     string
	OrderBy     string
	Filter      string
}

// NewMetricRequest creates a new metric request with defaults for optional parameters
func NewMetricRequest(metricName string) MetricRequest {
	return MetricRequest{
		MetricName: metricName,
	}
}

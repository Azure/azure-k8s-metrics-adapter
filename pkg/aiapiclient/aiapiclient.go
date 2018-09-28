package aiapiclient

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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
func (ai AiAPIClient) GetMetric(metricInfo MetricRequest) (*insights.MetricsResult, error) {
	if ai.useADAuthorizer {
		glog.V(2).Infoln("No application insights key provided - using Azure GO SDK auth.")
		return getMetricUsingADAuthorizer(ai, metricInfo)
	}

	glog.V(2).Infoln("Application insights key has been provided - using Application Insights REST API.")
	return getMetricUsingAPIKey(ai, metricInfo)
}

func getMetricUsingADAuthorizer(ai AiAPIClient, metricInfo MetricRequest) (*insights.MetricsResult, error) {

	authorizer, err := auth.NewAuthorizerFromEnvironmentWithResource(azureAdResource)
	if err != nil {
		glog.Errorf("unable to retrieve an authorizer from environment: %v", err)
		return nil, err
	}

	metricsClient := insights.NewMetricsClient()
	metricsClient.Authorizer = authorizer

	metricsBodyParameter := insights.MetricsPostBodySchemaParameters{
		Interval: &metricInfo.Interval,
		Timespan: &metricInfo.Timespan,
		MetricID: insights.MetricID(metricInfo.MetricName),
	}

	requestSchemaIdentifier := generateRequestSchemaUniqueIdentifier()
	metricsBody := []insights.MetricsPostBodySchema{
		insights.MetricsPostBodySchema{
			ID:         &requestSchemaIdentifier,
			Parameters: &metricsBodyParameter,
		},
	}

	metricsResultsItem, err := metricsClient.GetMultiple(context.Background(), ai.appID, metricsBody)
	if err != nil {
		glog.Errorf("unable to get retrive metric: %v", err)
		return nil, err
	}

	// check there is a metric result
	if len(*metricsResultsItem.Value) == 0 {
		return nil, errors.New("response from metrics request is empty")
	}

	// take only the first result (as we ask for a specific metric, there is only one result)
	metricsResult := (*metricsResultsItem.Value)[0]

	// check the body is not nil
	if metricsResult.Body == nil {
		return nil, errors.New("response from metrics request is empty")
	}

	return metricsResult.Body, nil
}

func generateRequestSchemaUniqueIdentifier() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return uuid
}

func getMetricUsingAPIKey(ai AiAPIClient, metricInfo MetricRequest) (*insights.MetricsResult, error) {
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
		glog.Errorf("unable to retrive metric: %v", err)
		return nil, err
	}

	// check the response status is OK. If not, return the error
	if resp.StatusCode != 200 {
		defer resp.Body.Close()
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			glog.Errorf("unable to retrieve metric: %s", err)
			return nil, err
		}

		respMessage := string(respBody)
		err = fmt.Errorf(respMessage)
		return nil, err
	}
	// return the response unmarshaled
	metricsResult := insights.MetricsResult{}
	return unmarshalResponse(resp.Body, &metricsResult)
}

func unmarshalResponse(body io.ReadCloser, metricsResult *insights.MetricsResult) (*insights.MetricsResult, error) {
	defer body.Close()
	respBody, err := ioutil.ReadAll(body)

	if err != nil {
		glog.Errorf("unable to get read metric response body: %v", err)
		return nil, err
	}

	err = json.Unmarshal(respBody, metricsResult)
	if err != nil {
		return nil, errors.New("unknown response format")
	}

	return metricsResult, nil
}

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

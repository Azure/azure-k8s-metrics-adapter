package custommetrics

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

// AzureAppInsightsClient provides methods for accessing App Insights via AD auth or App API Key
type AzureAppInsightsClient interface {
	GetCustomMetric(request MetricRequest) (float64, error)
}

// appinsightsClient is used to call Application Insights Api
type appinsightsClient struct {
	appID           string
	appKey          string
	useADAuthorizer bool
}

// NewClient creates a client for calling Application
// insights api
func NewClient() AzureAppInsightsClient {
	defaultAppInsightsAppID := os.Getenv("APP_INSIGHTS_APP_ID")
	appInsightsKey := os.Getenv("APP_INSIGHTS_KEY")

	// if no application insights key has been specified, then we will use AD authentication
	return appinsightsClient{
		appID:           defaultAppInsightsAppID,
		appKey:          appInsightsKey,
		useADAuthorizer: appInsightsKey == "",
	}
}

// GetCustomMetric calls to Application Insights to retrieve the value of the metric requested
func (c appinsightsClient) GetCustomMetric(request MetricRequest) (float64, error) {

	// get the last 5 mins and chunking into 30 seconds
	// this seems to be the best way to get the closest average rate at time of request
	// any smaller time intervals and the values come back null
	// TODO make this configurable?
	request.Timespan = "PT5M"
	request.Interval = "PT30S"

	metricsResult, err := c.getMetric(request)
	if err != nil {
		return 0, err
	}

	if metricsResult.Value == nil || metricsResult.Value.Segments == nil {
		return 0, errors.New("metrics result is nil")
	}

	segments := *metricsResult.Value.Segments
	if len(segments) <= 0 {
		glog.V(2).Info("segments length = 0")
		return 0, nil
	}

	// grab just the last value which will be the latest value of the metric
	metric := segments[len(segments)-1].AdditionalProperties[request.MetricName]
	metricMap := metric.(map[string]interface{})
	value := metricMap["avg"]
	normalizedValue := normalizeValue(value)

	glog.V(2).Infof("found metric value: %f", normalizedValue)
	return normalizedValue, nil
}

func normalizeValue(value interface{}) float64 {
	switch t := value.(type) {
	case int32:
		return float64(value.(int32))
	case float32:
		return float64(value.(float32))
	case float64:
		return value.(float64)
	case int64:
		return float64(value.(int64))
	default:
		glog.V(0).Infof("unexpected type: %T", t)
		return 0
	}
}

// GetMetric calls to API to retrieve a specific metric
func (ai appinsightsClient) getMetric(metricInfo MetricRequest) (*insights.MetricsResult, error) {
	if ai.useADAuthorizer {
		glog.V(2).Infoln("No application insights key provided - using Azure GO SDK auth.")
		return getMetricUsingADAuthorizer(ai, metricInfo)
	}

	glog.V(2).Infoln("Application insights key has been provided - using Application Insights REST API.")
	return getMetricUsingAPIKey(ai, metricInfo)
}

func getMetricUsingADAuthorizer(ai appinsightsClient, metricInfo MetricRequest) (*insights.MetricsResult, error) {

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

func getMetricUsingAPIKey(ai appinsightsClient, metricInfo MetricRequest) (*insights.MetricsResult, error) {
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

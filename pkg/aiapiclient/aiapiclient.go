package aiapiclient

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/golang/glog"
)

// AiAPIClient is used to call Application Insights Api
type AiAPIClient struct {
	appID  string
	appKey string
}

const (
	defaultAPIUrl = "api.applicationinsights.io"
	apiVersion    = "v1"
)

// NewAiAPIClient creates a client for calling Application
// insights api
func NewAiAPIClient() AiAPIClient {
	defaultAppInsightsAppID := os.Getenv("APP_INSIGHTS_APP_ID")
	appInsightsKey := os.Getenv("APP_INSIGHTS_KEY")

	return AiAPIClient{
		appID:  defaultAppInsightsAppID,
		appKey: appInsightsKey,
	}
}

// GetMetric calls to API to retrieve a specific metric
func (ai AiAPIClient) GetMetric(metric, aggregation string) (*int64, error) {

	client := &http.Client{}

	request := fmt.Sprintf("/%s/apps/%s/metrics/%s", apiVersion, ai.appID, metric)

	req, _ := http.NewRequest("GET", fmt.Sprintf("https://%s%s", defaultAPIUrl, request), nil)
	req.Header.Add("x-api-key", ai.appKey)

	timespan := "PT5M"
	interval := "PT30S"

	q := req.URL.Query()
	q.Add("timespan", timespan)
	q.Add("interval", interval)
	req.URL.RawQuery = q.Encode()

	glog.V(2).Infoln("request to: ", req.URL)
	resp, err := client.Do(req)
	if err != nil {
		glog.Errorf("unable to get retrive metric: %v", err)
		return nil, err
	}

	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	response := string(respBody[:])

	glog.V(2).Infoln("response", response)

	//segments := *(result.Value.Segments)
	//value := segments[len(segments)-1].AdditionalProperties["performanceCounters/requestsPerSecond"]

	//glog.V(2).Infof("perf/rps", value)
	//v := value.(map[string]int64)

	v := int64(10)

	return &v, nil
}

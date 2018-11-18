package servicebus

import (
	"fmt"
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/labels"
)

const subscription = "test-sub"
const topic = "test-topic"
const namespace = "test-namespace"
const resourceGroup = "test-resource-group"
const metricName = "metric-name"
const subscriptionID = "1234-5678"

var validLabelSelector = fmt.Sprintf("subscription=%s,topic=%s,namespace=%s,resourceGroup=%s,metricName=%s", subscription, topic, namespace, resourceGroup, metricName)

type testArguments struct {
	metricSelector        string
	defaultSubscriptionID string
}

// List of test cases that will be run for validating the parsing of metric configuration
var testCases = []struct {
	name     string
	args     testArguments
	want     AzureMetricRequest
	wantErr  bool
	validate bool
}{
	// Begin test cases
	{
		name: "Test if metricSelector is nil do not fail",
		args: testArguments{
			defaultSubscriptionID: "",
			metricSelector:        "",
		},
		want:     AzureMetricRequest{},
		wantErr:  true,
		validate: false,
	},
	{
		name: "Test insufficient data expect an error",
		args: testArguments{
			defaultSubscriptionID: "",
			metricSelector:        "namespace=testing",
		},
		want:     AzureMetricRequest{},
		wantErr:  true,
		validate: true,
	},
	{
		name: "Test valid case with overriding subscription ID passed in",
		args: testArguments{
			defaultSubscriptionID: subscriptionID,
			metricSelector:        validLabelSelector,
		},
		want: AzureMetricRequest{
			Namespace:      namespace,
			Subscription:   subscription,
			Topic:          topic,
			MetricName:     metricName,
			ResourceGroup:  resourceGroup,
			SubscriptionID: subscriptionID,
		},
		wantErr:  false,
		validate: true,
	},
	{
		name: "Test valid case with overriding subscription ID in selector",
		args: testArguments{
			defaultSubscriptionID: subscriptionID,
			metricSelector:        fmt.Sprintf("%s,subscriptionID=%s", validLabelSelector, subscriptionID),
		},
		want: AzureMetricRequest{
			Namespace:      namespace,
			Subscription:   subscription,
			Topic:          topic,
			MetricName:     metricName,
			ResourceGroup:  resourceGroup,
			SubscriptionID: subscriptionID,
		},
		wantErr:  false,
		validate: true,
	},
	{
		name: "Test valid case with overriding subscription ID in selector",
		args: testArguments{
			defaultSubscriptionID: subscriptionID,
			metricSelector:        fmt.Sprintf("%s,subscriptionID=%s", validLabelSelector, subscriptionID),
		},
		want: AzureMetricRequest{
			Namespace:      namespace,
			Subscription:   subscription,
			Topic:          topic,
			MetricName:     metricName,
			ResourceGroup:  resourceGroup,
			SubscriptionID: subscriptionID,
		},
		wantErr:  false,
		validate: true,
	},
}

// Test the parsing of the External Metric Configuration that is expected from
// the Custom Resource Definition for the External Metric Provider
func TestParsingAzureExternalMetricConfiguration(t *testing.T) {
	// Run through the test cases and valid expected outcomes
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			metricSelector, err := labels.Parse(tt.args.metricSelector)
			if err != nil {
				t.Errorf("ParseAzureMetric() error parsing metricSelector %s", metricSelector)
			}

			if len(tt.args.metricSelector) == 0 {
				metricSelector = nil
			}

			got, parseErr := ParseAzureMetric(metricSelector, tt.args.defaultSubscriptionID)

			if tt.validate {
				err = got.Validate()
				if (err != nil) != tt.wantErr {
					t.Errorf("ParseAzureMetric() validation error = %v", err)
				}
			}

			if err != nil {
				return
			}

			if (parseErr != nil) != tt.wantErr {
				t.Errorf("ParseAzureMetric() error = %v, wantErr %v", parseErr, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseAzureMetric() = %v, want %v", got, tt.want)
				return
			}
		})
	}
}

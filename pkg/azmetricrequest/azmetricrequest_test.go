package azmetricrequest

import (
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/labels"
)

func Test_azureMetricRequest_metricResourceURI(t *testing.T) {
	tests := []struct {
		name string
		amr  AzureMetricRequest
		want string
	}{
		{
			name: "valid metric",
			amr: AzureMetricRequest{
				SubscriptionID:            "1234-1234-234-12414",
				ResourceGroup:             "test-rg",
				ResourceProviderNamespace: "Microsoft.Servicebus",
				ResourceType:              "namespaces",
				ResourceName:              "sb-external-ns",
			},
			want: "/subscriptions/1234-1234-234-12414/resourceGroups/test-rg/providers/Microsoft.Servicebus/namespaces/sb-external-ns",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.amr.MetricResourceURI(); got != tt.want {
				t.Errorf("azureMetricRequest.metricResourceURI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseAzureMetric(t *testing.T) {
	type args struct {
		metricSelector        labels.Selector
		defaultSubscriptionID string
	}
	tests := []struct {
		name    string
		args    args
		want    AzureMetricRequest
		wantErr bool
	}{
		{
			name: "if metricSelector is nil do not fail",
			args: args{
				defaultSubscriptionID: "",
				metricSelector:        nil,
			},
			want:    AzureMetricRequest{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAzureMetric(tt.args.metricSelector, tt.args.defaultSubscriptionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseAzureMetric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseAzureMetric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ParseInsufficientDataGetError(t *testing.T) {
	// This doesn't have the all requeired selectors so should report that it is missing
	selector, _ := labels.Parse("resourceProviderNamespace=Microsoft.Servicebus")

	metric, _ := ParseAzureMetric(selector, "")

	err := metric.Validate()

	if !IsInvalidMetricRequestError(err) {
		t.Errorf("should be InvalidMetricRequest error got %v, want InvalidMetricRequestError", err)
	}
}

func Test_ParseWithSubIdPassedIsValid(t *testing.T) {
	// This doesn't have the all requeired selectors so should report that it is missing
	selector, _ := labels.Parse("resourceProviderNamespace=Microsoft.Servicebus,resourceType=namespaces,aggregation=Total,filter=EntityName_eq_externalq,resourceGroup=sb-external-example,resourceName=sb-external-ns,metricName=Messages")

	metric, _ := ParseAzureMetric(selector, "1234")

	err := metric.Validate()

	if err != nil {
		t.Errorf("validate got error %v, want nil", err)
	}
}

func Test_ParseWithSubIdOnSelectorPassedIsValid(t *testing.T) {
	// This doesn't have the all requeired selectors so should report that it is missing
	selector, _ := labels.Parse("subscriptionID=1234,resourceProviderNamespace=Microsoft.Servicebus,resourceType=namespaces,aggregation=Total,filter=EntityName_eq_externalq,resourceGroup=sb-external-example,resourceName=sb-external-ns,metricName=Messages")

	metric, _ := ParseAzureMetric(selector, "1234")

	err := metric.Validate()

	if err != nil {
		t.Errorf("validate got error %v, want nil", err)
	}
}

func Test_ParseWithNoSubIdPassedIsFails(t *testing.T) {
	// This doesn't have the all requeired selectors so should report that it is missing
	selector, _ := labels.Parse("resourceProviderNamespace=Microsoft.Servicebus,resourceType=namespaces,aggregation=Total,filter=EntityName_eq_externalq,resourceGroup=sb-external-example,resourceName=sb-external-ns,metricName=Messages")

	metric, _ := ParseAzureMetric(selector, "")

	err := metric.Validate()

	if !IsInvalidMetricRequestError(err) {
		t.Errorf("should be InvalidMetricRequest error got %v, want InvalidMetricRequestError", err)
	}
}

package azmetricrequest

import "testing"

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

package azureMetricClient

import (
	"testing"
)

func Test_normalizeValue(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "int64 to float64",
			args: args{
				value: int64(42),
			},
			want: float64(42),
		},
		{
			name: "float64 to float64",
			args: args{
				value: float64(42.0),
			},
			want: float64(42),
		},
		{
			name: "int32 to float64",
			args: args{
				value: int32(42),
			},
			want: float64(42),
		},
		{
			name: "float32 to float64",
			args: args{
				value: float32(42.0),
			},
			want: float64(42),
		},
		{
			name: "if something random like a string, return 0",
			args: args{
				value: "this is not the answer",
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeValue(tt.args.value); got != tt.want {
				t.Errorf("normalizeValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_azureMetricRequest_metricResourceURI(t *testing.T) {
	tests := []struct {
		name string
		amr  azureMetricRequest
		want string
	}{
		{
			name: "valid metric",
			amr: azureMetricRequest{
				subscriptionID:            "1234-1234-234-12414",
				resourceGroup:             "test-rg",
				resourceProviderNamespace: "Microsoft.Servicebus",
				resourceType:              "namespaces",
				resourceName:              "sb-external-ns",
			},
			want: "/subscriptions/1234-1234-234-12414/resourceGroups/test-rg/providers/Microsoft.Servicebus/namespaces/sb-external-ns",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.amr.metricResourceURI(); got != tt.want {
				t.Errorf("azureMetricRequest.metricResourceURI() = %v, want %v", got, tt.want)
			}
		})
	}
}

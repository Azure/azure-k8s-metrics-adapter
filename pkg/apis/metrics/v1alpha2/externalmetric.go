package v1alpha2

import (
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +genclient:skipVerbs=patch
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ExternalMetric describes a ExternalMetric resource
type ExternalMetric struct {
	// TypeMeta is the metadata for the resource, like kind and apiversion
	meta_v1.TypeMeta `json:",inline"`

	// ObjectMeta contains the metadata for the particular object (name, namespace, self link, labels, etc)
	meta_v1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the custom resource spec
	Spec ExternalMetricSpec `json:"spec"`
}

// ExternalMetricSpec is the spec for a ExternalMetric resource
type ExternalMetricSpec struct {
	MetricConfig ExternalMetricConfig `json:"metric"`
	AzureConfig  AzureConfig          `json:"azure"`
	Type         string               `json:"type,omitempty"`
}

// ExternalMetricConfig holds azure monitor metric configuration
type ExternalMetricConfig struct {
	// Shared
	MetricName string `json:"metricName,omitempty"`
	// Azure Monitor
	Aggregation string `json:"aggregation,omitempty"`
	Filter      string `json:"filter,omitempty"`
}

// AzureConfig holds Azure configuration for an External Metric
type AzureConfig struct {
	// Shared
	ResourceGroup  string `json:"resourceGroup"`
	SubscriptionID string `json:"subscriptionID"`
	// Azure Monitor
	ResourceName              string `json:"resourceName,omitempty"`
	ResourceProviderNamespace string `json:"resourceProviderNamespace,omitempty"`
	ResourceType              string `json:"resourceType,omitempty"`
	// Azure Service Bus Topic Subscription
	ServiceBusNamespace    string `json:"serviceBusNamespace,omitempty"`
	ServiceBusTopic        string `json:"serviceBusTopic,omitempty"`
	ServiceBusSubscription string `json:"serviceBusSubscription,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ExternalMetricList is a list of ExternalMetric resources
type ExternalMetricList struct {
	meta_v1.TypeMeta `json:",inline"`
	meta_v1.ListMeta `json:"metadata"`

	Items []ExternalMetric `json:"items"`
}

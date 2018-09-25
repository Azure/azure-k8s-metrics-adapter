package v1alpha1

import (
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
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
	MetricConfig MetricConfig `json:"metric"`
	AzureConfig  AzureConfig  `json:"azure"`
}

type MetricConfig struct {
	MetricName  string `json:"metricName"`
	Aggregation string `json:"aggregation"`
	Filter      string `json:"filter"`
}

type AzureConfig struct {
	ResourceGroup             string `json:"resourceGroup"`
	ResourceName              string `json:"resourceName"`
	ResourceProviderNamespace string `json:"resourceProviderNamespace"`
	ResourceType              string `json:"resourceType"`
	SubscriptionID            string `json:"subscriptionID"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ExternalMetricsList is a list of ExternalMetric resources
type ExternalMetricList struct {
	meta_v1.TypeMeta `json:",inline"`
	meta_v1.ListMeta `json:"metadata"`

	Items []ExternalMetric `json:"items"`
}

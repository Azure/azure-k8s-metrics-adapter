package v1alpha1

import (
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CustomMetric describes a configuration for Application insights
type CustomMetric struct {
	// TypeMeta is the metadata for the resource, like kind and apiversion
	meta_v1.TypeMeta `json:",inline"`

	// ObjectMeta contains the metadata for the particular object (name, namespace, self link, labels, etc)
	meta_v1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the custom resource spec
	Spec CustomMetricSpec `json:"spec"`
}

// CustomMetricSpec is the spec for a CustomMetric resource
type CustomMetricSpec struct {
	MetricConfig CustomMetricConfig `json:"metric"`
}

// CustomMetricConfig holds app insights configuration
type CustomMetricConfig struct {
	metricName    string `json:"metricName"`
	applicationID string `json:"applicationID"`
	query         string `json:"query"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CustomMetricList is a list of CustomMetric resources
type CustomMetricList struct {
	meta_v1.TypeMeta `json:",inline"`
	meta_v1.ListMeta `json:"metadata"`

	Items []CustomMetric `json:"items"`
}

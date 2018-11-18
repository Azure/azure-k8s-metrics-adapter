package servicebus

import (
	"errors"
	"fmt"

	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

type AzureMetricRequest struct {
	MetricName     string
	ResourceGroup  string
	Namespace      string
	Topic          string
	Subscription   string
	SubscriptionID string
}

func ParseAzureMetric(metricSelector labels.Selector, defaultSubscriptionID string) (AzureMetricRequest, error) {
	glog.V(4).Infof("Parsing a received AzureMetric")
	glog.V(6).Infof("%v", metricSelector)

	if metricSelector == nil {
		return AzureMetricRequest{}, fmt.Errorf("metricSelector cannot be nil")
	}

	// Using selectors to pass required values thorugh
	// to retain camel case as azure provider is case sensitive.
	//
	// There is are restrictions so using some conversion
	// restrictions here
	// note: requirement values are already validated by apiserver
	merticReq := AzureMetricRequest{
		SubscriptionID: defaultSubscriptionID,
	}
	requirements, _ := metricSelector.Requirements()
	for _, request := range requirements {
		if request.Operator() != selection.Equals {
			return AzureMetricRequest{}, errors.New("selector type not supported. only equals is supported at this time")
		}

		value := request.Values().List()[0]

		switch request.Key() {
		case "metricName":
			glog.V(4).Infof("AzureMetric metricName: %s", value)
			merticReq.MetricName = value
		case "resourceGroup":
			glog.V(4).Infof("AzureMetric resourceGroup: %s", value)
			merticReq.ResourceGroup = value
		case "namespace":
			glog.V(4).Infof("AzureMetric namespace: %s", value)
			merticReq.Namespace = value
		case "topic":
			glog.V(4).Infof("AzureMetric topic: %s", value)
			merticReq.Topic = value
		case "subscription":
			glog.V(4).Infof("AzureMetric subscription: %s", value)
			merticReq.Subscription = value
		case "subscriptionID":
			// if sub id is passed via label selectors then it takes precedence
			glog.V(4).Infof("AzureMetric override azure subscription id with : %s", value)
			merticReq.SubscriptionID = value
		default:
			return AzureMetricRequest{}, fmt.Errorf("selector label '%s' not supported", request.Key())
		}
	}

	glog.V(2).Infof("Successfully parsed AzureMetric %s", merticReq.MetricName)

	return merticReq, nil
}

type InvalidMetricRequestError struct {
	err string
}

func (i InvalidMetricRequestError) Error() string {
	return fmt.Sprintf(i.err)
}

func IsInvalidMetricRequestError(err error) bool {
	if _, ok := err.(InvalidMetricRequestError); ok {
		return true
	}
	return false
}

func (amr AzureMetricRequest) Validate() error {
	if amr.MetricName == "" {
		return InvalidMetricRequestError{err: "metricName is required"}
	}
	if amr.ResourceGroup == "" {
		return InvalidMetricRequestError{err: "resourceGroup is required"}
	}
	if amr.Namespace == "" {
		return InvalidMetricRequestError{err: "namespace is required"}
	}
	if amr.Topic == "" {
		return InvalidMetricRequestError{err: "topic is required"}
	}
	if amr.Subscription == "" {
		return InvalidMetricRequestError{err: "subscription is required"}
	}
	if amr.SubscriptionID == "" {
		return InvalidMetricRequestError{err: "subscriptionID is required. set a default or pass via label selectors"}
	}

	// if here then valid!
	return nil
}

package azmetricrequest

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

type AzureMetricRequest struct {
	MetricName                string
	ResourceGroup             string
	ResourceName              string
	ResourceProviderNamespace string
	ResourceType              string
	Aggregation               string
	Timespan                  string
	Filter                    string
	SubscriptionID            string
}

func ParseAzureMetric(metricSelector labels.Selector, defaultSubscriptionID string) (AzureMetricRequest, error) {
	glog.V(2).Infof("begin parsing metric")

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
		Timespan:       TimeSpan(),
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
			glog.V(2).Infof("metricName: %s", value)
			merticReq.MetricName = value
		case "resourceGroup":
			glog.V(2).Infof("resourceGroup: %s", value)
			merticReq.ResourceGroup = value
		case "resourceName":
			glog.V(2).Infof("resourceName: %s", value)
			merticReq.ResourceName = value
		case "resourceProviderNamespace":
			glog.V(2).Infof("resourceProviderNamespace: %s", value)
			merticReq.ResourceProviderNamespace = value
		case "resourceType":
			glog.V(2).Infof("resourceType: %s", value)
			merticReq.ResourceType = value
		case "aggregation":
			glog.V(2).Infof("aggregation: %s", value)
			merticReq.Aggregation = value
		case "filter":
			// TODO: Should handle filters by converting equality and setbased label selectors
			// to  oData syntax: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#label-selectors
			glog.V(2).Infof("filter: %s", value)
			filterStrings := strings.Split(value, "_")
			merticReq.Filter = fmt.Sprintf("%s %s '%s'", filterStrings[0], filterStrings[1], filterStrings[2])
			glog.V(2).Infof("filter formatted: %s", merticReq.Filter)
		case "subscriptionID":
			// if sub id is passed via label selectors then it takes precedence
			glog.V(2).Infof("override azure subscription id with : %s", value)
			merticReq.SubscriptionID = value
		default:
			return AzureMetricRequest{}, fmt.Errorf("selector label '%s' not supported", request.Key())
		}
	}

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
	if amr.ResourceName == "" {
		return InvalidMetricRequestError{err: "resourceName is required"}
	}
	if amr.ResourceProviderNamespace == "" {
		return InvalidMetricRequestError{err: "resourceProviderNamespace is required"}
	}
	if amr.ResourceType == "" {
		return InvalidMetricRequestError{err: "resourceType is required"}
	}
	if amr.Aggregation == "" {
		return InvalidMetricRequestError{err: "aggregation is required"}
	}
	if amr.Timespan == "" {
		return InvalidMetricRequestError{err: "timespan is required"}
	}
	if amr.Filter == "" {
		return InvalidMetricRequestError{err: "filter is required"}
	}

	if amr.SubscriptionID == "" {
		return InvalidMetricRequestError{err: "subscriptionID is required. set a default or pass via label selectors"}
	}

	// if here then valid!
	return nil
}

func (amr AzureMetricRequest) MetricResourceURI() string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/%s/%s/%s",
		amr.SubscriptionID,
		amr.ResourceGroup,
		amr.ResourceProviderNamespace,
		amr.ResourceType,
		amr.ResourceName)
}

// TimeSpan sets the default time to aggregate a metric
func TimeSpan() string {
	// defaults to last five minutes.
	// TODO support configuration via config
	endtime := time.Now().UTC().Format(time.RFC3339)
	starttime := time.Now().Add(-(5 * time.Minute)).UTC().Format(time.RFC3339)
	return fmt.Sprintf("%s/%s", starttime, endtime)
}

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

	// Using selectors to pass required values thorugh
	// to retain camel case as azure provider is case sensitive.
	//
	// There is are restrictions so using some conversion
	// restrictions here
	// note: requirement values are already validated by apiserver
	merticReq := AzureMetricRequest{
		Timespan:       timeSpan(),
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

	err := merticReq.Validate()
	if err != nil {
		return AzureMetricRequest{}, err
	}
	return merticReq, nil
}

func (amr AzureMetricRequest) Validate() error {
	if amr.MetricName == "" {
		return fmt.Errorf("metricName is required")
	}
	if amr.ResourceGroup == "" {
		return fmt.Errorf("resourceGroup is required")
	}
	if amr.ResourceName == "" {
		return fmt.Errorf("resourceName is required")
	}
	if amr.ResourceProviderNamespace == "" {
		return fmt.Errorf("resourceProviderNamespace is required")
	}
	if amr.ResourceType == "" {
		return fmt.Errorf("resourceType is required")
	}
	if amr.Aggregation == "" {
		return fmt.Errorf("aggregation is required")
	}
	if amr.Timespan == "" {
		return fmt.Errorf("timespan is required")
	}
	if amr.Filter == "" {
		return fmt.Errorf("filter is required")
	}

	if amr.SubscriptionID == "" {
		return fmt.Errorf("subscriptionID is required. set a default or pass via label selectors")
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

func timeSpan() string {
	// defaults to last five minutes.
	// TODO support configuration via config
	endtime := time.Now().UTC().Format(time.RFC3339)
	starttime := time.Now().Add(-(5 * time.Minute)).UTC().Format(time.RFC3339)
	return fmt.Sprintf("%s/%s", starttime, endtime)
}

package externalmetrics

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/golang/glog"
)

type servicebusSubscriptionsClient interface {
	Get(ctx context.Context, resourceGroupName string, namespaceName string, topicName string, subscriptionName string) (result servicebus.SBSubscription, err error)
}

type servicebusClient struct {
	client                servicebusSubscriptionsClient
	DefaultSubscriptionID string
}

func NewServiceBusSubscriptionClient(defaultSubscriptionID string) AzureExternalMetricClient {
	glog.V(2).Info("Creating a new Azure Service Bus Subscriptions client")
	client := servicebus.NewSubscriptionsClient(defaultSubscriptionID)
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		client.Authorizer = authorizer
	}

	return &servicebusClient{
		client:                client,
		DefaultSubscriptionID: defaultSubscriptionID,
	}
}

func newServiceBusSubscriptionClient(defaultsubscriptionID string, client servicebusSubscriptionsClient) servicebusClient {
	return servicebusClient{
		client:                client,
		DefaultSubscriptionID: defaultsubscriptionID,
	}
}

func (c *servicebusClient) GetAzureMetric(azMetricRequest AzureExternalMetricRequest) (AzureExternalMetricResponse, error) {
	glog.V(6).Infof("Received metric request:\n%v", azMetricRequest)
	err := azMetricRequest.Validate()
	if err != nil {
		return AzureExternalMetricResponse{}, err
	}

	glog.V(2).Infof("Requesting Service Bus Subscription %s to topic %s in namespace %s from resource group %s", azMetricRequest.Subscription, azMetricRequest.Topic, azMetricRequest.Namespace, azMetricRequest.ResourceGroup)
	subscriptionResult, err := c.client.Get(
		context.Background(),
		azMetricRequest.ResourceGroup,
		azMetricRequest.Namespace,
		azMetricRequest.Topic,
		azMetricRequest.Subscription,
	)
	if err != nil {
		return AzureExternalMetricResponse{}, err
	}

	glog.V(2).Infof("Successfully retrieved Service Bus Subscription %s to topic %s in namespace %s from resource group %s", azMetricRequest.Subscription, azMetricRequest.Topic, azMetricRequest.Namespace, azMetricRequest.ResourceGroup)
	glog.V(6).Infof("%v", subscriptionResult.Response)

	activeMessageCount := float64(*subscriptionResult.SBSubscriptionProperties.CountDetails.ActiveMessageCount)

	glog.V(4).Infof("Service Bus Subscription active message count: %f", activeMessageCount)

	// TODO set Value based on aggregations type
	return AzureExternalMetricResponse{
		Total: activeMessageCount,
	}, nil
}

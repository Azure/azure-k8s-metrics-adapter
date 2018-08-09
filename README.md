[![CircleCI](https://circleci.com/gh/jsturtevant/azure-k8-metrics-adapter.svg?style=svg)](https://circleci.com/gh/jsturtevant/azure-k8-metrics-adapter)
[![GitHub (pre-)release](https://img.shields.io/github/release/jsturtevant/azure-k8-metrics-adapter/all.svg)](https://github.com/jsturtevant/azure-k8-metrics-adapter/releases)

# Azure Kubernetes Metrics Adapter

An implementation of the Kubernetes [Custom Metrics API and External Metrics API](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/#support-for-metrics-apis) for Azure Services. 

This adapter enables you to scale your [application deployment pods](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/) running on [AKS](https://docs.microsoft.com/en-us/azure/aks/) using the [Horizontal Pod Autoscaler](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/) (HPA) with metrics from Azure Resources (such as [Service Bus Queues](https://docs.microsoft.com/en-us/azure/service-bus-messaging/service-bus-dotnet-get-started-with-queues)) and custom metrics stored in Application Insights. Learn more about [using an HPA to autoscale with external and custom metrics](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale-walkthrough/#autoscaling-on-metrics-not-related-to-kubernetes-objects).

Checkout a [video showing how scaling works with the adapter](https://www.youtube.com/watch?v=5pNpzwLLzW4&feature=youtu.be), [deploy the adapter](#deploy) or [learn by going through the walkthrough](samples/servicebus-queue/readme.md).

This was build using the [Custom Metric Adapter Server Boilerplate project.](https://github.com/kubernetes-incubator/custom-metrics-apiserver) 

## Project Status: Alpha

## Walkthrough
Check out this [walkthrough](samples/servicebus-queue/readme.md) to try it out.

## Deploy
Requires some [set up on your AKS Cluster](#azure-setup) and [Metric Server deployed](https://github.com/kubernetes-incubator/metrics-server#deployment) to your cluster.

```
kubectl apply -f https://raw.githubusercontent.com/jsturtevant/azure-k8-metrics-adapter/master/deploy/adapter.yaml
```

After deployment you can create an Horizontal Pod Auto Scaler (HPA) to scale of your [external metric](#external-metrics) of choice:

```yaml
apiVersion: autoscaling/v2beta1
kind: HorizontalPodAutoscaler
metadata:
 name: consumer-scaler
spec:
 scaleTargetRef:
   apiVersion: extensions/v1beta1
   kind: Deployment
   name: consumer
 minReplicas: 1
 maxReplicas: 10
 metrics:
  - type: External
    external:
      metricName: queuemessages
      metricSelector:
        matchLabels:
          metricName: Messages
          resourceGroup: sb-external-example
          resourceName: sb-external-ns
          resourceProviderNamespace: Microsoft.Servicebus
          resourceType: namespaces
          aggregation: Total
          filter: EntityName_eq_externalq
      targetValue: 30
```

And your that's it to enable auto scaling on External Metric.  Checkout the [samples](samples) for more examples.

### Verifying the deployment
You can also can also query the api to available metrics:

```bash
kubectl get --raw "/apis/custom.metrics.k8s.io/v1beta1" | jq .
kubectl get --raw "/apis/external.metrics.k8s.io/v1beta1" | jq .
```

To Query for a specific custom metric (*not currently supported*):

```
kubectl get --raw "/apis/custom.metrics.k8s.io/v1beta1/namespaces/test/pods/*/custom-metric" | jq .
```

To query for a specific external metric:

```bash
kubectl  get --raw "/apis/external.metrics.k8s.io/v1beta1/namespaces/test/queuemessages?labelSelector=resourceProviderNamespace=Microsoft.Servicebus,resourceType=namespaces,aggregation=Total,filter=EntityName_eq_externalq,resourceGroup=sb-external-example,resourceName=sb-external-ns,metricName=Messages" | jq .
```

## External Metrics

Requires k8 1.10+

See a full [list of hundreds of available azure external metrics](https://docs.microsoft.com/en-us/azure/monitoring-and-diagnostics/monitoring-supported-metrics) that can be used.  

Common external metrics to use for autoscaling are:

- [Azure ServiceBus Queue](https://docs.microsoft.com/en-us/azure/monitoring-and-diagnostics/monitoring-supported-metrics#microsoftservicebusnamespaces)  - Message Count - [example](samples/servicebus-queue)
- [Azure Storage Queue](https://docs.microsoft.com/en-us/azure/monitoring-and-diagnostics/monitoring-supported-metrics#microsoftstoragestorageaccountsqueueservices) - Message Count
- [Azure Eventhubs](https://docs.microsoft.com/en-us/azure/monitoring-and-diagnostics/monitoring-supported-metrics#microsofteventhubnamespaces)

## Custom Metrics
Custom Metrics are not currently implemented.

## Azure Setup

### Security
Authenticating with Azure Monitor can be achieved via a variety of authentication mechanisms. ([full list](https://github.com/Azure/azure-sdk-for-go#more-authentication-details))

We recommend to use one of the following options:
- [Azure Managed Service Identity](#using-azure-managed-service-identity-msi) (MSI)
- [Azure AD Application ID and Secret](#using-azure-ad-application-id-and-secret)
- [Azure AD Application ID and X.509 Certificate](#azure-ad-application-id-and-x509-certificate)

The Azure AD entity needs to have `Monitoring Reader` permission on the resource group that will be queried. More information can be found [here](https://docs.microsoft.com/en-us/azure/monitoring-and-diagnostics/monitoring-roles-permissions-security).

#### Using Azure Managed Service Identity (MSI)
Enable [Managed Service Identity](https://docs.microsoft.com/en-us/azure/active-directory/managed-service-identity/tutorial-linux-vm-access-arm) on each of your AKS vms: 

> There is a known issue when upgrading a AKS cluster with MSI enabled.  After the AKS upgrade you will lose your MSI setting and need to re-enable it. An alternative may be to use [aad-pod-identity](https://github.com/Azure/aad-pod-identity)

```bash
export RG=<aks resource group> 
export CLUSTER=<aks cluster name> 

NODE_RG="$(az aks show -n $CLUSTER -g $RG | jq -r .nodeResourceGroup)"
az vm list -g $NODE_RG
VMS="$(az vm list -g $NODE_RG | jq -r '.[] | select(.tags.creationSource | . and contains("aks")) | .name')"

while read -r vm; do
    echo "updating vm $vm..."
    msi="$(az vm identity assign -g $NODE_RG -n $vm | jq -r .systemAssignedIdentity)"
done <<< "$VMS"
```

Give access to the resource the MSI needs to access for each vm: 

```bash
export RG=<aks resource group> 
export CLUSTER=<aks cluster name> 
export ACCESS_RG=<resource group with metrics>

NODE_RG="$(az aks show -n $CLUSTER -g $RG | jq -r .nodeResourceGroup)"
az vm list -g $NODE_RG
VMS="$(az vm list -g $NODE_RG | jq -r '.[] | select(.tags.creationSource | . and contains("aks")) | .name')"

while read -r vm; do
    echo "getting vm identity $vm..."
    msi="$(az vm identity show -g $NODE_RG -n $vm | jq -r .principalId)"

    echo "adding access with msi $msi..."
    az role assignment create --role "Monitoring Reader" --assignee-object-id $msi --resource-group $ACCESS_RG
done <<< "$VMS"
```

#### Using Azure AD Application ID and Secret
See how to create an [example deployment](samples/azure-authentication).

Required environment variables:
- `AZURE_TENANT_ID`: Specifies the Tenant to which to authenticate.
- `AZURE_CLIENT_ID`: Specifies the app client ID to use.
- `AZURE_CLIENT_SECRET`: Specifies the app secret to use.
    
#### Azure AD Application ID and X.509 Certificate
Required environment variables:
- `AZURE_TENANT_ID`: Specifies the Tenant to which to authenticate.
- `AZURE_CLIENT_ID`: Specifies the app client ID to use.
- `AZURE_CERTIFICATE_PATH`: Specifies the certificate Path to use.
- `AZURE_CERTIFICATE_PASSWORD`: Specifies the certificate password to use.

## Subscription Information

The use the adapter your Azure Subscription must be provided.  There are a few ways to provide this information:

- [Azure Instance Metadata](https://docs.microsoft.com/en-us/azure/virtual-machines/windows/instance-metadata-service) - If you are running the adapter on a VM in Azure (for instance in an AKS cluster) there is nothing you need to do.  The Subscription Id will be automatically picked up from the Azure Instance Metadata endpoint
- Environment Variable - If you are outside of Azure or want full control of the subscription that is used you can set the Environment variable `SUBSCRIPTION_ID`  on the adapter deployment.  This takes precedence over the Azure Instance Metadata.
- [On each HPA](samples/hpa-examples) - you can work with multiple subscriptions by supplying the metric selector `subscriptionID` on each HPA.  This overrides Environment variables and Azure Instance Metadata settings.

## Contributing

See [Contributing](CONTRIBUTING.md) for more information.




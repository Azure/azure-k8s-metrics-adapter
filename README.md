[![CircleCI](https://circleci.com/gh/jsturtevant/azure-k8-metrics-adapter.svg?style=svg)](https://circleci.com/gh/jsturtevant/azure-k8-metrics-adapter)
[![GitHub (pre-)release](https://img.shields.io/github/release/jsturtevant/azure-k8-metrics-adapter/all.svg)](https://github.com/jsturtevant/azure-k8-metrics-adapter/releases)

# Azure Kubernetes Metrics Adapter

An implementation of the Kubernetes [Custom Metrics API and External Metrics API](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/#support-for-metrics-apis) for Azure Services. 

This adapter enables you to scale your [application deployment pods](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/) running on [AKS](https://docs.microsoft.com/en-us/azure/aks/) using the [Horizontal Pod Autoscaler](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/) (HPA) with metrics from Azure Resources (such as [Service Bus Queues](https://docs.microsoft.com/en-us/azure/service-bus-messaging/service-bus-dotnet-get-started-with-queues)) and custom metrics stored in Application Insights. Learn more about [using an HPA to autoscale with with external and custom metrics](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale-walkthrough/#autoscaling-on-metrics-not-related-to-kubernetes-objects).

Checkout a [video showing how scaling works with the adapter](https://www.youtube.com/watch?v=5pNpzwLLzW4&feature=youtu.be), [deploy the adapter](#deploy) or [learn by going through the walkthrough](samples/servicebus-queue/readme.md).

This was build using the [Custom Metric Adapter Server Boilerplate project.](https://github.com/kubernetes-incubator/custom-metrics-apiserver) 

## External Metrics

Requires k8 1.10+

See a full [list of hundreds of available azure external metrics](https://docs.microsoft.com/en-us/azure/monitoring-and-diagnostics/monitoring-supported-metrics) that can be used.  

Common external metrics to use for autoscaling are:

- [Azure ServiceBus Queue](https://docs.microsoft.com/en-us/azure/monitoring-and-diagnostics/monitoring-supported-metrics#microsoftservicebusnamespaces)  - Message Count - [example](samples/servicebus-queue)
- [Azure Storage Queue](https://docs.microsoft.com/en-us/azure/monitoring-and-diagnostics/monitoring-supported-metrics#microsoftstoragestorageaccountsqueueservices) - Message Count
- [Azure Eventhubs](https://docs.microsoft.com/en-us/azure/monitoring-and-diagnostics/monitoring-supported-metrics#microsofteventhubnamespaces)


## Custom Metrics
Custom Metrics are not currently implemented.

## Walkthrough
Check out this [walkthrough](samples/servicebus-queue/readme.md) to try it out.

## Deploy
Requires some [set up on your AKS Cluster](#azure-setup) and [Metric Server deployed](https://github.com/kubernetes-incubator/metrics-server#deployment) to your cluster.

```
kubectl apply -f https://raw.githubusercontent.com/jsturtevant/azure-k8-metrics-adapter/master/deploy/adapter.yaml
```

After deployment you can query the api to avaliable metrics:

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

## Azure Setup

### Security
Authenticating with Azure Monitor can be achieved via a variety of authentication mechanisms. ([full list](https://github.com/Azure/azure-sdk-for-go#more-authentication-details))

We recommend to use one of the following options:
- **Azure Managed Service Identity (MSI)**
- **Azure AD Application ID and Secret**
- **Azure AD Application ID and X.509 Certificate**

The Azure AD entity needs to have `Monitoring Reader` permission on the resource group that will be queried. More information can be found [here](https://docs.microsoft.com/en-us/azure/monitoring-and-diagnostics/monitoring-roles-permissions-security).

#### Using Azure Managed Service Identity (MSI)
Enable [Managed Service Identity](https://docs.microsoft.com/en-us/azure/active-directory/managed-service-identity/tutorial-linux-vm-access-arm) on each of your AKS vms: 

> There is a known issue when upgrading a AKS cluster with MSI enabled.  After the AKS upgrade you will lose your MSI setting and need to re-enable it.  

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
    az role assignment create --role Reader --assignee-object-id $msi --resource-group $ACCESS_RG
done <<< "$VMS"
```

## Development

### Get the source

```bash
go get github.com/jsturtevant/azure-k8-metrics-adapter
cd $GOPATH/github.com/jsturtevant/azure-k8-metrics-adapter
```

### Use Skaffold
Before you run the command below be sure to:

- Download [skaffold](https://github.com/GoogleContainerTools/skaffold#installation) 
- Log in to your container registry: `docker login`
- Have your K8 context set to the cluster you want to deploy to: `kubectl config use-context`

```bash
make dev
```


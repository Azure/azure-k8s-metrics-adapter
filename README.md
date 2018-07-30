[![CircleCI](https://circleci.com/gh/jsturtevant/azure-k8-metrics-adapter.svg?style=svg)](https://circleci.com/gh/jsturtevant/azure-k8-metrics-adapter)

# Azure Kuberenetes Metrics Adapter

An implementation of the Kubernetes [Custom Metrics API and External Metrics API](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/#support-for-metrics-apis) for Azure Services. 

This adapter enables you to scale your application [deployments](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/) running on [AKS](https://docs.microsoft.com/en-us/azure/aks/) using the [Horizontal Pod Autoscaler](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/) (HPA) with metrics from Azure Resources (such as [Service Bus Queues](https://docs.microsoft.com/en-us/azure/service-bus-messaging/service-bus-dotnet-get-started-with-queues)) and custom metrics stored in Application Insights. 

## External Metrics

Requires k8 1.10+

See a [list of available azure external metrics](https://docs.microsoft.com/en-us/azure/monitoring-and-diagnostics/monitoring-supported-metrics#microsofteventhubnamespaces) that can be used.  See section on [autoscaling with custom metrics](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale-walkthrough/#autoscaling-on-metrics-not-related-to-kubernetes-objects) to learn how to create an HPA to use external metrics.  A schema for external metrics can be found [here](https://raw.githubusercontent.com/kubernetes/kubernetes/master/api/openapi-spec/swagger.json).  

Common external metrics are:

- Azure ServiceBus Queue Message Count - [example](samples/servicebus-queue)
- Azure Storage Queue Message Count 
- Azure Eventhubs

## Custom Metrics
Custom Metrics are not currently implemented.

## Walkthrough
Check out this [walkthrough](samples/servicebus-queue) to try it out.

## Deploy
Requires some [set up on your AKS Cluster](#azure-setup)

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
kubectl  get --raw "/apis/external.metrics.k8s.io/v1beta1/namespaces/test/queuemessages?labelSelector=resourceProviderNamespace=Microsoft.Servicebus,resourceType=namespaces,aggregation=Total,filter=EntityName_eq_helloworld,resourceName=k8custom,resourceGroup=k8metrics,resourceName=k8custom,metricName=Messages" | jq .
```

## Azure Setup

Enable [Managed Service Identity](https://docs.microsoft.com/en-us/azure/active-directory/managed-service-identity/tutorial-linux-vm-access-arm) on each of your AKS vms and give access to the resource the MSI access for each vm:

```bash
export RG=<aks resource group> 
export CLUSTER=<aks cluster name> 
export ACCESS_RG=<rg to give read access to>

NODE_RG="$(az aks show -n $CLUSTER -g $RG | jq -r .nodeResourceGroup)"
az vm list -g $NODE_RG
VMS="$(az vm list -g $NODE_RG | jq -r '.[] | select(.tags.creationSource | . and contains("aks")) | .name')"

while read -r vm; do
    echo "updating vm $vm..."
    msi="$(az vm identity assign -g $NODE_RG -n $vm | jq -r .systemAssignedIdentity)"

    echo "adding access with $msi..."
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


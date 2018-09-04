[![CircleCI](https://circleci.com/gh/Azure/azure-k8s-metrics-adapter/tree/master.svg?style=svg)](https://circleci.com/gh/Azure/azure-k8s-metrics-adapter/tree/master)
[![GitHub (pre-)release](https://img.shields.io/github/release/Azure/azure-k8s-metrics-adapter/all.svg)](https://github.com/Azure/azure-k8s-metrics-adapter/releases)

# Azure Kubernetes Metrics Adapter

An implementation of the Kubernetes [Custom Metrics API and External Metrics API](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/#support-for-metrics-apis) for Azure Services. 

This adapter enables you to scale your [application deployment pods](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/) running on [AKS](https://docs.microsoft.com/en-us/azure/aks/) using the [Horizontal Pod Autoscaler](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/) (HPA) with [External Metrics](#external-metrics) from Azure Resources (such as [Service Bus Queues](https://docs.microsoft.com/en-us/azure/service-bus-messaging/service-bus-dotnet-get-started-with-queues)) and [Custom Metrics](#custom-metrics) stored in Application Insights. 

Learn more about [using an HPA to autoscale with external and custom metrics](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale-walkthrough/#autoscaling-on-metrics-not-related-to-kubernetes-objects).

Checkout a [video showing how scaling works with the adapter](https://www.youtube.com/watch?v=5pNpzwLLzW4&feature=youtu.be), [deploy the adapter](#deploy) or [learn by going through the walkthrough](samples/servicebus-queue/readme.md).

This was build using the [Custom Metric Adapter Server Boilerplate project.](https://github.com/kubernetes-incubator/custom-metrics-apiserver) 

## Project Status: Alpha

## Walkthrough
Try out scaling with [External Metrics](#external-metrics) using the a Azure Service Bus Queue in this [walkthrough](samples/servicebus-queue).

Try out scaling with [Custom Metrics](#custom-metrics) using Requests per Second and Application Insights in this [walkthrough](samples/request-per-second) 

## Deploy
Requires some [set up on your AKS Cluster](#azure-setup) and [Metric Server deployed](https://github.com/kubernetes-incubator/metrics-server#deployment) to your cluster.

```
kubectl apply -f https://raw.githubusercontent.com/Azure/azure-k8s-metrics-adapter/master/deploy/adapter.yaml
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

Custom metrics are currently retrieved from Application Insights.  Currently you will need to use an [AppId and API key](https://dev.applicationinsights.io/documentation/Authorization/API-key-and-App-ID) to authenticate and retrieve metrics from Application Insights.  View a list of basic metrics that come out of the box and see sample values at the [AI api explorer](https://dev.applicationinsights.io/apiexplorer/metrics).  

> note: Metrics that have multi parts currently need will be passed as with a separator of `-` in place of AI separator of `/`.  An example is `performanceCounters/requestsPerSecond` should be specified as `performanceCounters-requestsPerSecond`

Common Custom Metrics are:

- Requests per Second (RPS) - [example](samples/request-per-second) 

## Azure Setup

### Security

Authenticating with Azure Monitor can be achieved via a variety of authentication mechanisms. ([full list](https://github.com/Azure/azure-sdk-for-go#more-authentication-details))

We recommend to use one of the following options:

- [Azure Managed Service Identity](#using-azure-managed-service-identity-msi) (MSI)
- [Azure AD Pod Identity](#using-azure-ad-pod-identity) (aad-pod-identity)
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

#### Using Azure AD Pod Identity

[aad-pod-identity](https://github.com/Azure/aad-pod-identity) is currently in beta and allows to bind a user managed identity or a service principal to a pod. That means that instead of using the same managed identity for all the pod running on a node like explained above, you are able to get a specific identity with specific RBAC for a specific pod.

Using this project requires to deploy a bit of infrastructure first. You can do it following the [Get started page](https://github.com/Azure/aad-pod-identity#get-started) of the project.

Once the aad-pod-identity infrastructure is running, you need to create an Azure identity scoped to the resource group you are monitoring:

```bash
az identity create -g {ResourceGroup1} -n k8s-custom-metrics-identity
```

Assign `Monitoring Reader` to it

```bash
az role assignment create --role "Monitoring Reader" --assignee <principalId> --scopes /subscriptions/{SubID}/resourceGroups/{ResourceGroup1}
```

As documented [here](https://github.com/Azure/aad-pod-identity#providing-required-permissions-for-mic) *aad-pod-identity* uses the service principal of  your Kubernetes cluster to access the Azure resources. You need to give this service principal the rights to use the managed identity created before:

```bash
az role assignment create --role "Managed Identity Operator" --assignee <servicePrincipalId> --scope /subscriptions/{SubID}/resourceGroups/{ResourceGroup1}/providers/Microsoft.ManagedIdentity/userAssignedIdentities/k8s-custom-metrics-identity
```

Install the Azure Identity to your Kubernetes cluster:

```yaml
apiVersion: "aadpodidentity.k8s.io/v1"
kind: AzureIdentity
metadata:
  name: k8s-custom-metrics-identity
  namespace: custom-metrics
spec:
  type: 0
  ResourceID: /subscriptions/{SubID}/resourceGroups/{ResourceGroup1}/providers/Microsoft.ManagedIdentity/userAssignedIdentities/k8s-custom-metrics-identity
  ClientID: <clientid>
```

Install the Azure Identity Binding on your Kubernetes cluster:

```yaml
apiVersion: "aadpodidentity.k8s.io/v1"
kind: AzureIdentityBinding
metadata:
  name: k8s-custom-metrics-identity-binding
spec:
  AzureIdentity: k8s-custom-metrics-identity
  Selector: k8s-custom-metrics-identity
```

#### Using Azure AD Application ID and Secret

See how to create an [example deployment](samples/azure-authentication).

Create a service principal scoped to the resource group the resource you monitoring and assign  `Monitoring Reader` to it:

```
az ad sp create-for-rbac -n "adapter-sp" --role "Monitoring Reader" --scopes /subscriptions/{SubID}/resourceGroups/{ResourceGroup1}
```

Required environment variables:
- `AZURE_TENANT_ID`: Specifies the Tenant to which to authenticate.
- `AZURE_CLIENT_ID`: Specifies the app client ID to use.
- `AZURE_CLIENT_SECRET`: Specifies the app secret to use.

Deploy the environment variables via secret: 

```
 kubectl create secret generic adapter-service-principal -n custom-metrics \
  --from-literal=azure-tenant-id=<tenantid> \
  --from-literal=azure-client-id=<clientid>  \
  --from-literal=azure-client-secret=<secret>
 ```
    
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

## Issues
Report any issues in the [Github issues](https://github.com/Azure/azure-k8s-metrics-adapter/issues).  

## Roadmap
See the Projects tab for [current roadmap](https://github.com/Azure/azure-k8s-metrics-adapter/projects).

## Reporting Security Issues

Security issues and bugs should be reported privately, via email, to the Microsoft Security
Response Center (MSRC) at [secure@microsoft.com](mailto:secure@microsoft.com). You should
receive a response within 24 hours. If for some reason you do not, please follow up via
email to ensure we received your original message. Further information, including the
[MSRC PGP](https://technet.microsoft.com/en-us/security/dn606155) key, can be found in
the [Security TechCenter](https://technet.microsoft.com/en-us/security/default).
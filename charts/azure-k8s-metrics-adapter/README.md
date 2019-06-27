# Azure Kubernetes Metrics Adapter

An implementation of the Kubernetes [Custom Metrics API and External Metrics API](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/#support-for-metrics-apis) for Azure Services. 

This adapter enables you to scale your [application deployment pods](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/) running on [AKS](https://docs.microsoft.com/en-us/azure/aks/) using the [Horizontal Pod Autoscaler](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/) (HPA) with metrics from Azure Resources (such as [Service Bus Queues](https://docs.microsoft.com/en-us/azure/service-bus-messaging/service-bus-dotnet-get-started-with-queues)) and custom metrics stored in Application Insights. Learn more about [using an HPA to autoscale with with external and custom metrics](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale-walkthrough/#autoscaling-on-metrics-not-related-to-kubernetes-objects).

## Installing the Chart
Clone this repository and cd to the root folder:  

```
go get -u github.com/Azure/azure-k8s-metrics-adapter
cd $GOPATH/src/github.com/Azure/azure-k8s-metrics-adapter/
```

> If you are not planning on modify the project and don't have Go installed, you can clone this via the `git clone https://github.com/Azure/azure-k8s-metrics-adapter.git` and `cd` to the root of the project.

Next create a namespace and install:

```sh
kubectl create namespace custom-metrics
helm install --name my-release charts/azure-k8s-metrics-adapter --namespace custom-metrics
```

## Uninstalling the Chart

```sh
helm delete my-release
```

## Configuration

### All Available Configuration Values

| Parameter | Description | Default |
| --------- | ----------- | ------- |
| `image.repository` | Image repository | `jsturtevant/azure-k8-metrics-adapter-amd64` |
| `image.tag` | Image tag | `latest` |
| `image.pullPolicy` | Image pull policy | `Always` |
| `logLevel` | Log level for V logs | `2` |
| `replicaCount`  | Number of azure-k8s-metrics-adapter replicas  | `1` |
| `adapterSecurePort` | Port on which the adapter is listening | `6443` |
| `apiServiceInsecureSkipTLSVerify` | Disables TLS certificate verification when communicating with the apiService | `true` |
| `apiServiceGroupPriorityMinimum` | The priority the APIService group should have at least | `100` |
| `apiServiceVersionPriority` | Controls the ordering of this API version inside of its group | `100` |
| `azureAuthentication.method` | Which method to use when authenticating with the Azure Resource Monitory service. Valid options are `msi`, `clientSecret`, and `clientCertificate` | `msi` |
| `azureAuthentication.clientID` | The Azure Active Directory Application client ID to use if using the authentication method `clientSecret` or `clientCertificate` | `''` |
| `azureAuthentication.tenantID` | Specifies the Tenant to which to authenticate if using the authentication method `clientSecret` or `clientCertificate | `''` |
| `azureAuthentication.clientSecret` | Specifies the app secret to use if using the authentication method `clientSecret` | `''` |
| `azureAuthentication.clientCertificate` | Specifies the contents of the certificate if using the authentication method `clientCertificate`  | `''` |
| `azureAuthentication.clientCertificatePath` | Specifies certificate Path to use if using the authentication method `clientCertificate`  | `''` |
| `azureAuthentication.azureClientCertificatePassword` | Specifies the certificate password to use  if using the authentication method `clientCertificate`  | `''` |
| `defaultSubscriptionId` | Specifies the subscription to use instead of using Azure Instance Metadata  | `''` |
| `extraArgs` | Optional flags for azure-k8s-metrics-adapter | `{}` |
| `extraEnv` | Optional environment variables for azure-k8s-metrics-adapter | `{}` |
| `rbac.create` | If `true`, create and use RBAC resources | `true` |
| `serviceAccount.create` | If `true`, create a new service account | `true` |
| `serviceAccount.name` | Service account to be used. If not set and `serviceAccount.create` is `true`, a name is generated using the fullname template |  |
| `resources` | CPU/memory resource requests/limits | `requests: {cpu: 10m, memory: 32Mi}` |
| `nodeSelector` | Node labels for pod assignment | `{}` |
| `affinity` | Node affinity for pod assignment | `{}` |
| `tolerations` | Node tolerations for pod assignment | `[]` |

## Authentication to Azure Montior

This project offers 3 methods for authenticating to the Azure Monitor API: [MSI](https://github.com/Azure/azure-k8s-metrics-adapter#using-azure-ad-pod-identity), Azure AD Application with Client Secret, and Azure AD Application with Client Certificate. By default this chart will use MSI authentication.

Optionally the `azureAuthentication.method` value can be specified and either Azure AD Application with Client Secret or Azure AD Application with Client Certificate can be set using the following values for `azureAuthentication.method`
- `clientSecret` for Azure AD Application with Client Secre. These additional values must be set to use the client certificate
    - `azureAuthentication.clientID`
    - `azureAuthentication.tenantID`
    - `azureAuthentication.clientSecret`
- `clientCertificate` Azure AD Application with Client Certificate. These additional values must be set to use the client certificate
    - `azureAuthentication.clientID`
    - `azureAuthentication.tenantID`
    - `azureAuthentication.clientCertificate`
    - `azureAuthentication.clientCertificatePath`
    - `azureAuthentication.clientCertificatePassword`

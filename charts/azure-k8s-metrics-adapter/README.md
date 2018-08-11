# Azure Kubernetes Metrics Adapter

An implementation of the Kubernetes [Custom Metrics API and External Metrics API](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/#support-for-metrics-apis) for Azure Services. 

This adapter enables you to scale your [application deployment pods](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/) running on [AKS](https://docs.microsoft.com/en-us/azure/aks/) using the [Horizontal Pod Autoscaler](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/) (HPA) with metrics from Azure Resources (such as [Service Bus Queues](https://docs.microsoft.com/en-us/azure/service-bus-messaging/service-bus-dotnet-get-started-with-queues)) and custom metrics stored in Application Insights. Learn more about [using an HPA to autoscale with with external and custom metrics](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale-walkthrough/#autoscaling-on-metrics-not-related-to-kubernetes-objects).

## Installing the Chart

```sh
helm install --name my-release charts/azure-k8s-metrics-adapter
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
| `azureTenantId` | Specifies the Tenant to which to authenticate | `''`
| `azureClientId` | Specifies the app client ID to use | `''`
| `azureClientSecret` | Specifies the app secret to use | `''`
| `azureClientCertificate` | Specifies the contents of the certificate to use | `''`
| `azureClientCertificatePath` | Specifies the certificate Path to use | `''`
| `azureClientCertificatePassword` | Specifies the certificate password to use | `''`
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

### Authenticating with an Azure AD Application ID and Secret

To authenticate with an Azure AD Application ID and Secret set the following configuration values
- `azureTenantId`
- `azureClientId`
- `azureClientSecret`

### Authenticating with an Azure AD Application ID and X.509 Certificate

To authenticate with an Azure AD Application ID and X.509 Certificate set the following configuration values
- `azureTenantId`
- `azureClientId`
- `azureClientCertificate`
- `azureClientCertificatePath`
- `azureClientCertificatePassword`

# Integration testing against different Kubernetes versions

**TODO: edit repo links based on final location of this README**

## Prerequisites

In order to run make teste2e

* a Kubernetes cluster that your kubectl context points to
* [Helm](https://docs.helm.sh/using_helm/) installed locally and on your cluster (or [Helm for RBAC-enabled AKS clusters](https://docs.microsoft.com/en-us/azure/aks/kubernetes-helm))
* jq (used in parsing responses from the endpoint)
* Docker
* [Kubernetes Metrics Server](https://github.com/kubernetes-incubator/metrics-server#deployment) deployed on your cluster (it is deployed by default with most deployments)

If testing locally on Minikube, you may find you need [socat](.azdevops/0_install/install-misc.sh). 
If testing locally on Minikube using Kubernetes 1.11, you may find you require [crictl](.azdevops/0_install/install-crictl.sh) and [ebtables](.azdevops/0_install/install-misc.sh).

## Environment variables

Edit the [local dev values](local-dev-values.yaml.example) file to create `local-dev-values.yaml`. 

**Currently, these scripts only support deployment with a Service Principal.** To use an alternate method (MSI, Azure AD Pod Identity), check the [walkthrough for the Service Bus sample](https://github.com/Azure/azure-k8s-metrics-adapter/tree/master/samples/servicebus-queue#enable-access-to-azure-resources) for differences. You should just need to alter the values file and the adapter deployment step.

| Variable name | Description |  Optional? |
| ------------- | ----------- |  --------- |
| `SERVICEBUS_CONNECTION_STRING` | Connection string for the service bus namespace | No |
| `SERVICEBUS_RESOURCE_GROUP` | Resource group that holds the service bus namespace | No |
| `SERVICEBUS_NAMESPACE` | Service bus namespace | No |
| `SERVICEBUS_QUEUE_NAME` | Name of the service bus queue | Yes, defaults to `externalq` if not set |
| `GOPATH` | Golang project directory | Yes, defaults to `$HOME/go` if not set |


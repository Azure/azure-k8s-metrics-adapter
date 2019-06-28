# Contributing

This project welcomes contributions and suggestions. Most contributions require you to
agree to a Contributor License Agreement (CLA) declaring that you have the right to,
and actually do, grant us the rights to use your contribution. For details, visit
https://cla.microsoft.com.

When you submit a pull request, a CLA-bot will automatically determine whether you need
to provide a CLA and decorate the PR appropriately (e.g., label, comment). Simply follow the
instructions provided by the bot. You will only need to do this once across all repositories using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/)
or contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.

## Development
To do development you will need:

- [Golang](https://golang.org/doc/install) - same as current [Kubernetes version ](https://github.com/kubernetes/community/blob/master/contributors/devel/development.md#go)
- Kubernetes cluster - [minikube](https://github.com/kubernetes/minikube), [Docker for Mac with Kubernetes support](https://docs.docker.com/docker-for-mac/kubernetes/),  [Docker for Windows with Kubernetes support](https://docs.docker.com/docker-for-windows/kubernetes/), [AKS](https://docs.microsoft.com/en-us/azure/aks/kubernetes-walkthrough)
- [git](https://git-scm.com/downloads) 
- [mercurial](https://www.mercurial-scm.org/downloads)  

### Get the source

```bash
go get github.com/Azure/azure-k8s-metrics-adapter
cd $GOPATH/src/github.com/Azure/azure-k8s-metrics-adapter
```

### Add your fork

[Fork this project in GitHub](https://help.github.com/articles/fork-a-repo/). Then add your fork:

```bash
cd $GOPATH/github.com/Azure/azure-k8s-metrics-adapter
git remote rename origin upstream #rename to upstream so you can sync 
git remote add origin <your-fork-url>
git checkout -b <your-feature-branch>
```

Renaming the `origin` set by `go get` to `upstream` let's you use [upstream to sync your repository](https://help.github.com/articles/syncing-a-fork/) so you can keep your project uptodate with changes.

### Building the project
To build the project locally, with out creating a docker image:

```bash
make build-local
```

To build the docker image use:

```bash
export REGISTRY=<your registry name> ("" if using DockerHub)
export IMAGE=azure-k8s-metrics-adapter-testimage
make build
```

You can then use `make push`:

```bash
export DOCKER_USER=<your docker username>
(optional, will prompt otherwise) export DOCKER_PASS=<your docker password>
make push
```

### End-to-end testing
You can run `make teste2e` to check that the adapter deploys properly, uses given metrics, and pulls metric information. This script uses the [Service Bus Queue example](samples/servicebus-queue/readme.md).

To run `make teste2e`, you need the following:

* [Helm](https://docs.helm.sh/using_helm/) installed locally and on your cluster (or [Helm for RBAC-enabled AKS clusters](https://docs.microsoft.com/en-us/azure/aks/kubernetes-helm))
* jq (used in parsing responses from the endpoint)
* [Kubernetes Metrics Server](https://github.com/kubernetes-incubator/metrics-server#deployment) deployed on your cluster (it is deployed by default with most deployments)
* An Azure [Service Bus Queue](https://docs.microsoft.com/en-us/azure/service-bus-messaging/service-bus-dotnet-get-started-with-queues)
* An Azure Topic with Subscription

#### Environment variables for e2e tests

Build the project with a custom repository:

```
make build
```

Edit the [local dev values](local-dev-values.yaml.example) file to create `local-dev-values.yaml`. If using custom image be sure to set the values (`export REGISTRY=<your registry name>` ("" if using DockerHub)
`export IMAGE=azure-k8s-metrics-adapter-testimage`) before building.  The `pullPolicy: IfNotPresent` lets you use the local image on your minikube cluster.  If you are not using a local cluster you can use `pullPolicy: Always` to use an image that is in a remote repository. 

Example of the `image` setting in the `local-dev-values.yaml` using a custom image:

```
image:
  repository: metrics-adapter
  tag: latest
  pullPolicy: IfNotPresent
```

Set the following Environment Variables:

| Variable name | Description |  Optional? |
| ------------- | ----------- |  --------- |
| `SERVICEBUS_CONNECTION_STRING` | Connection string for the service bus namespace | No |
| `SERVICEBUS_RESOURCE_GROUP` | Resource group that holds the service bus namespace | No |
| `SERVICEBUS_NAMESPACE` | Service bus namespace | No |
| `SERVICEBUS_QUEUE_NAME` | Name of the service bus queue | Yes, defaults to `externalq` if not set |
| `GOPATH` | Golang project directory | Yes, defaults to `$HOME/go` if not set |
| `SERVICEBUS_TOPIC_NAME` | Name of the service bus topic | Yes, defaults to `example-topic` if not set |
| `SERVICEBUS_SUBSCRIPTION_NAME` | Name of the service bus subscription |  Yes, defaults to `externalsub` if not set |

## Adding dependencies

Add the dependency to the Gopkg.toml file and then run:

```
make vendor
```

### Use Skaffold
To create a fast dev cycle you can use skaffold with a local cluster (minikube or Docker for win/mac).  Before you run the command below be sure to:

- Download [skaffold](https://github.com/GoogleContainerTools/skaffold#installation) 
- Have your K8s context set to the local cluster you want to deploy to: `kubectl config use-context`
- If using minikube run `eval $(minikube docker-env)`
- Create a Service Principle for local development: `az ad sp create-for-rbac -n "adapter-sp" --role "Monitoring Reader" --scopes /subscriptions/{SubID}/resourceGroups/{ResourceGroup1}` where the resource group contains resources (queue or app insights) you want to retrieve metrics for
- Make a copy of `local-dev-values.yaml.example` and call it `local-dev-values.yaml` (`cp local-dev-values.yaml.example local-dev-values.yaml`) and replace the values with your Service Principle and subscription id.  

Then run: 

```bash
make dev
```

## Releasing

1. Switch to the `master` branch and run `make version SEMVER=<sem-version-to-bump>`. Options for SEMVER are `SEMVER=major`, `SEMVER=minor` or `SEMVER=patch`
2. Then run `git push --follow-tags`
3. Everything is automated after the `git push`.  `make version` will bump the version and tag the commit.  The Circle CI will recognize the tagged master branch and push to the repository.

> note: you must be on the master branch and it must be clean. 

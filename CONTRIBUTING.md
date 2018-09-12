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
- Kubernetes cluster - [minikube](https://github.com/kubernetes/minikube), , [Docker for Mac with Kubernetes support](https://docs.docker.com/docker-for-mac/kubernetes/),  [Docker for Windows with Kubernetes support](https://docs.docker.com/docker-for-windows/kubernetes/), [AKS](https://docs.microsoft.com/en-us/azure/aks/kubernetes-walkthrough)
- [git](https://git-scm.com/downloads) 
- [mercurial](https://www.mercurial-scm.org/downloads)  

### Get the source

```bash
go get github.com/Azure/azure-k8s-metrics-adapter
cd $GOPATH/src/github.com/Azure/azure-k8s-metrics-adapter
```

### Building the project
To build the project locally, with out creating a docker image:

```bash
make build-local
```

To build the docker image use:

```bash
export REGISTRY=<your registry name>
export IMAGE=azure-k8s-metrics-adapter-testimage
make build
```

You can then login into your registry (`docker login`) and run:

```bash
make push
```

## Adding dependencies

Add the dependency to the Gopkg.toml file and then run:

```
make vendor
```

### Use Skaffold
To create a fast dev cycle you can use skaffold.  Before you run the command below be sure to:

- Download [skaffold](https://github.com/GoogleContainerTools/skaffold#installation) 
- Log in to your container registry: `docker login`
- Have your K8 context set to the cluster you want to deploy to: `kubectl config use-context`

```bash
make dev
```

## Releasing

1. Switch to the `master` branch and run `make version SEMVER=<sem-version-to-bump>`. Options for SEMVER are `SEMVER=major`, `SEMVER=minor` or `SEMVER=patch`
2. Then run `git push --follow-tags`
3. Everything is automated after the `git push`.  `make version` will bump the version and tag the commit.  The Circle CI will recognize the tagged master branch and push to the repository.

> note: you must be on the master branch and it must be clean. 
# Integration testing against different Kubernetes versions

## Prerequisites

In order to run the following, you need:

* a Kubernetes cluster that your kubectl context points to
* [Helm](https://docs.helm.sh/using_helm/) installed locally and on your cluster
* jq (used in parsing responses from the endpoint)
* Docker
* Go

If testing locally on Minikube, you may find you need [socat](.azdevops/0_install/install-misc.sh). 
If additionally using Kubernetes 1.11, you may find you require [crictl](.azdevops/0_install/install-crictl.sh) and [ebtables](.azdevops/0_install/install-misc.sh).


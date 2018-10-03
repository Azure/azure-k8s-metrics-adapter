#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

installminikube() {
    # from minikube docs: https://github.com/kubernetes/minikube
    curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64 && chmod +x minikube && sudo cp minikube /usr/local/bin/ && rm minikube
    curl -Lo kubectl https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl && chmod +x kubectl && sudo cp kubectl /usr/local/bin/ && rm kubectl

    export MINIKUBE_WANTUPDATENOTIFICATION=false
    export MINIKUBE_WANTREPORTERRORPROMPT=false
    export MINIKUBE_HOME=$HOME
    export CHANGE_MINIKUBE_NONE_USER=true
    mkdir -p $HOME/.kube
    mkdir -p $HOME/.minikube
    touch $HOME/.kube/config

    export KUBECONFIG=$HOME/.kube/config
    sudo -E minikube start --vm-driver=none

    # this for loop waits until kubectl can access the api server that Minikube has created
    for i in {1..150}; do # timeout for 5 minutes
    kubectl get po &> /dev/null
    if [ $? -ne 1 ]; then
        break
    fi
    sleep 2
    done
}

installhelm(){
    ./install-helm.sh
    helm init --wait
}

installMetricAdapter() {
    helm install --name e2e-test ../charts/azure-k8s-metrics-adapter \
    --namespace custom-metrics \
    --set azureAuthentication.method=clientSecret \
    --set azureAuthentication.tenantID=$TENANTID \
    --set azureAuthentication.clientID=$CLIENTID \
    --set azureAuthentication.clientSecret=$CLIENTSECRET \
    --set azureAuthentication.createSecret=true \
    --set image.repository=$FULL_IMAGE \
    --set image.tag=$VERSION
}

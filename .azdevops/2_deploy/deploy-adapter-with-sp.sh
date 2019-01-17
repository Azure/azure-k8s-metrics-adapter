#!/bin/bash

set -o nounset
set -o errexit

GOPATH="${GOPATH:-$HOME/go}"

echo; echo "Deploying metrics adapter..."
cd $GOPATH/src/github.com/Azure/azure-k8s-metrics-adapter
helm install --name adapter \
    -f ./local-values.yaml \
    ./charts/azure-k8s-metrics-adapter

echo; echo "Waiting for deployment to be available..."
until [[ `kubectl get deploy/adapter-azure-k8s-metrics-adapter -o jsonpath="{@.status.availableReplicas}"` == 1 ]]; do 
    kubectl get deploy/adapter-azure-k8s-metrics-adapter
    sleep 15
done
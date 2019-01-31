#!/bin/bash

set -o nounset
set -o errexit

GOPATH="${GOPATH:-$HOME/go}"

echo; echo "Deploying metrics adapter..."
cd $GOPATH/src/github.com/Azure/azure-k8s-metrics-adapter
helm install --name adapter \
    -f ./local-dev-values.yaml \
    ./charts/azure-k8s-metrics-adapter

echo; echo "Waiting for deployment to be available..."
START=`date +%s`

kubectl get deploy/adapter-azure-k8s-metrics-adapter
while [[ ! `kubectl get deploy/adapter-azure-k8s-metrics-adapter -o jsonpath="{@.status.availableReplicas}"` == 1 ]] && \
        [[ $(( $(date +%s) - 105 )) -lt $START ]]; do 
    sleep 15
    kubectl get deploy/adapter-azure-k8s-metrics-adapter --no-headers
done

if [[ ! `kubectl get deploy/adapter-azure-k8s-metrics-adapter -o jsonpath="{@.status.availableReplicas}"` == 1 ]]; then
    kubectl describe deploy/adapter-azure-k8s-metrics-adapter
    echo
    kubectl logs deploy/adapter-azure-k8s-metrics-adapter
fi
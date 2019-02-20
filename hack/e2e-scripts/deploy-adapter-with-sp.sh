#!/bin/bash

set -o nounset
set -o errexit

GOPATH="${GOPATH:-$HOME/go}"

echo; echo "Deploying metrics adapter..."
cd $GOPATH/src/github.com/Azure/azure-k8s-metrics-adapter
# If SP values are set via env vars, use them; otherwise, just use values file
if [[ -v SP_CLIENT_ID ]]; then
    helm install --name adapter \
        --set azureAuthentication.tenantID=$SP_TENANT_ID \
        --set azureAuthentication.clientID=$SP_CLIENT_ID \
        --set azureAuthentication.clientSecret=$SP_CLIENT_SECRET \
        -f ./local-dev-values.yaml \
        ./charts/azure-k8s-metrics-adapter
else
    helm install --name adapter \
        -f ./local-dev-values.yaml \
        ./charts/azure-k8s-metrics-adapter
fi

echo; echo "Waiting for deployment to be available..."
START=`date +%s`

kubectl get deploy/adapter-azure-k8s-metrics-adapter
while [[ ! `kubectl get deploy/adapter-azure-k8s-metrics-adapter -o jsonpath="{@.status.availableReplicas}"` = 1 ]] && \
        [[ $(( $(date +%s) - 55 )) -lt $START ]]; do 
    sleep 5
    kubectl get deploy/adapter-azure-k8s-metrics-adapter --no-headers
done

if [[ ! `kubectl get deploy/adapter-azure-k8s-metrics-adapter -o jsonpath="{@.status.availableReplicas}"` = 1 ]]; then
    echo; echo "Deployment failed, output debug information"
    kubectl describe deploy/adapter-azure-k8s-metrics-adapter
    echo
    kubectl logs deploy/adapter-azure-k8s-metrics-adapter
    exit 1
fi
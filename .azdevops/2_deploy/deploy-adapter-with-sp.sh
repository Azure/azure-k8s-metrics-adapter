#!/bin/bash

set -o nounset
set -o errexit

GOPATH="${GOPATH:-$HOME/go}"

echo; echo "Making image..."
cd $GOPATH/src/github.com/Azure/azure-k8s-metrics-adapter/
export REGISTRY="integration"
export REGISTRY_PATH=""
export VERSION="local"
make build-simple

echo; echo "Deploying metrics adapter..."
helm install --name adapter \
    ./charts/azure-k8s-metrics-adapter \
    --set azureAuthentication.method=clientSecret \
    --set azureAuthentication.tenantID=$SP_TENANT_ID \
    --set azureAuthentication.clientID=$SP_CLIENT_ID \
    --set azureAuthentication.clientSecret=$SP_CLIENT_SECRET \
    --set azureAuthentication.createSecret=true \
    --set defaultSubscriptionId=$SUBSCRIPTION_ID \
    --set image.repository=integration/adapter \
    --set image.tag=local \
    --set image.pullPolicy=IfNotPresent

echo; echo "Waiting for deployment to be available..."
until [[ `kubectl get deploy/adapter-azure-k8s-metrics-adapter -o jsonpath="{@.status.availableReplicas}"` == 1 ]]; do 
    kubectl get deploy/adapter-azure-k8s-metrics-adapter
    sleep 15
done
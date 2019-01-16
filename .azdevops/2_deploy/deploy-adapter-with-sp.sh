#!/bin/bash

set -o nounset
set -o errexit

GOPATH="${GOPATH:-$HOME/go}"

echo; echo "Deploying metrics adapter..."
helm install --name adapter \
    ./charts/azure-k8s-metrics-adapter \
    --set azureAuthentication.method=clientSecret \
    --set azureAuthentication.tenantID=$SP_TENANT_ID \
    --set azureAuthentication.clientID=$SP_CLIENT_ID \
    --set azureAuthentication.clientSecret=$SP_CLIENT_SECRET \
    --set azureAuthentication.createSecret=true \
    --set defaultSubscriptionId=$SUBSCRIPTION_ID \
    --set image.repository="$ACR_ADDR/$ACR_PATH" \
    --set image.tag=$IMAGE_BUILDNUMBER \
    --set image.pullPolicy=IfNotPresent

echo; echo "Waiting for deployment to be available..."
until [[ `kubectl get deploy/adapter-azure-k8s-metrics-adapter -o jsonpath="{@.status.availableReplicas}"` == 1 ]]; do 
    kubectl get deploy/adapter-azure-k8s-metrics-adapter
    sleep 15
done
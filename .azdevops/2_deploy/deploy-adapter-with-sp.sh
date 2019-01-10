#!/bin/bash

: "${SP_TENANT_ID:?Must set SP_TENANT_ID}"
: "${SP_CLIENT_ID:?Must set SP_CLIENT_ID}"
: "${SP_CLIENT_SECRET:?Must set SP_CLIENT_SECRET}"
: "${SUBSCRIPTION_ID:?Must set SUBSCRIPTION_ID}"

echo; echo "Making image..."
cd $HOME/go/src/github.com/Azure/azure-k8s-metrics-adapter/
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

until kubectl get deploy/adapter-azure-k8s-metrics-adapter -o jsonpath="{@.status.availableReplicas}"
    do sleep 15
    echo "waiting for deployment to be available"
done
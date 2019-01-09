#!/bin/bash

: "${SP_TENANT_ID:?Must set SP_TENANT_ID}"
: "${SP_CLIENT_ID:?Must set SP_CLIENT_ID}"
: "${SP_CLIENT_SECRET:?Must set SP_CLIENT_SECRET}"
: "${SUBSCRIPTION_ID:?Must set SUBSCRIPTION_ID}"

echo; echo "Deploying metrics adapter..."
cd $HOME/go/src/github.com/Azure/azure-k8s-metrics-adapter/
kubectl create secret generic subscriptionid --from-literal=subscription_id=$SUBSCRIPTION_ID
helm install --name adapter \
    ./charts/azure-k8s-metrics-adapter \
    --set azureAuthentication.method=clientSecret \
    --set azureAuthentication.tenantID=$SP_TENANT_ID \
    --set azureAuthentication.clientID=$SP_CLIENT_ID \
    --set azureAuthentication.clientSecret=$SP_CLIENT_SECRET \
    --set azureAuthentication.createSecret=true

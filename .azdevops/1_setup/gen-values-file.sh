#!/bin/bash

set -o nounset
set -o errexit

GOPATH="${GOPATH:-$HOME/go}"
FNAME="$GOPATH/src/github.com/Azure/azure-k8s-metrics-adapter/local-values.yaml"

echo; echo "Creating local values file..."

echo "azureAuthentication:" > $FNAME
echo "  method: clientSecret" >> $FNAME
echo "  createSecret: true" >> $FNAME
echo "  tenantID: \"$SP_TENANT_ID\"" >> $FNAME
echo "  clientID: \"$SP_CLIENT_ID\"" >> $FNAME
echo "  clientSecret: \"$SP_CLIENT_SECRET\"" >> $FNAME
echo >> $FNAME

# Set subscription ID if set by user
if [[ -v SUBSCRIPTION_ID ]]; then
    echo "defaultSubscriptionId: \"$SUBSCRIPTION_ID\"" >> $FNAME
    echo >> $FNAME
fi

# Set image address w/ default to 'adapter' as image name
# If the address is changed, use pullPolicy IfNotPresent to allow use of local images? TODO
if [[ -v REGISTRY && -v REGISTRY_PATH ]] || [[ -v VERSION ]]; then 
    echo "image:" >> $FNAME
fi

if [[ -v REGISTRY && -v REGISTRY_PATH ]]; then
    IMAGE="${IMAGE:-adapter}"

    if [[ "$REGISTRY_PATH" = "" ]]; then
        FULL_IMAGE=$REGISTRY/$IMAGE
    else
        FULL_IMAGE=$REGISTRY/$REGISTRY_PATH/$IMAGE
    fi

    echo "  repository: $FULL_IMAGE" >> $FNAME
    echo "  pullPolicy: IfNotPresent" >> $FNAME
fi 

if [[ -v VERSION ]]; then
    echo "  tag: $VERSION" >> $FNAME
fi
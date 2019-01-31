#!/bin/bash

set -o nounset
set -o errexit

GOPATH="${GOPATH:-$HOME/go}"
FNAME="$GOPATH/src/github.com/Azure/azure-k8s-metrics-adapter/local-dev-values.yaml"

if [[ -f $FNAME ]]; then
    echo "local-dev-values.yaml already exists and will not be altered"
    exit 1
fi

echo; echo "Creating local values file..."

# Currently only gens values for SP auth
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

# Set image address and/or version if either are set
if [[ -v IMAGE_REPOSITORY ]] || [[ -v VERSION ]]; then 
    echo "image:" >> $FNAME
    
    if [[ -v IMAGE_REPOSITORY ]]; then
        echo "  repository: $IMAGE_REPOSITORY" >> $FNAME
    fi 
    
    if [[ -v VERSION ]]; then
        echo "  tag: $VERSION" >> $FNAME
    fi
fi

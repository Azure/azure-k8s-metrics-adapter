#!/bin/bash

set -o nounset
set -o errexit

GOPATH="${GOPATH:-$HOME/go}"
FNAME="$GOPATH/src/github.com/Azure/azure-k8s-metrics-adapter/local-dev-values.yaml"

if [[ -f $FNAME ]]; then
    echo "local-dev-values.yaml already exists and will not be altered"
    exit
fi

echo; echo "Creating local values file..."

# Currently only gens values for SP auth
echo "azureAuthentication:" > $FNAME
echo "  method: clientSecret" >> $FNAME
echo "  createSecret: true" >> $FNAME
echo >> $FNAME

# Set subscription ID if set by user
if [[ -v SUBSCRIPTION_ID ]]; then
    echo "defaultSubscriptionId: \"$SUBSCRIPTION_ID\"" >> $FNAME
    echo >> $FNAME
fi

# Set image address and/or version if either are set
if [[ -v IMAGE ]] || [[ -v VERSION ]]; then 
    echo "image:" >> $FNAME
    
    if [[ -v IMAGE ]]; then
        echo "  repository: $IMAGE" >> $FNAME
    fi 
    
    if [[ -v VERSION ]]; then
        echo "  tag: $VERSION" >> $FNAME
    fi
fi

if [[ -v DOCKER_USER ]] && [[ -v DOCKER_PASS ]]; then
    echo "imageCredentials:" >> $FNAME
    echo "  createSecret: true" >> $FNAME

    if [[ ! "$REGISTRY" = "" ]]; then
        echo "  registry: $REGISTRY" >> $FNAME
    fi

    echo "  username: $DOCKER_USER" >> $FNAME
    echo "  password: $DOCKER_PASS" >> $FNAME
fi
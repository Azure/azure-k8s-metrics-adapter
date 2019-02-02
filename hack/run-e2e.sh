#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

GOPATH="${GOPATH:-$HOME/go}"

DIVIDER="============================================================"

echo "Checking cluster"
kubectl config current-context
JSONPATH='{range .items[*]}{@.metadata.name}{"\t"}Ready={@.status.conditions[?(@.type=="Ready")].status}{"\n"}{end}'
if kubectl get nodes -o jsonpath="$JSONPATH" | grep "Ready=False"; then 
    exit 1
fi

echo; echo "Checking helm and tiller install"
helm version

echo; echo "Checking for other pre-reqs"
which jq        # Used to format some command output
which docker

echo; echo "Checking that the Metrics Server is deployed"
kubectl get --raw "/apis/metrics.k8s.io/v1beta1/nodes" | jq .

# If the script fails on error here, the deployment won't be cleaned up properly
set +o errexit 

echo; echo "Running deployment scripts"
cd $GOPATH/src/github.com/Azure/azure-k8s-metrics-adapter/.azdevops
chmod +x 2_deploy/*.sh
chmod +x 3_output/*.sh

./2_deploy/deploy-adapter-with-sp.sh
if [[ $? = 0 ]]; then
    ./2_deploy/configure-metrics.sh

    echo "Testing deployment"

    ./3_output/gen-and-check-messages.sh
    if [[ $? = 0 ]];
        then echo $DIVIDER; echo "PASS"; echo $DIVIDER
        else echo $DIVIDER; echo "FAIL"; echo $DIVIDER; 
    fi

    ./3_output/run-consumer.sh
fi

echo "Removing adapter deployment"
helm delete --purge adapter

# TODO add final reporting message for clarity

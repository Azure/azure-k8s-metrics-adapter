#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

GOPATH="${GOPATH:-$HOME/go}"

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
cd $GOPATH/src/github.com/Azure/azure-k8s-metrics-adapter/.azdevops/2_deploy
chmod +x *.sh

./deploy-adapter-with-sp.sh
if [[ $? = 0 ]]; then
    ./configure-metrics.sh

    echo "Testing deployment"
    cd $GOPATH/src/github.com/Azure/azure-k8s-metrics-adapter/.azdevops/3_output
    chmod +x *.sh

    ./gen-and-check-messages.sh
    ./run-consumer.sh
fi

echo; echo "Removing adapter deployment"
helm delete --purge adapter

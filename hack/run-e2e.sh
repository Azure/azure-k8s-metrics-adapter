#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

GOPATH="${GOPATH:-$HOME/go}"

DIVIDER="============================================================"

echo "Checking for a current cluster context"
kubectl config current-context

echo; echo "Checking that cluster nodes are ready"
JSONPATH='{range .items[*]}{@.metadata.name}{"\t"}Ready={@.status.conditions[?(@.type=="Ready")].status}{"\n"}{end}'
kubectl get nodes -o jsonpath="$JSONPATH"
if kubectl get nodes -o jsonpath="$JSONPATH" | grep "Ready=False"; then
    exit 1
fi

echo; echo "Checking helm and tiller install"
helm version

echo; echo "Checking for other pre-reqs"
which jq        # Used to format some command output
which docker

# If the script fails on error here, the deployment won't be cleaned up properly
set +o errexit 

echo; echo "Running deployment scripts"
cd $GOPATH/src/github.com/Azure/azure-k8s-metrics-adapter/hack/e2e-scripts
chmod +x *.sh

TEST_FAILED=0

./deploy-adapter-with-sp.sh
if [[ $? = 0 ]]; then
    echo "Testing Queue (Azure Monitor) metrics"
    ./configure-queue-metrics.sh
    ./gen-and-check-queue-messages.sh
    if [[ $? = 0 ]];
        then echo $DIVIDER; echo "PASS"; echo $DIVIDER
        else echo $DIVIDER; echo "FAIL"; echo $DIVIDER; TEST_FAILED=1;
    fi

    ./run-queue-consumer.sh

    echo "Testing Topic Subscriptions metrics"
    ./configure-topic-subscriptions-metrics.sh
    ./gen-and-check-topic-subscriptions-messages.sh
    if [[ $? = 0 ]];
        then echo $DIVIDER; echo "PASS"; echo $DIVIDER
        else echo $DIVIDER; echo "FAIL"; echo $DIVIDER; TEST_FAILED=1;
    fi

    ./run-topic-consumer.sh
fi

echo "Removing adapter deployment"
helm delete --purge adapter

if [[ $TEST_FAILED == 1 ]]; then
    echo $DIVIDER; echo "FAIL"; echo $DIVIDER;
    exit 1
fi

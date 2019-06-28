#!/bin/bash

set -o nounset

GOPATH="${GOPATH:-$HOME/go}"
SERVICEBUS_QUEUE_NAME="${SERVICEBUS_QUEUE_NAME:-externalq}"

echo; echo "Running queue consumer to clear queue"
cd $GOPATH/src/github.com/Azure/azure-k8s-metrics-adapter/samples/servicebus-queue/

timeout 30 ./bin/consumer $SERVICEBUS_QUEUE_NAME> /dev/null

# Exit status 124 just means timeout completed, which is what we expect
if [[ $? = 124 ]]; then 
    echo "Consumer timed out as expected"
    exit 0
fi


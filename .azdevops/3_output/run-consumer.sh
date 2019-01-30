#!/bin/bash

set -o nounset
set -o errexit

GOPATH="${GOPATH:-$HOME/go}"
SERVICEBUS_QUEUE_NAME="${SERVICEBUS_QUEUE_NAME:-externalq}"

echo; echo "Running queue consumer..."
cd $GOPATH/src/github.com/Azure/azure-k8s-metrics-adapter/samples/servicebus-queue/

timeout 30 ./bin/consumer

# Exit status 124 just means timeout completed, which is what we expect
if [[ $? = 124 ]]; then 
    echo "Consumer timed out as expected"
fi
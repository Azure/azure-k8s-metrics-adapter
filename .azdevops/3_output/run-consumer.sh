#!/bin/bash

set -o nounset
set -o errexit

GOPATH="${GOPATH:-$HOME/go}"

echo; echo "Running queue consumer..."
cd $GOPATH/src/github.com/Azure/azure-k8s-metrics-adapter/samples/servicebus-queue/

timeout 30 ./bin/consumer

# Exit status 124 just means timeout completed, which is what we expect
if [[ $? == 124 ]]; then 
    exit 0
fi
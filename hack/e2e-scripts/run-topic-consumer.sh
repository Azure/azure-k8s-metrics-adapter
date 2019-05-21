#!/bin/bash

set -o nounset

GOPATH="${GOPATH:-$HOME/go}"

SERVICEBUS_TOPIC_NAME="${SERVICEBUS_TOPIC_NAME:-example-topic}"
SERVICEBUS_SUBSCRIPTION_NAME="${SERVICEBUS_SUBSCRIPTION_NAME:-externalsub}"

echo; echo "Running consumer to clear topic"
cd $GOPATH/src/github.com/Azure/azure-k8s-metrics-adapter/samples/servicebus-topic/

timeout 30 ./bin/consumer $SERVICEBUS_TOPIC_NAME $SERVICEBUS_SUBSCRIPTION_NAME> /dev/null

# Exit status 124 just means timeout completed, which is what we expect
if [[ $? = 124 ]]; then 
    echo "Consumer timed out as expected"
    exit 0
fi


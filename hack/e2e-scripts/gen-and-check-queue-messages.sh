#!/bin/bash

set -o nounset
set -o errexit

GOPATH="${GOPATH:-$HOME/go}"
SERVICEBUS_QUEUE_NAME="${SERVICEBUS_QUEUE_NAME:-externalq}"

cd $GOPATH/src/github.com/Azure/azure-k8s-metrics-adapter/samples/servicebus-queue/

echo; echo "Building producer and consumer..."
make

echo; echo "Creating random number for producer..."
NUM=$(( ($RANDOM % 30 )  + 1 ))

echo; echo "Sending $NUM messages..."
./bin/producer 0 $NUM $SERVICEBUS_QUEUE_NAME> /dev/null

echo; echo "Checking metrics endpoint for 4 minutes..."

MSGCOUNT=$(kubectl get --raw "/apis/external.metrics.k8s.io/v1beta1/namespaces/default/queuemessages" | jq .items[0].value)
START=`date +%s`

while [[ ! "$MSGCOUNT" = "\"$NUM\"" ]] && [[ $(( $(date +%s) - 225 )) -lt $START ]]; do
  sleep 15
  MSGCOUNT=$(kubectl get --raw "/apis/external.metrics.k8s.io/v1beta1/namespaces/default/queuemessages" | jq .items[0].value)
  echo "Endpoint returned $MSGCOUNT messages"
done

if [[ ! "$MSGCOUNT" = "\"$NUM\"" ]]; then
    echo "Timed out, message count ($MSGCOUNT) not equal to number of messages sent ($NUM)"
    exit 1
else
    echo "Message count equal to number of messages sent, metrics adapter working correctly"
fi

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

MSGCOUNT=$(kubectl get --raw "/apis/external.metrics.k8s.io/v1beta1/namespaces/default/queuemessages-total" | jq .items[0].value)
START=`date +%s`

while [[ ! "$MSGCOUNT" = "\"$NUM\"" ]] && [[ $(( $(date +%s) - 225 )) -lt $START ]]; do
  sleep 15
  MSGCOUNT=$(kubectl get --raw "/apis/external.metrics.k8s.io/v1beta1/namespaces/default/queuemessages-total" | jq .items[0].value)
  echo "Endpoint returned $MSGCOUNT messages"
done

AGGREGATE_TYPE=( "average" "maximum" "minimum" "total" )
for AGGREGATE in "${AGGREGATE_TYPE[@]}"
do
    METRIC_NAME="queuemessages-${AGGREGATE}"
    VALUE=$(kubectl get --raw "/apis/external.metrics.k8s.io/v1beta1/namespaces/default/${METRIC_NAME}" | jq .items[0].value)
    if [[ ! "$VALUE" = "\"$NUM\"" ]]; then
        echo "Timed out, message aggregate type: ${AGGREGATE} value: ${VALUE} not equal to number of messages sent ($NUM)"
        exit 1
    else
        echo "message aggregate type: ${AGGREGATE} value: ${VALUE} is equal to number of messages sent ($NUM), metrics adapter working correctly"
    fi
done

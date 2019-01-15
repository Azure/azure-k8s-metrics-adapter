#!/bin/bash

set -o nounset
set -o errexit

GOPATH="${GOPATH:-$HOME/go}"

cd $GOPATH/src/github.com/Azure/azure-k8s-metrics-adapter/samples/servicebus-queue/

echo; echo "Creating random number for producer script..."
NUM=$(( ($RANDOM % 30 )  + 1 ))
sed -i 's|20000|'$(( NUM + 1 ))'|g' producer/main.go

echo; echo "Building producer..."
make

echo; echo "Sending $NUM messages..."
./bin/producer 0 > /dev/null

echo; echo "Checking metrics endpoint..."

START=`date +%s`
MSGCOUNT=$(kubectl get --raw "/apis/external.metrics.k8s.io/v1beta1/namespaces/default/queuemessages" | jq .items[0].value)

while [[ "$MSGCOUNT" == "\"0\"" && $(( $(date +%s) - 155 )) -lt $START ]]; do
  sleep 15
  MSGCOUNT=$(kubectl get --raw "/apis/external.metrics.k8s.io/v1beta1/namespaces/default/queuemessages" | jq .items[0].value)
  echo "Endpoint returned $MSGCOUNT messages"
done

if [[ ! "$MSGCOUNT" == "\"$NUM\"" ]]; then
    echo "Message count ($MSGCOUNT) not equal to number of messages sent ($NUM)" 1>&2
else
    echo "Message count equal to number of messages sent, metrics adapter working correctly"
fi

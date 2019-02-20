#!/bin/bash

set -o nounset
set -o errexit

GOPATH="${GOPATH:-$HOME/go}"
SERVICEBUS_QUEUE_NAME="${SERVICEBUS_QUEUE_NAME:-externalq}"

cd $GOPATH/src/github.com/Azure/azure-k8s-metrics-adapter/samples/servicebus-queue/

# Copy original files
cp producer/main.go producer/main.go.copy
cp consumer/main.go consumer/main.go.copy

echo; echo "Creating random number for producer script..."
NUM=$(( ($RANDOM % 30 )  + 1 ))
sed -i 's|20000|'$(( NUM + 1 ))'|g' producer/main.go

echo; echo "Replacing queue name in consumer and producer..."
sed -i 's|externalq|'${SERVICEBUS_QUEUE_NAME}'|g' consumer/main.go
sed -i 's|externalq|'${SERVICEBUS_QUEUE_NAME}'|g' producer/main.go

echo; echo "Building producer and consumer..."
make

# Replace edited files with copies of original so git doesn't complain
rm producer/main.go; mv producer/main.go.copy producer/main.go
rm consumer/main.go; mv consumer/main.go.copy consumer/main.go

echo; echo "Sending $NUM messages..."
./bin/producer 0 > /dev/null

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

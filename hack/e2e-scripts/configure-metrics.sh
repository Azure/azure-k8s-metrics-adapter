#!/bin/bash

set -o nounset
set -o errexit

GOPATH="${GOPATH:-$HOME/go}"
SERVICEBUS_QUEUE_NAME="${SERVICEBUS_QUEUE_NAME:-externalq}"


echo; echo "Configuring external metric (queuemessages)..."
cd $GOPATH/src/github.com/Azure/azure-k8s-metrics-adapter/samples/servicebus-queue/
cp deploy/externalmetric.yaml deploy/externalmetric.yaml.copy

sed -i 's|sb-external-ns|'${SERVICEBUS_NAMESPACE}'|g' deploy/externalmetric.yaml
sed -i 's|sb-external-example|'${SERVICEBUS_RESOURCE_GROUP}'|g' deploy/externalmetric.yaml
sed -i 's|externalq|'${SERVICEBUS_QUEUE_NAME}'|g' deploy/externalmetric.yaml
kubectl apply -f deploy/externalmetric.yaml

rm deploy/externalmetric.yaml; mv deploy/externalmetric.yaml.copy deploy/externalmetric.yaml

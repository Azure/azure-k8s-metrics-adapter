#!/bin/bash

set -o nounset
set -o errexit

echo; echo "Configuring external metric (queuemessages)..."
cd $HOME/go/src/github.com/Azure/azure-k8s-metrics-adapter/samples/servicebus-queue/
sed -i 's|sb-external-ns|'${SERVICEBUS_NAMESPACE}'|g' deploy/externalmetric.yaml
sed -i 's|sb-external-example|'${SERVICEBUS_RESOURCE_GROUP}'|g' deploy/externalmetric.yaml
kubectl apply -f deploy/externalmetric.yaml

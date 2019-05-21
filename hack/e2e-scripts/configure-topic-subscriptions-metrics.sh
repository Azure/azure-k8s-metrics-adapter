#!/bin/bash

set -o nounset
set -o errexit

GOPATH="${GOPATH:-$HOME/go}"
SERVICEBUS_TOPIC_NAME="${SERVICEBUS_TOPIC_NAME:-example-topic}"
SERVICEBUS_SUBSCRIPTION_NAME="${SERVICEBUS_SUBSCRIPTION_NAME:-externalsub}"

echo; echo "Configuring external metric (subscriptionmessages)..."
cd $GOPATH/src/github.com/Azure/azure-k8s-metrics-adapter/samples/resources/externalmetric-examples/
cp servicebussubscription-example.yaml servicebussubscription-example.yaml.copy

sed -i 's|sb-external-ns|'${SERVICEBUS_NAMESPACE}'|g' servicebussubscription-example.yaml
sed -i 's|sb-external-example|'${SERVICEBUS_RESOURCE_GROUP}'|g' servicebussubscription-example.yaml
sed -i 's|example-topic|'${SERVICEBUS_TOPIC_NAME}'|g' servicebussubscription-example.yaml
sed -i 's|example-sub|'${SERVICEBUS_SUBSCRIPTION_NAME}'|g' servicebussubscription-example.yaml
kubectl apply -f servicebussubscription-example.yaml

rm servicebussubscription-example.yaml; mv servicebussubscription-example.yaml.copy servicebussubscription-example.yaml

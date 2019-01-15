#!/bin/bash

set -o nounset
set -o errexit

GOPATH="${GOPATH:-$HOME/go}"

echo; echo "Deploying queue consumer..."
cd $GOPATH/src/github.com/Azure/azure-k8s-metrics-adapter/samples/servicebus-queue/
kubectl create secret generic servicebuskey --from-literal=sb-connection-string=$SERVICEBUS_CONNECTION_STRING
kubectl apply -f deploy/consumer-deployment.yaml

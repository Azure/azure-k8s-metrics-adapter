#!/bin/bash

set -o nounset
set -o errexit

echo; echo "Deploying queue consumer..."
cd $HOME/go/src/github.com/Azure/azure-k8s-metrics-adapter/samples/servicebus-queue/
kubectl create secret generic servicebuskey --from-literal=sb-connection-string=$SERVICEBUS_CONNECTION_STRING
kubectl apply -f deploy/consumer-deployment.yaml

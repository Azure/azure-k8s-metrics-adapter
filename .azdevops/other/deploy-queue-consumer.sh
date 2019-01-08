#!/bin/bash

if [[ ! -v SERVICEBUS_CONNECTION_STRING ]]; then
    echo; echo "Must set SERVICEBUS_CONNECTION_STRING"
    exit 1
fi

echo; echo "Deploying queue consumer..."
cd $HOME/go/src/github.com/Azure/azure-k8s-metrics-adapter/samples/servicebus-queue/
kubectl create secret generic servicebuskey --from-literal=sb-connection-string=$SERVICEBUS_CONNECTION_STRING
kubectl apply -f deploy/consumer-deployment.yaml

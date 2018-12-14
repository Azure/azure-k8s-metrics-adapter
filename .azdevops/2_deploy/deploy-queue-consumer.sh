#!/bin/bash

if [[ ! -v GOPATH ]]; then
	echo "Must set GOPATH (/home/vsts/go on Azure Pipelines)"
	exit 1
fi

if [[ ! -v SERVICEBUS_CONNECTION_STRING ]]; then
    echo "Must set SERVICEBUS_CONNECTION_STRING"
    exit 1
fi

echo "Deploying queue consumer..."
cd $GOPATH/src/github.com/Azure/azure-k8s-metrics-adapter/samples/servicebus-queue/
kubectl create secret generic servicebuskey --from-literal=sb-connection-string=$(SERVICEBUS_CONNECTION_STRING)
kubectl apply -f deploy/consumer-deployment.yaml
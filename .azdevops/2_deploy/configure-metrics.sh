#!/bin/bash

if [[ ! -v GOPATH ]]; then
	echo "Must set GOPATH (/home/vsts/go on Azure Pipelines)"
	exit 1
fi

echo "Configuring external metric (queuemessages)..."
cd $GOPATH/src/github.com/Azure/azure-k8s-metrics-adapter/samples/servicebus-queue/
sed -i 's|sb-external-ns|'$(SERVICEBUS_NS)'|g' deploy/externalmetric.yaml
kubectl apply -f deploy/externalmetric.yaml
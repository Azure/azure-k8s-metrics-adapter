#!/bin/bash

if [[ ! -v GOPATH ]]; then
	echo "Must set GOPATH (/home/vsts/go on Azure Pipelines)"
	exit 1
fi

echo "Getting metrics adapter code..."
go get -d github.com/Azure/azure-k8s-metrics-adapter
echo "Getting metrics server code..."
go get -d github.com/kubernetes-incubator/metrics-server/...

echo "Getting go service bus library..."
cd $GOPATH/src/github.com/Azure/azure-k8s-metrics-adapter/samples/servicebus-queue
go get -u github.com/Azure/azure-service-bus-go

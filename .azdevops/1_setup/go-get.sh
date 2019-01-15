#!/bin/bash

GOPATH="${GOPATH:-$HOME/go}"

echo; echo "Getting metrics server code..."
go get -d github.com/kubernetes-incubator/metrics-server/...

echo; echo "Getting go service bus library..."
cd $GOPATH/src/github.com/Azure/azure-k8s-metrics-adapter/samples/servicebus-queue
go get -u github.com/Azure/azure-service-bus-go

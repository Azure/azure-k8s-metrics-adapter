#!/bin/bash

set -o nounset
set -o errexit

GOPATH="${GOPATH:-$HOME/go}"

echo; echo "Getting go service bus library..."
cd $GOPATH/src/github.com/Azure/azure-k8s-metrics-adapter/samples/servicebus-queue
go get -u github.com/Azure/azure-service-bus-go
 
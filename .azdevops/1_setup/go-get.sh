#!/bin/bash

go get -dv github.com/Azure/azure-k8s-metrics-adapter
go get -dv github.com/kubernetes-incubator/metrics-server/...

cd $GOPATH/src/github.com/Azure/azure-k8s-metrics-adapter/samples/servicebus-queue
go get -u github.com/Azure/azure-service-bus-go

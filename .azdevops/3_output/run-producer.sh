#!/bin/bash

echo; echo "Building producer..."
cd $GOPATH/src/github.com/Azure/azure-k8s-metrics-adapter/samples/servicebus-queue/
make

echo; echo "Running producer..."
./bin/producer 0

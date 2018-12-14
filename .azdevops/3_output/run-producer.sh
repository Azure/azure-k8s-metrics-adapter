#!/bin/bash

echo "Building producer..."
cd $GOPATH?src/github.com/Azure/azure-k8s-metrics-adapter/samples/servicebus-queue/
make

echo "Running producer..."
./bin/producer 0

#!/bin/bash

cd $GOPATH/src/github.com/Azure/azure-k8s-metrics-adapter/samples/servicebus-queue/

echo; echo "Creating random number for producer script..."
NUM=$(( ($RANDOM % 30 )  + 1 ))
sed -i 's|20000|'${NUM}'|g' producer/main.go

echo; echo "Building producer..."
make

echo; echo "Sending $NUM messages..."
./bin/producer 0 > /dev/null

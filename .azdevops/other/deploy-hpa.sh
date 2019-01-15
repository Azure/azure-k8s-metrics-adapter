#!/bin/bash

set -o nounset
set -o errexit

echo; echo "Deploying HPA..."
cd $HOME/go/src/github.com/Azure/azure-k8s-metrics-adapter/samples/servicebus-queue/
kubectl apply -f deploy/hpa.yaml

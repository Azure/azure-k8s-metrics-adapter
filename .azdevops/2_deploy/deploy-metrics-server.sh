#!/bin/bash

echo "Deploying metrics server (for k8s v1.8+)..."
cd $GOPATH/src/github.com/kubernetes-incubator/metrics-server/
kubectl create -f deploy/1.8+/




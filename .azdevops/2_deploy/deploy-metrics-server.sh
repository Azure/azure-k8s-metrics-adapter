#!/bin/bash

if [[ ! -v GOPATH ]]; then
	echo; echo "Must set GOPATH (/home/vsts/go on Azure Pipelines)"
	exit 1
fi

echo; echo "Deploying metrics server (for k8s v1.8+)..."
cd $GOPATH/src/github.com/kubernetes-incubator/metrics-server/
kubectl create -f deploy/1.8+/

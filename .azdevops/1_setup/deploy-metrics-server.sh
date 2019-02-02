#!/bin/bash

set -o nounset
set -o errexit

GOPATH="${GOPATH:-$HOME/go}"

echo; echo "Deploying metrics server (for k8s v1.8+)..."
cd $GOPATH/src/github.com/kubernetes-incubator/metrics-server/
kubectl create -f deploy/1.8+/

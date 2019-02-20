#!/bin/bash

set -o nounset
set -o errexit

echo; echo "Installing minikube..."
curl -Lo minikube https://storage.googleapis.com/minikube/releases/v$MINIKUBE_VERSION/minikube-linux-amd64
curl -Lo minikube.sha256 https://storage.googleapis.com/minikube/releases/v$MINIKUBE_VERSION/minikube-linux-amd64.sha256

expected=$(cat minikube.sha256) 
echo "$expected minikube" | sha256sum -c

chmod +x minikube
sudo mv minikube /usr/local/bin

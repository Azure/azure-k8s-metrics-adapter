#!/bin/bash

set -o nounset
set -o errexit

echo; echo "Installing minikube..."
curl -Lo minikube https://storage.googleapis.com/minikube/releases/v$MINIKUBE_VERSION/minikube-linux-amd64
chmod +x minikube
sudo mv minikube /usr/local/bin

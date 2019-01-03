#!/bin/bash

if [[ ! -v MINIKUBE_VERSION ]]; then
	echo; echo "Must set MINIKUBE VERSION (i.e. 0.32.0)"
	exit 1
fi

echo; echo "Installing minikube..."
curl -Lo minikube https://storage.googleapis.com/minikube/releases/v$MINIKUBE_VERSION/minikube-linux-amd64
chmod +x minikube
sudo mv minikube /usr/local/bin

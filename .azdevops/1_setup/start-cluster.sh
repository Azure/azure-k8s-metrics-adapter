#!/bin/bash

set -o nounset
set -o errexit

echo; echo "Starting minikube..."
sudo minikube start --vm-driver=none --bootstrapper=kubeadm --kubernetes-version=v$KUBERNETES_VERSION

echo; echo "Fixing permissions..."
sudo chown -R $USER $HOME/.kube
sudo chgrp -R $USER $HOME/.kube

sudo chown -R $USER $HOME/.minikube
sudo chgrp -R $USER $HOME/.minikube

minikube update-context

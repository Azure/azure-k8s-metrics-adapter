#!/bin/bash
if [[ ! -v KUBERNETES_VERSION ]]; then
	echo; echo "Must set KUBERNETES_VERSION (i.e. 1.12.0)"
	exit 1
fi

echo; echo "Starting minikube..."
sudo minikube start --vm-driver=none --bootstrapper=kubeadm --kubernetes-version=v$KUBERNETES_VERSION

echo; echo "Fixing permissions..."
sudo chown -R $USER $HOME/.kube
sudo chgrp -R $USER $HOME/.kube

sudo chown -R $USER $HOME/.minikube
sudo chgrp -R $USER $HOME/.minikube

minikube update-context
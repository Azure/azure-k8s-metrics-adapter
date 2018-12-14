#!/bin/bash

if [[ ! -v HELM_VERSION ]]; then
	echo; echo "Must set HELM_VERSION (i.e. 2.12.0)"
	exit 1
fi

echo; echo "Installing helm..."
curl -Lo helm.tar.gz https://storage.googleapis.com/kubernetes-helm/helm-v$HELM_VERSION-linux-amd64.tar.gz
tar -xvf helm.tar.gz
sudo mv linux-amd65/helm /usr/local/bin/helm
rm -rf linux-amd64 helm.tar.gz

#!/bin/bash

set -o nounset
set -o errexit

echo; echo "Installing helm..."
curl -Lo helm.tar.gz https://storage.googleapis.com/kubernetes-helm/helm-v$HELM_VERSION-linux-amd64.tar.gz
tar -xvf helm.tar.gz
sudo mv linux-amd64/helm /usr/local/bin/helm
rm -rf linux-amd64 helm.tar.gz

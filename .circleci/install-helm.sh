#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

curl -O https://storage.googleapis.com/kubernetes-helm/helm-v2.10.0-linux-amd64.tar.gz
curl -O https://storage.googleapis.com/kubernetes-helm/helm-v2.10.0-linux-amd64.tar.gz.sha256
expected=$(cat helm-v2.10.0-linux-amd64.tar.gz.sha256) 
echo "$expected helm-v2.10.0-linux-amd64.tar.gz" | sha256sum -c
tar -zxvf helm-v2.10.0-linux-amd64.tar.gz
sudo mv linux-amd64/helm /usr/local/bin/helm
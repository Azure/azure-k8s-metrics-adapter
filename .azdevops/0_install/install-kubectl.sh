#!/bin/bash

set -o nounset
set -o errexit

echo; echo "Installing kubectl..."
curl -Lo kubectl https://storage.googleapis.com/kubernetes-release/release/v$KUBERNETES_VERSION/bin/linux/amd64/kubectl
chmod +x kubectl
sudo mv kubectl /usr/local/bin

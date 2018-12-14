#!/bin/bash

if [[ ! -v KUBERNETES_VERSION ]]; then
	echo "Must set KUBERNETES_VERSION (i.e. 1.12.0)"
	exit 1
fi

echo "Installing kubectl..."
curl -Lo kubectl https://storage.googleapis.com/kubernetes-release/release/v$KUBERNETES_VERSION/bin/linux/amd64/kubectl
chmod +x kubectl
sudo mv kubectl /usr/local/bin

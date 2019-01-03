#!/bin/bash

if [[ ! -v KUBERNETES_VERSION ]]; then
	echo; echo "Must set KUBERNETES_VERSION (i.e. 1.12.4)"
	exit 1
fi

echo; echo "Installing kubectl..."
curl -Lo kubectl https://storage.googleapis.com/kubernetes-release/release/v$KUBERNETES_VERSION/bin/linux/amd64/kubectl
chmod +x kubectl
sudo mv kubectl /usr/local/bin

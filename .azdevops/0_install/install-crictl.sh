#!/bin/bash

if [[ ! $KUBERNETES_VERSION == 1.11.* ]]; then
	echo; echo "Not installing crictl (not needed)"
	exit 1
fi

echo; echo "Installing crictl..."
curl -Lo crictl.tar.gz https://github.com/kubernetes-sigs/cri-tools/releases/download/v1.11.1/crictl-v1.11.1-linux-amd64.tar.gz
sudo tar -C /usr/local/bin -xzf crictl.tar.gz
rm crictl.tar.gz

#!/bin/bash

set -o nounset
set -o errexit

if [[ ! $KUBERNETES_VERSION = 1.11.* ]]; then
	echo; echo "Not installing crictl (not needed)"
	exit
fi

echo; echo "Installing crictl..."
curl -Lo crictl.tar.gz https://github.com/kubernetes-sigs/cri-tools/releases/download/v1.13.0/crictl-v1.13.0-linux-amd64.tar.gz

expected="9bdbea7a2b382494aff2ff014da328a042c5aba9096a7772e57fdf487e5a1d51"
echo "$expected crictl.tar.gz" | sha256sum -c

sudo tar -C /usr/local/bin -xzf crictl.tar.gz
rm crictl.tar.gz

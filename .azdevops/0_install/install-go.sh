#!/bin/bash

if [[ ! -v GO_VERSION ]]; then
	echo; echo "Must set GO_VERSION (i.e. 1.11.3)"
	exit 1
fi

echo; echo "Installing Go..."
curl -Lo go.tar.gz https://dl.google.com/go/go$GO_VERSION.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go.tar.gz
rm go.tar.gz
export PATH=$PATH:/usr/local/go/bin

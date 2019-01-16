#!/bin/bash

set -o nounset
set -o errexit

echo; echo "Installing Go..."
curl -Lo go.tar.gz https://dl.google.com/go/go$GO_VERSION.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go.tar.gz
rm go.tar.gz
export PATH=$PATH:/usr/local/go/bin

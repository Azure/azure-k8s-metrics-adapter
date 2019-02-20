#!/bin/bash

set -o nounset
set -o errexit

echo; echo "Installing Tiller on cluster..."
helm init --upgrade --wait

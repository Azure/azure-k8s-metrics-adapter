#!/bin/bash

set -o nounset
set -o errexit

echo; echo "Installing socat, jq, & ebtables..."
sudo apt-get install socat jq ebtables

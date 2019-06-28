#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

# make sure we have correct version of code generator
cd $GOPATH/src/k8s.io/code-generator/
git fetch --all
git checkout tags/kubernetes-1.12.9 -b kubernetes-1.12.9





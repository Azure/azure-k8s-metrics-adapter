#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

# make sure we have correct version of code generator
cd $GOPATH/src/k8s.io/code-generator/

# https://stackoverflow.com/a/6245587/697126
currentbranch=$(git branch | grep \* | cut -d ' ' -f2)

if [[ $currentbranch != "release-1.14" ]]
then
  git fetch --all
  git checkout -t origin/release-1.14
fi




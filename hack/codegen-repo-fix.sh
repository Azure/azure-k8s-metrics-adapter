#!/usr/bin/env bash

# this is a hack for fixing casing why PR (https://github.com/kubernetes/kubernetes/pull/68484) get merged upstream.
# currently the code-generator does not support package names 
# with capitals (https://github.com/kubernetes/code-generator/issues/20#issuecomment-389311739)
# I made a fix on my own fork so this swaps out that code until I get the PR
# merge upstream

set -o errexit
set -o nounset
set -o pipefail


cd $GOPATH/src/k8s.io/code-generator/

remotes=$(git remote)

if [[ $remotes != *bugfix* ]]
then
    git remote add bugfix https://github.com/jsturtevant/code-generator.git
fi

# https://stackoverflow.com/a/6245587/697126
currentbranch=$(git branch | grep \* | cut -d ' ' -f2)

if [[ $currentbranch != "fix-casing" ]]
then
    git fetch bugfix
    git checkout -b fix-casing bugfix/fix-casing
fi




#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

$GOPATH/src/k8s.io/code-generator/generate-groups.sh all \
    github.com/Azure/azure-k8s-metrics-adapter/pkg/client \
    github.com/Azure/azure-k8s-metrics-adapter/pkg/apis \
    externalmetric:v1alpha1

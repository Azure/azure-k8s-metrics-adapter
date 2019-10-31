#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
CODEGEN_PKG=${CODEGEN_PKG:-$(cd "${SCRIPT_ROOT}"; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../code-generator)}
chmod +x ${CODEGEN_PKG}/generate-groups.sh

"${CODEGEN_PKG}"/generate-groups.sh all \
    github.com/Azure/azure-k8s-metrics-adapter/pkg/client \
    github.com/Azure/azure-k8s-metrics-adapter/pkg/apis \
    metrics:v1alpha2 \
    --go-header-file "${SCRIPT_ROOT}"/hack/boilerplate.go.txt

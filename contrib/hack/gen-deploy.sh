#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

SCRIPT_ROOT=$(dirname "${BASH_SOURCE}")
REPO_ROOT="${SCRIPT_ROOT}/../.."

gen() {
	OUTPUT=$1
    VALUES_PATH=""
	TMP_OUTPUT=$(mktemp)
	mkdir -p "$(dirname ${OUTPUT})"
    if [[ ! -z "$2" ]]; then
        VALUES=$2
        VALUES_PATH="--values=${SCRIPT_ROOT}/deploy/values/${VALUES}.yaml"
    fi
    helm template \
        "${REPO_ROOT}/charts/azure-k8s-metrics-adapter/" \
        --namespace "custom-metrics" \
        --name "azure-k8s-metrics-adapter" \
        --set "fullnameOverride=azure-k8s-metrics-adapter" \
		--set "createNamespaceResource=true" > "${TMP_OUTPUT}" \
        ${VALUES_PATH}
    
	mv "${TMP_OUTPUT}" "${OUTPUT}"
}

gen "${REPO_ROOT}/deploy/adapter.yaml" "adapter"

#!/bin/bash

set -o nounset
set -o errexit

echo; echo "Showing HPA pod data for 10 min..."
kubectl get hpa consumer-scaler
START=`date +%s`
while [ $(( $(date +%s) - 285 )) -lt $START ]; do
    sleep 15
    kubectl get hpa consumer-scaler --no-headers
done

#!/bin/bash

echo "Checking metrics endpoint..."
kubectl get --raw "/apis/external.metrics.k8s.io/v1beta1" | jq .

echo
kubectl  get --raw "/apis/external.metrics.k8s.io/v1beta1/namespaces/default/queuemessages" | jq .
#!/bin/bash

echo; echo "Checking metrics endpoint..."
until kubectl get --raw "/apis/external.metrics.k8s.io/v1beta1" | jq . 2>&1 | grep -q "external.metrics.k8s.io/v1beta1"
    do sleep 1
    echo "waiting for endpoint to return"
done

kubectl get --raw "/apis/external.metrics.k8s.io/v1beta1" | jq .

kubectl  get --raw "/apis/external.metrics.k8s.io/v1beta1/namespaces/default/queuemessages" | jq .
sleep 1
kubectl  get --raw "/apis/external.metrics.k8s.io/v1beta1/namespaces/default/queuemessages" | jq .
sleep 1
kubectl  get --raw "/apis/external.metrics.k8s.io/v1beta1/namespaces/default/queuemessages" | jq .
sleep 1
kubectl  get --raw "/apis/external.metrics.k8s.io/v1beta1/namespaces/default/queuemessages" | jq .
sleep 1
kubectl  get --raw "/apis/external.metrics.k8s.io/v1beta1/namespaces/default/queuemessages" | jq .
sleep 1

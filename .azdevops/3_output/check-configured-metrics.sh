#!/bin/bash

kubectl get deploy

kubectl describe deploy adapter-azure-k8s-metrics-adapter

echo; echo "Listing configured external metrics..."
kubectl get aem

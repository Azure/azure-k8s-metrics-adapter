#!/bin/bash

echo; echo "> kubectl get deploy"
kubectl get deploy

echo; echo "> kubectl describe deploy adapter-azure-k8s-metrics-adapter"
kubectl describe deploy adapter-azure-k8s-metrics-adapter

echo; echo "> kubectl get aem"
kubectl get aem

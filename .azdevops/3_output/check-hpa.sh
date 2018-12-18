#!/bin/bash

echo; echo "Showing HPA pod data for 10 min..."
START=`date +%s`
while [ $(( $(date +%s) - 600 )) -lt $START ]; do
    kubectl get hpa consumer-scaler
    sleep 15
done
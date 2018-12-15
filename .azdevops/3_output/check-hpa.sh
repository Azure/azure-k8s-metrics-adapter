#!/bin/bash

echo; echo "Showing HPA pod data for 5 min..."
timeout 300s kubectl get hpa consumer-scaler -w
#!/bin/bash

echo; echo "Showing HPA pod data for 5 min..."
timeout 600s kubectl get hpa consumer-scaler -w

# 124 is the exit status if the command times out rather than completing when using 
# `timeout`. Since this one can't complete and we just want to let it run for some
# time, timing out is successful operation
if [ $? = 124 ]; then
    exit 0
fi
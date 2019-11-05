#!/bin/bash

set -o nounset
set -o errexit

GOPATH="${GOPATH:-$HOME/go}"
SERVICEBUS_QUEUE_NAME="${SERVICEBUS_QUEUE_NAME:-externalq}"


echo; echo "Configuring external metric (queuemessages)..."
cd $GOPATH/src/github.com/Azure/azure-k8s-metrics-adapter/samples/servicebus-queue/
cp deploy/externalmetric.yaml deploy/externalmetric.yaml.copy

sed -i 's|sb-external-ns|'${SERVICEBUS_NAMESPACE}'|g' deploy/externalmetric.yaml
sed -i 's|sb-external-example|'${SERVICEBUS_RESOURCE_GROUP}'|g' deploy/externalmetric.yaml
sed -i 's|externalq|'${SERVICEBUS_QUEUE_NAME}'|g' deploy/externalmetric.yaml

AGGREGATE_TYPE=( "Average" "Maximum" "Minimum" "Total" )
# supported aggregates https://github.com/Azure/azure-sdk-for-go/blob/0acfc1d1083d148a606d380143176e218d437728/services/preview/monitor/mgmt/2018-03-01/insights/models.go#L38
for aggregate in "${AGGREGATE_TYPE[@]}"
do
    filePath="deploy/${aggregate}.externalmetric.yaml"
    echo "Creating ${aggregate} external metric with file: $filePath"
    cp deploy/externalmetric.yaml $filePath
    sed -i 's|Total|'${aggregate}'|g' $filePath

    kubectl apply -f $filePath
done

rm deploy/*externalmetric.yaml; mv deploy/externalmetric.yaml.copy deploy/externalmetric.yaml

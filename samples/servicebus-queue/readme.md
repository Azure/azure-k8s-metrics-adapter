# Service Bus Queue External Metric Scaling

This is an example of how to scale using Service Bus Queue as an external metric.  

- [Service Bus Queue External Metric Scaling](#service-bus-queue-external-metric-scaling)
    - [Walkthrough](#walkthrough)
    - [Setup Service Bus](#setup-service-bus)
    - [Setup AKS Cluster](#setup-aks-cluster)
        - [Enable Managed Service Identity (MSI)](#enable-managed-service-identity-msi)
        - [Start the producer](#start-the-producer)
        - [Configure Secret for consumer pod](#configure-secret-for-consumer-pod)
        - [Deploy Consumer](#deploy-consumer)
    - [Set up Azure Metrics Adapter](#set-up-azure-metrics-adapter)
        - [Deploy the adapter](#deploy-the-adapter)
        - [Deploy the HPA](#deploy-the-hpa)
    - [Scale!](#scale)
    - [Clean up](#clean-up)

## Walkthrough

Prequisites:

- provisioned an [AKS Cluster](https://docs.microsoft.com/en-us/azure/aks/kubernetes-walkthrough
- your `kubeconfig` points to your cluster.  
- [Metric Server deployed](https://github.com/kubernetes-incubator/metrics-server#deployment) to your cluster ([aks does not come with it deployed](https://github.com/Azure/AKS/issues/318)). Validate by running `kubectl get --raw "/apis/metrics.k8s.io/v1beta1/nodes" | jq .`

Get this repository and cd to this folder (on your GOPATH):

```
go get -u github.com/jsturtevant/azure-k8-metrics-adapter
cd $GOPATH/src/github.com/jsturtevant/azure-k8-metrics-adapter/samples/servicebus-queue/
```

## Setup Service Bus
Create a service bus in azure:

``` 
az group create -n sb-external-example -l eastus
az servicebus namespace create -n sb-external-ns -g sb-external-example
az servicebus queue create -n externalq --namespace-name sb-external-ns -g sb-external-example
```

Create an auth rules for queue:

```
az servicebus queue authorization-rule create --resource-group sb-external-example --namespace-name sb-external-ns --queue-name externalq  --name demorule --rights Listen Manage Send

#save for connection string for later
export SERVICEBUS_CONNECTION_STRING="$(az servicebus queue authorization-rule keys list --resource-group sb-external-example --namespace-name sb-external-ns --name demorule  --queue-name externalq | jq -r .primaryConnectionString)"
```

> note: this gives full access to the queue for ease of use of demo.  You should create more fine grained control for each component of your app.  For example the consumer app should only have `Listen` rights and the producer app should only have `Send` rights.

## Setup AKS Cluster

### Enable Managed Service Identity (MSI)
Run the scripts provided to [enable MSI](https://github.com/jsturtevant/azure-k8-metrics-adapter#azure-setup) with the following environment variables for giving the MSI access to the Service Bus Namespace Insights provider.

```
export ACCESS_RG=sb-external-example
```

### Start the producer
Make sure you have cloned this repository and are in the folder `samples/servicebus-queue` for remainder of this walkthrough.

Build the project:

```bash
#build
go get -u github.com/Azure/azure-service-bus-go
make
```

Run the producer to create a few queue items, then hit `ctl-c` after a few message have been sent to stop it:

```
./bin/producer 500
```

Check the queue has values:

```
az servicebus queue show --resource-group sb-external-example --namespace-name sb-external-ns --name externalq | jq .messageCount
```

### Configure Secret for consumer pod
Create a secret with the connection string (from [previous step](#setup-service-bus)) for the service bus:

```
kubectl create secret generic servicebuskey --from-literal=sb-connection-string=$SERVICEBUS_CONNECTION_STRING
```

### Deploy Consumer 
Deploy the consumer:

```
kubectl apply -f deploy/consumer-deployment.yaml
```

Check that the consumer was able to recieve messages:

```
k logs -l app=consumer

# output should look something like
connecting to queue:  externalq
setting up listener
listening...
received message:  the answer is 42
number message left:  6
received message:  the answer is 42
number message left:  5
received message:  the answer is 42
number message left:  4
received message:  the answer is 42
number message left:  3
received message:  the answer is 42
number message left:  2
received message:  the answer is 42
```

## Set up Azure Metrics Adapter

### Deploy the adapter
Deploy the adapter:

```
kubectl apply -f https://raw.githubusercontent.com/jsturtevant/azure-k8-metrics-adapter/master/deploy/adapter.yaml
```

Check you can hit the external metric endpoint.  The resources will be empty as it [is not implemented yet](https://github.com/jsturtevant/azure-k8-metrics-adapter/issues/3) but you should receive a result.

```
kubectl get --raw "/apis/external.metrics.k8s.io/v1beta1" | jq .

# output should look something like
{
  "kind": "APIResourceList",
  "apiVersion": "v1",
  "groupVersion": "external.metrics.k8s.io/v1beta1",
  "resources": []
}
```

### Deploy the HPA
Deploy the HPA:

```
kubectl apply -f deploy/hpa.yaml
```

After a few seconds, validate that the HPA is configured.  If the `targets` shows `<unknown>` wait longer and try again.

```
kubectl get hpa consumer-scaler
```

You can also check the queue value returns manually by running:

```
kubectl  get --raw "/apis/external.metrics.k8s.io/v1beta1/namespaces/test/queuemessages?labelSelector=resourceProviderNamespace=Microsoft.Servicebus,resourceType=namespaces,aggregation=Total,filter=EntityName_eq_externalq,resourceGroup=sb-external-example,resourceName=sb-external-ns,metricName=Messages" | jq .
```

## Scale!

Put some load on the queue. Note this will add 20,000 message then exit.

```
./bin/producer 0
```

Now check your queue is loaded:

```
az servicebus queue show --resource-group sb-external-example --namespace-name sb-external-ns --name externalq | jq .messageCount

// should have a good 19,000 or more
19,858
```

Now watch your HPA pick up the queue count and scal the replicas.  This will take 2 or 3 minutes due to the fact that the HPA check happens every 30 seconds and must have a value over target value.  It also takes a few seconds for the metrics to be reported to the Azure endpoint that is queried.  It will then scale back down after a few minutes as well.

```
kubectl get hpa consumer-scaler -w

#output similiar to
NAME              REFERENCE             TARGETS   MINPODS   MAXPODS   REPLICAS   AGE  
consumer-scaler   Deployment/consumer   0/30      1         10        1          1h   
consumer-scaler   Deployment/consumer   27278/30   1         10        1         1h   
consumer-scaler   Deployment/consumer   26988/30   1         10        4         1h   
consumer-scaler   Deployment/consumer   26988/30   1         10        4         1h           consumer-scaler   Deployment/consumer   26702/30   1         10        4         1h           
consumer-scaler   Deployment/consumer   26702/30   1         10        4         1h           
consumer-scaler   Deployment/consumer   25808/30   1         10        4         1h           
consumer-scaler   Deployment/consumer   25808/30   1         10        4         1h           consumer-scaler   Deployment/consumer   24784/30   1         10        8         1h           consumer-scaler   Deployment/consumer   24784/30   1         10        8         1h          
consumer-scaler   Deployment/consumer   23775/30   1         10        8         1h           
consumer-scaler   Deployment/consumer   22065/30   1         10        8         1h           
consumer-scaler   Deployment/consumer   22065/30   1         10        8         1h           
consumer-scaler   Deployment/consumer   20059/30   1         10        8         1h           
consumer-scaler   Deployment/consumer   20059/30   1         10        10        1h
```

Once it is scaled up you can check the deployment:

```
kubectl get deployment consumer

NAME       DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
consumer   10         10         10            10           23m
```

And check out the logs for the consumers:

```
k logs -l app=consumer --tail 100
```

## Clean up
Once the queue is empty (will happen pretty quickly after scaled up to 10) you should see your deployment scale back down.

Once you are done with this experiment you can delete kubernetes deployments and  the resource group:

```
kubectl delete -f deploy/hpa.yaml
kubectl delete -f deploy/consumer-deployment.yaml
kubectl detele -f https://raw.githubusercontent.com/jsturtevant/azure-k8-metrics-adapter/master/deploy/adapter.yaml

az group delete -n sb-external-example
```
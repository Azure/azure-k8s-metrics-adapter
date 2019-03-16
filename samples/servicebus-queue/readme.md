# Service Bus Queue External Metric Scaling

This is an example of how to scale using Service Bus Queue as an external metric.  

- [Service Bus Queue External Metric Scaling](#service-bus-queue-external-metric-scaling)
  - [Walkthrough](#walkthrough)
  - [Setup Service Bus](#setup-service-bus)
  - [Setup AKS Cluster](#setup-aks-cluster)
    - [Enable Access to Azure Resources](#enable-access-to-azure-resources)
    - [Start the producer](#start-the-producer)
    - [Configure Secret for consumer pod](#configure-secret-for-consumer-pod)
    - [Deploy Consumer](#deploy-consumer)
  - [Set up Azure Metrics Adapter](#set-up-azure-metrics-adapter)
    - [Deploy the adapter](#deploy-the-adapter)
    - [Configure Metric Adapter with metrics](#configure-metric-adapter-with-metrics)
    - [Deploy the HPA](#deploy-the-hpa)
  - [Scale!](#scale)
  - [Clean up](#clean-up)

## Walkthrough

Prerequisites:

- provisioned an [AKS Cluster](https://docs.microsoft.com/en-us/azure/aks/kubernetes-walkthrough)
- your `kubeconfig` points to your cluster.  
- [Metric Server deployed](https://github.com/kubernetes-incubator/metrics-server#deployment) to your cluster. Validate by running `kubectl get --raw "/apis/metrics.k8s.io/v1beta1/nodes" | jq .`
- [helm](https://docs.helm.sh/using_helm/#quickstart-guide) or install [helm on aks](https://docs.microsoft.com/en-us/azure/aks/kubernetes-helm) (If you have RBAC enabled, you will need to configure permissions as outlined in the second link)

Get this repository and cd to this folder (on your GOPATH):

```
go get -u github.com/Azure/azure-k8s-metrics-adapter
cd $GOPATH/src/github.com/Azure/azure-k8s-metrics-adapter/samples/servicebus-queue/
```

## Setup Service Bus
Create a service bus in azure:

``` 
export SERVICEBUS_NS=sb-external-ns-<your-initials>
az group create -n sb-external-example -l eastus
az servicebus namespace create -n $SERVICEBUS_NS -g sb-external-example
az servicebus queue create -n externalq --namespace-name $SERVICEBUS_NS -g sb-external-example
```

Create an auth rules for queue:

```
az servicebus queue authorization-rule create --resource-group sb-external-example --namespace-name $SERVICEBUS_NS --queue-name externalq  --name demorule --rights Listen Manage Send

#save for connection string for later
export SERVICEBUS_CONNECTION_STRING="$(az servicebus queue authorization-rule keys list --resource-group sb-external-example --namespace-name $SERVICEBUS_NS --name demorule  --queue-name externalq -o json | jq -r .primaryConnectionString)"
```

> note: this gives full access to the queue for ease of use of demo.  You should create more fine grained control for each component of your app.  For example the consumer app should only have `Listen` rights and the producer app should only have `Send` rights.

## Setup AKS Cluster

### Enable Access to Azure Resources

Run the scripts provided to either [enable MSI](https://github.com/Azure/azure-k8s-metrics-adapter#azure-setup), [Azure AD Pod Identity](https://github.com/Azure/azure-k8s-metrics-adapter#azure-setup#using-azure-ad-pod-identity
) or [configure a Service Principal](https://github.com/Azure/azure-k8s-metrics-adapter/blob/master/README.md#using-azure-ad-application-id-and-secret) with the following environment variables for giving the access to the Service Bus Namespace Insights provider.

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
# 0 delay in sending messages, send 5 messages to queue 'externalq'
./bin/producer 0 5 externalq
```

Check the queue has values:

```
az servicebus queue show --resource-group sb-external-example --namespace-name $SERVICEBUS_NS --name externalq -o json | jq .messageCount
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

Check that the consumer was able to receive messages:

```
kubectl logs -l app=consumer

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
helm install --name sample-release ../../charts/azure-k8s-metrics-adapter --namespace custom-metrics 
```


>Note: if you used a Service Principal you will need the deployment with a service principal configured and a secret deployed with the service principal values 
>
>```
>helm install --name sample-release ../../charts/azure-k8s-metrics-adapter --namespace custom-metrics --set azureAuthentication.method=clientSecret --set azureAuthentication.tenantID=<your tenantid> --set azureAuthentication.clientID=<your clientID> --set azureAuthentication.clientSecret=<your clientSecret> --set azureAuthentication.createSecret=true`
>```

> Note: if you used [Azure AD Pod Identity](../../README.md#using-azure-ad-pod-identity) you need to use the specific adapter template file that declares the Azure Identity Binding on [Line 49](../../deploy/adapter-aad-pod-identity.yaml#L49) and [Line 61](../../deploy/adapter-aad-pod-identity.yaml#L61).
> ```bash
> kubectl apply -f >https://raw.githubusercontent.com/Azure/azure-k8s-metrics-adapter/master/deploy/adapter-aad-pod-identity.yaml
>```


Check you can hit the external metric endpoint.  The resources will be empty as it [is not implemented yet](https://github.com/Azure/azure-k8s-metrics-adapter/issues/3) but you should receive a result.

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

### Configure Metric Adapter with metrics
The metric adapter deploys a CRD called ExternalMetric which you can use to configure metrics.  To deploy these metric we need to update the Service Bus namespace in the configuration then deploy it:

```
sed -i 's|sb-external-ns|'${SERVICEBUS_NS}'|g' deploy/externalmetric.yaml
kubectl apply -f deploy/externalmetric.yaml
```

> note: the ExternalMetric configuration is deployed per namespace.

You can list of the configured external metrics via:

```
kubectl get aem #shortcut for externalmetric
```

### Deploy the HPA
Deploy the HPA:

```
kubectl apply -f deploy/hpa.yaml
```

> note: the `external.metricName` defined on the HPA must match the `metadata.name` on the ExternalMetric declaration, in this case `queuemessages`

After a few seconds, validate that the HPA is configured.  If the `targets` shows `<unknown>` wait longer and try again.

```
kubectl get hpa consumer-scaler
```

You can also check the queue value returns manually by running:

```
kubectl  get --raw "/apis/external.metrics.k8s.io/v1beta1/namespaces/default/queuemessages" | jq .
```

## Scale!

Put some load on the queue. Note this will add 5,000 message then exit.

```
# 0 delay in sending messages, send 5000 messages to queue 'externalq'
./bin/producer 0 5000 externalq
```

Now check your queue is loaded:

```
az servicebus queue show --resource-group sb-external-example --namespace-name $SERVICEBUS_NS --name externalq -o json | jq .messageCount

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
helm delete --purge sample-release

az group delete -n sb-external-example
```

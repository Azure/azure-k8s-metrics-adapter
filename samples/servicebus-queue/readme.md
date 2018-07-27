# Service Bus Queue External Metric Scaling

This is an example of how to scale using Service Bus Queue as an external metric.  

## Walkthrough


## Setup Service Bus
Create a service bus in azure:

```
az servicebus namespace create
az servicebus queue create
```

### Configure Secret for consumer pod
Create a secret with the connection string for the service bus:

```
kubectl create secret generic servicebuskey --from-literal=sb-connection-string='<connection-string>'
```

> the quotes around the connection string are needed

### Start the producer
Build the project:

```bash
#build
go get -u github.com/Azure/azure-service-bus-go
make
```

In a separte windows, Run the producer to create a few queue items:

```
export SERVICEBUS_CONNECTION_STRING='your-connstring' 
./bin/producer 500
```

Check the queue has values:

```
az servicebus queue show --resource-group myresourcegroup --namespace-name mynamespace --name myqueue
```

### Deploy Consumer 
Deploy the consumer:

```
kubectl deploy -f deploy/consumer-deployment.yaml
```

Check the logs for the consumer to see it connected to the queue:

```
#list pods
kubectl get pod

#use pod name from list
kubectl logs <podname>
Hello World!
Hello World!
```



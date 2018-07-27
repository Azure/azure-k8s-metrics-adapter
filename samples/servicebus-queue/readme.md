# Service Bus Queue External Metric Scaling

This is an example of how to scale using Service Bus Queue as an external metric.  

## Setup

Create a service bus:

```
az servicebus namespace create
az servicebus queue create
```

## To Run

```bash
#build
go get github.com/Azure/azure-service-bus-go
make

#run
SERVICEBUS_CONNECTION_STRING='your-connstring' ./bin/producer
```


# Service Bus Queue External Metric Scaling

This is an example of how to scale using Service Bus Queue as an external metric.  The producer and consumer are from the quick start repository https://github.com/Azure/azure-service-bus-go.

## Setup

Create a service bus:

```
az servicebus create
```

## To Run

`go get github.com/Azure/azure-service-bus-go`

- from this directory execute `make`
- open two terminal windows
  - in the first terminal, execute `SERVICEBUS_CONNECTION_STRING='your-connstring' ./bin/consumer`
  - in the second terminal, execute `SERVICEBUS_CONNECTION_STRING='your-connstring' ./bin/producer`
  - in the second terminal, type some words and press enter
- see the words you typed in the second terminal in the first
- type 'exit' in the second terminal when you'd like to end your session

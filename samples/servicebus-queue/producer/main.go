package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/Azure/azure-service-bus-go"
)

func main() {
	connStr := os.Getenv("SERVICEBUS_CONNECTION_STRING")
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connStr))
	if err != nil {
		fmt.Println("namespace: ", err)
	}

	// Initialize and create a Service Bus Queue named helloworld if it doesn't exist
	queueName := "helloworld"
	q, err := ns.NewQueue(queueName)
	if err != nil {
		// handle queue creation error
		fmt.Println("create queue: ", err)
	}

	// Send message to the Queue named helloworld
	err = q.Send(context.Background(), servicebus.NewMessageFromString("Hello World!"))
	if err != nil {
		// handle message send error
		fmt.Println("send message: ", err)
	}

	// Receive message from queue named helloworld
	listenHandle, err := q.Receive(context.Background(),
		func(ctx context.Context, msg *servicebus.Message) servicebus.DispositionAction {
			fmt.Println(string(msg.Data))
			return msg.Complete()
		})
	if err != nil {
		// handle queue listener creation err
		fmt.Println("listener: ", err)
	}
	// Close the listener handle for the Service Bus Queue
	defer listenHandle.Close(context.Background())

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan
}

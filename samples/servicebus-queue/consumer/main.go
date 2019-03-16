package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/Azure/azure-service-bus-go"
)

func main() {
	queueName := os.Args[1]

	connStr := os.Getenv("SERVICEBUS_CONNECTION_STRING")
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connStr))
	if err != nil {
		fmt.Println("namespace: ", err)
		panic(err)
	}

	fmt.Println("connecting to queue: ", queueName)
	q, err := ns.NewQueue(queueName)
	if err != nil {
		// handle queue creation error
		fmt.Println("create queue: ", err)
	}

	fmt.Println("setting up listener")
	var messageHandler servicebus.HandlerFunc = func(ctx context.Context, msg *servicebus.Message) error {
		fmt.Println("received message: ", string(msg.Data))
		return msg.Complete(ctx)
	}

	err = q.Receive(context.Background(), messageHandler)
	if err != nil {
		// handle queue listener creation err
		fmt.Println("listener error: ", err)
	}

	fmt.Println("listening...")

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan
}

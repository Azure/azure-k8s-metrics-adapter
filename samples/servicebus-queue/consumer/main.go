package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/Azure/azure-service-bus-go"
)

func main() {
	connStr := os.Getenv("SERVICEBUS_CONNECTION_STRING")
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connStr))
	if err != nil {
		fmt.Println("namespace: ", err)
		panic(err)
	}

	queueName := "externalq"
	qm := ns.NewQueueManager()

	fmt.Println("connecting to queue: ", queueName)
	q, err := ns.NewQueue(queueName)
	if err != nil {
		// handle queue creation error
		fmt.Println("create queue: ", err)
	}

	fmt.Println("setting up listener")
	var messageHandler servicebus.HandlerFunc = func(ctx context.Context, msg *servicebus.Message) servicebus.DispositionAction {
		fmt.Println("received message: ", string(msg.Data))

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		qe, err := qm.Get(ctx, queueName)
		if err != nil {
			fmt.Println("create manager created error: ", err)
		}
		fmt.Println("number message left: ", *qe.MessageCount)
		return msg.Complete()
	}

	err = q.Receive(context.Background(), messageHandler)
	if err != nil {
		// handle queue listener creation err
		fmt.Println("listener: ", err)
	}

	fmt.Println("listening...")

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan
}

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/Azure/azure-service-bus-go"
)

// Printer is a type that prints the message
type Printer struct{}

// Handle takes the message and prints contents
func (p Printer) Handle(ctx context.Context, msg *servicebus.Message) error {
	fmt.Println(string(msg.Data))
	return msg.Complete(ctx)
}

func main() {
	topicName := os.Args[1]
	subscriptionName := os.Args[2]

	connStr := os.Getenv("SERVICEBUS_CONNECTION_STRING")
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connStr))
	if err != nil {
		fmt.Println("namespace: ", err)
		panic(err)
	}

	topic, err := ns.NewTopic(topicName)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("setting up subscription reciever %s on topic %s", subscriptionName, topicName)
	sub, err := topic.NewSubscription(subscriptionName)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = sub.Receive(context.Background(), Printer{})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("listening...")

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan
}

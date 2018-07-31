package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	servicebus "github.com/Azure/azure-service-bus-go"
)

func main() {
	speedArg := os.Args[1]
	speed, err := strconv.Atoi(speedArg)
	if err != nil {
		fmt.Println("Please provide speed in milliseconds")
		return
	}

	connStr := os.Getenv("SERVICEBUS_CONNECTION_STRING")
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connStr))
	if err != nil {
		fmt.Println("namespace: ", err)
	}

	// Initialize and create a Service Bus Queue named helloworld if it doesn't exist
	queueName := "externalq"
	fmt.Println("connecting to queue: ", queueName)
	q, err := ns.NewQueue(queueName)
	if err != nil {
		// handle queue creation error
		fmt.Println("create queue: ", err)
	}

	//https: //stackoverflow.com/a/18158859/697126
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	go func() {
		<-signalChan
		os.Exit(1)
	}()

	i := 0
	for {
		fmt.Println("sending message ", i)
		err = q.Send(context.Background(), servicebus.NewMessageFromString("the answer is 42"))
		if err != nil {
			// handle message send error
			fmt.Println("error sending message: ", err)
		}

		if i > 1000 {
			break
		}

		time.Sleep(time.Duration(speed) * time.Millisecond)
	}

}

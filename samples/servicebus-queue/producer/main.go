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

	messagesToSendArg := os.Args[2]
	messagesCount, err := strconv.Atoi(messagesToSendArg)
	if err != nil {
		fmt.Println("Please provide number of messages")
		return
	}

	queueName := os.Args[3]

	connStr := os.Getenv("SERVICEBUS_CONNECTION_STRING")
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connStr))
	if err != nil {
		fmt.Println("namespace: ", err)
	}

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

	for i := 1; i <= messagesCount; i++ {
		fmt.Println("sending message ", i)
		err = q.Send(context.Background(), servicebus.NewMessageFromString("the answer is 42"))
		if err != nil {
			// handle message send error
			fmt.Println("error sending message: ", err)
		}

		time.Sleep(time.Duration(speed) * time.Millisecond)
	}

}

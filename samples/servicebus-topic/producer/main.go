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

	topicName := os.Args[3]
	subscription := os.Args[4]

	connStr := os.Getenv("SERVICEBUS_CONNECTION_STRING")
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connStr))
	if err != nil {
		fmt.Println("namespace: ", err)
	}

	topic, err := ns.NewTopic(topicName)
	if err != nil {
		fmt.Println(err)
		return
	}

	sm, err := ns.NewSubscriptionManager(topicName)
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	_, err = ensureSubscription(ctx, sm, subscription)
	if err != nil {
		fmt.Println(err)
		return
	}

	// remove the default rule, which is the "TrueFilter" that accepts all messages
	err = sm.DeleteRule(ctx, subscription, "$Default")
	if err != nil {
		fmt.Printf("delete default rule err: %s", err)
		return
	}

	exp := fmt.Sprintf("subscription = '%s'", subscription)
	fmt.Printf("filter is %s\n", exp)
	_, err = sm.PutRule(ctx, subscription, subscription+"Rule", servicebus.SQLFilter{Expression: exp})
	if err != nil {
		fmt.Printf("add rule err: %s", err)
		return
	}

	//https: //stackoverflow.com/a/18158859/697126
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	go func() {
		<-signalChan
		os.Exit(1)
	}()

	fmt.Printf("sending %d messages ", messagesCount)
	for i := 1; i <= messagesCount; i++ {
		fmt.Printf("sending message %d to sub %s\n", i, subscription)
		m := servicebus.NewMessageFromString("the answer is 42")
		m.UserProperties = map[string]interface{}{"subscription": subscription}
		err = topic.Send(context.Background(), m)
		if err != nil {
			// handle message send error
			fmt.Println("error sending message: ", err)
		}

		time.Sleep(time.Duration(speed) * time.Millisecond)
	}

}

func ensureSubscription(ctx context.Context, sm *servicebus.SubscriptionManager, name string, opts ...servicebus.SubscriptionManagementOption) (*servicebus.SubscriptionEntity, error) {
	subEntity, err := sm.Get(ctx, name)
	if err == nil {
		_ = sm.Delete(ctx, name)
	}

	subEntity, err = sm.Put(ctx, name, opts...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return subEntity, nil
}

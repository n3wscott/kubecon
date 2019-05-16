package main

import (
	"context"
	"fmt"
	"github.com/cloudevents/sdk-go"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/transport/nats"
	"github.com/n3wscott/kubecon/helloworld/pkg/person"
)

func main() {
	t, err := nats.New("localhost:4222", "hello")
	if err != nil {
		panic(err)
	}
	client, _ := cloudevents.NewClient(t)

	if err := client.StartReceiver(context.Background(), Receive); err != nil {
		panic(err)
	}
}

func Receive(event cloudevents.Event) {
	data := &person.Hello{}
	if err := event.DataAs(data); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Hello, %s!\n", data.Name)

	fmt.Printf("\n---☁️  Event---\n%s\n\n", event)
}

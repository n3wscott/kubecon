package main

import (
	"context"
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

	if err := client.StartReceiver(context.Background(), person.Receive); err != nil {
		panic(err)
	}
}

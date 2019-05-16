package main

import (
	"context"
	"github.com/cloudevents/sdk-go"
	"github.com/n3wscott/kubecon/helloworld/pkg/person"
)

func main() {
	client, _ := cloudevents.NewDefaultClient()

	if err := client.StartReceiver(context.Background(), person.Receive); err != nil {
		panic(err)
	}
}

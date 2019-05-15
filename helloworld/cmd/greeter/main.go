package main

import (
	"context"
	"github.com/cloudevents/sdk-go"
	"github.com/n3wscott/kubecon/helloworld/pkg/hello"
)

func main() {
	client, err := cloudevents.NewDefaultClient()
	if err != nil {
		panic(err)
	}

	if err := client.StartReceiver(context.Background(), hello.Receive); err != nil {
		panic(err)
	}
}

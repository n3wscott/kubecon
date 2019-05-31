package main

import (
	"context"
	"fmt"
	"github.com/cloudevents/sdk-go"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/transport"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/transport/nats"
	"github.com/n3wscott/kubecon/helloworld/pkg/person"
)

func main() {
	client, _ := cloudevents.NewClient(
		makeOutboundHTTP(),
		cloudevents.WithUUIDs(),
	)

	event := cloudevents.NewEvent()
	event.SetSource("/kubecon/demo/barcelona-2019")
	event.SetType("com.n3wscott.hello")
	_ = event.SetData(&person.Hello{
		Name: "Hanna",
	})

	if _, err := client.Send(context.Background(), event); err != nil {
		panic(err)
	}
	fmt.Println("sent")
}

func makeOutboundNATS() transport.Transport {
	t, err := nats.New("localhost:4222", "hello")
	if err != nil {
		panic(err)
	}
	return t
}

func makeOutboundHTTP() transport.Transport {
	t, err := cloudevents.NewHTTPTransport(
		cloudevents.WithTarget("http://localhost:8181/"),
		cloudevents.WithBinaryEncoding(),
	)
	if err != nil {
		panic(err)
	}
	return t
}

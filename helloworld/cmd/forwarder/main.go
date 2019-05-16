package main

import (
	"context"
	"fmt"
	"github.com/cloudevents/sdk-go"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/transport"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/transport/nats"
)

func main() {
	inboundNATS, _ := cloudevents.NewClient(
		makeInboundNATS(),
		cloudevents.WithUUIDs(),
	)

	inboundHTTP, _ := cloudevents.NewClient(
		makeInboundHTTP(),
		cloudevents.WithUUIDs(),
	)

	outbound, _ := cloudevents.NewClient(
		makeOutboundHTTP(),
		cloudevents.WithUUIDs(),
	)

	forward := func(event cloudevents.Event) {
		fmt.Printf("forwarding %s\n", event.ID())
		_, _ = outbound.Send(context.Background(), event)
	}

	go func() {
		if err := inboundHTTP.StartReceiver(context.Background(), forward); err != nil {
			panic(err)
		}
	}()

	go func() {
		if err := inboundNATS.StartReceiver(context.Background(), forward); err != nil {
			panic(err)
		}
	}()
	// block
	<-context.Background().Done()
}

func makeInboundHTTP() transport.Transport {
	t, err := cloudevents.NewHTTPTransport(
		cloudevents.WithPort(8181),
		cloudevents.WithBinaryEncoding(),
	)
	if err != nil {
		panic(err)
	}
	return t
}

func makeInboundNATS() transport.Transport {
	t, err := nats.New("localhost:4222", "hello")
	if err != nil {
		panic(err)
	}
	return t
}

func makeOutboundHTTP() transport.Transport {
	t, err := cloudevents.NewHTTPTransport(
		cloudevents.WithTarget("http://localhost:8080/"),
		cloudevents.WithBinaryEncoding(),
	)
	if err != nil {
		panic(err)
	}
	return t
}

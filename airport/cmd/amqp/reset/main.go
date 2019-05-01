package main

import (
	"context"
	"log"
	"os"

	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/transport/amqp"
	"github.com/kelseyhightower/envconfig"
	"github.com/n3wscott/kubecon/airport/pkg/events"
	qp "pack.ag/amqp"
)

const (
	count = 10
)

type envConfig struct {
	// AMQPServer URL to connect to the amqp server.
	AMQPServer string `envconfig:"AMQP_SERVER" default:"amqp://localhost:5672/" required:"true"`

	// Queue is the amqp queue name to publish cloudevents on.
	Queue string `envconfig:"AMQP_QUEUE"`

	AccessKeyName string `envconfig:"AMQP_ACCESS_KEY_NAME" default:"guest"`
	AccessKey     string `envconfig:"AMQP_ACCESS_KEY" default:"password"`
}

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		os.Exit(1)
	}
	t, err := amqp.New(env.AMQPServer, env.Queue,
		amqp.WithConnOpt(qp.ConnSASLPlain(env.AccessKeyName, env.AccessKey)),
	)
	if err != nil {
		log.Fatalf("failed to create amqp transport, %s", err.Error())
	}
	t.Encoding = amqp.BinaryV03
	c, err := cloudevents.NewClient(t, cloudevents.WithTimeNow(), cloudevents.WithUUIDs())
	if err != nil {
		log.Fatalf("failed to create client, %s", err.Error())
	}

	event := cloudevents.NewEvent(cloudevents.VersionV03)
	event.SetSource("kn")
	event.SetType(events.ResetType)

	if _, err := c.Send(context.Background(), event); err != nil {
		log.Fatalf("failed to send: %v", err)
	}
}

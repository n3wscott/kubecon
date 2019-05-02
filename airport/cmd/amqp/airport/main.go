package main

import (
	"context"
	"github.com/n3wscott/kubecon/airport/pkg/airport"
	"github.com/n3wscott/kubecon/airport/pkg/cache"
	"log"
	"os"
	"strings"

	"github.com/bradfitz/gomemcache/memcache"
	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/transport/amqp"
	"github.com/kelseyhightower/envconfig"
	qp "pack.ag/amqp"
)

type envConfig struct {
	// AMQPServer URL to connect to the amqp server.
	AMQPServer string `envconfig:"AMQP_SERVER" default:"amqp://localhost:5672/" required:"true"`

	// Queue is the amqp queue name to publish cloudevents on.
	Queue string `envconfig:"AMQP_QUEUE"`

	AccessKeyName string `envconfig:"AMQP_ACCESS_KEY_NAME" default:"guest"`
	AccessKey     string `envconfig:"AMQP_ACCESS_KEY" default:"password"`

	Role string `envconfig:"AIRPORT_ROLE" default:"barista"`

	MemcacheServers string `envconfig:"MEMCACHE_SERVERS" default:"localhost:11211"` // comma separated.
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

	mc := memcache.New(strings.Split(env.MemcacheServers, ",")...)

	log.Println("Starting as a", env.Role)

	a := airport.NewKnAirport(c, cache.NewCache(mc), env.Role)
	if err := a.Start(context.Background()); err != nil {
		log.Printf("start returned error, %v", err)
		os.Exit(1)
	}
	os.Exit(0)
}

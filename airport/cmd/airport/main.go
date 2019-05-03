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
	"github.com/kelseyhightower/envconfig"
)

type envConfig struct {
	Role              string `envconfig:"AIRPORT_ROLE" default:"barista"`
	Broker            string `envconfig:"BROKER" default:"http://localhost:8080/" required:"true"`
	SinkAccessKeyName string `envconfig:"SINK_ACCESS_KEY_NAME" default:"sak" required:"true"`
	SinkAccessKey     string `envconfig:"SINK_ACCESS_KEY" default:"sak" required:"true"`
	MemcacheServers   string `envconfig:"MEMCACHE_SERVERS" default:"localhost:11211" required:"true"` // comma separated.
}

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		os.Exit(1)
	}

	t, err := cloudevents.NewHTTPTransport(
		cloudevents.WithTarget(env.Broker),
		cloudevents.WithBinaryEncoding(),
	)
	if err != nil {
		log.Fatalf("failed to create transport, %s", err.Error())
	}

	c, err := cloudevents.NewClient(t, cloudevents.WithTimeNow(), cloudevents.WithUUIDs())
	if err != nil {
		log.Fatalf("failed to create client, %s", err.Error())
	}

	mc := memcache.New(strings.Split(env.MemcacheServers, ",")...)

	log.Println("Starting as a", env.Role)

	a := airport.NewKnAirport(c, cache.NewCache(mc), env.Role, env.SinkAccessKeyName, env.SinkAccessKey)

	if err := a.Start(context.Background()); err != nil {
		log.Printf("start returned error, %v", err)
		os.Exit(1)
	}
	os.Exit(0)
}

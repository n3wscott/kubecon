package airport

import (
	"context"
	"github.com/cloudevents/sdk-go"
	"github.com/n3wscott/kubecon/airport/pkg/events"
	"log"
)

type Truck struct {
	ConnectedRole

	provider string
	name     string
}

func (a *Truck) Connect() {
	if a.Role != TruckRole {
		return
	}

	event := cloudevents.NewEvent(cloudevents.VersionV03)
	event.SetType(events.ConnectionType)
	event.SetSource(a.provider)

	data := events.ConnectionData{
		System:       a.provider,
		Organization: a.name,
		Logo:         events.KnLogo,
	}

	if err := event.SetData(data); err != nil {
		log.Fatalf("failed to set data, %s", err.Error())
	}

	if _, err := a.Client.Send(context.Background(), event); err != nil {
		log.Fatalf("failed to send: %v", err)
	}
}

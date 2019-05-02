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

func (a *Truck) Receive(event cloudevents.Event) {
	//fmt.Printf("CloudEvent:\n%s", event)
	//
	//fmt.Printf("----------------------------\n")

	switch event.Type() {
	case events.TransferActionType:
		data := &events.TransferActionData{}

		if err := event.DataAs(data); err != nil {
			log.Printf("failed to get transfer data, %v", err)
		}

		if data.ToLocation == a.provider {
			switch data.ActionStatus {
			case events.ActionStatusPotential:
				log.Println("More", data.Offer, "on the way!")
			case events.ActionStatusArrived:
				a.ShipmentArrived(event, data)
			case events.ActionStatusCompleted:
				a.UpdateOfferLevel(event.ID(), data.Offer)
			}
		}
	}
}

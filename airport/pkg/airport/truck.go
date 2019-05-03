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
	event.SetExtension(a.SinkAccessKeyName, a.SinkAccessKey)

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
	switch event.Type() {
	case events.ResetType:
		log.Println("Truck,", a.provider, "resetting.")
		a.Connect()

	case events.DisconnectType:
		if event.Subject() == a.provider {
			a.Connect()
		}

	case events.TransferActionType:
		data := &events.TransferActionData{}

		if err := event.DataAs(data); err != nil {
			log.Printf("failed to get transfer data, %v", err)
		}

		// see if we service this.
		if !a.inRoute(data) {
			return
		}

		switch data.ActionStatus {
		case events.ActionStatusPotential:
			a.PickUpShipment(event, data)

		case events.ActionStatusArrived:
			a.DeliverShipment(event, data)

		}
	}
}

func (a *Truck) inRoute(data *events.TransferActionData) bool {
	for _, route := range a.Cache.GetCarrierRoute(a.provider) {
		if route.ToLocation == data.ToLocation && route.FromLocation == data.FromLocation {
			log.Println("In route:", data.ToLocation, "-->", data.FromLocation)
			return true
		}
	}
	log.Println("Not in route:", data.ToLocation, "-->", data.FromLocation)
	return false
}

func (a *Truck) PickUpShipment(from cloudevents.Event, box *events.TransferActionData) {
	if a.Role != TruckRole {
		return
	}

	event := cloudevents.NewEvent(cloudevents.VersionV03)
	event.SetType(events.TransferActionType)
	event.SetSource(a.provider)
	event.SetSubject(from.Subject())
	event.SetExtension(events.ExtCause, from.ID())
	event.SetExtension(a.SinkAccessKeyName, a.SinkAccessKey)

	data := events.TransferActionData{
		ActionStatus: events.ActionStatusActive,
		ToLocation:   box.ToLocation,
		FromLocation: box.FromLocation,
		Offer:        box.Offer,
	}

	log.Println("Picked up", data.Offer, "and going", data.ToLocation, "-->", data.FromLocation)

	if err := event.SetData(data); err != nil {
		log.Fatalf("failed to set data, %s", err.Error())
	}

	if _, err := a.Client.Send(context.Background(), event); err != nil {
		log.Fatalf("failed to send: %v", err)
	}
}

func (a *Truck) DeliverShipment(from cloudevents.Event, box *events.TransferActionData) {
	if a.Role != TruckRole {
		return
	}

	event := cloudevents.NewEvent(cloudevents.VersionV03)
	event.SetType(events.TransferActionType)
	event.SetSource(a.provider)
	event.SetSubject(from.Subject())
	event.SetExtension(events.ExtCause, from.ID())
	event.SetExtension(a.SinkAccessKeyName, a.SinkAccessKey)

	data := events.TransferActionData{
		ActionStatus: events.ActionStatusCompleted,
		ToLocation:   box.ToLocation,
		FromLocation: box.FromLocation,
		Offer:        box.Offer,
	}

	log.Println("Delivered", data.Offer, "with route", data.ToLocation, "-->", data.FromLocation)

	if err := event.SetData(data); err != nil {
		log.Fatalf("failed to set data, %s", err.Error())
	}

	if _, err := a.Client.Send(context.Background(), event); err != nil {
		log.Fatalf("failed to send: %v", err)
	}
}

package airport

import (
	"context"
	"github.com/cloudevents/sdk-go"
	"github.com/google/uuid"
	"github.com/n3wscott/kubecon/airport/pkg/events"
	"log"
	"strings"
)

type Warehouse struct {
	ConnectedRole

	provider string
	name     string
}

func (a *Warehouse) Connect() {
	if a.Role != WarehouseRole {
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

func (a *Warehouse) Receive(event cloudevents.Event) {
	switch event.Type() {
	case events.ResetType:
		log.Println("Warehouse,", a.provider, "resetting.")
		a.Connect()

	case events.DisconnectType:
		if event.Subject() == a.provider {
			a.Connect()
		}

	case events.OrderType:
		data := &events.OrderData{}
		if err := event.DataAs(data); err != nil {
			log.Printf("failed to get order data, %v", err)
		}

		// we only service retail
		if !strings.HasPrefix(data.Customer, events.RetailerPrefix) {
			return
		}

		// see if we stock this.
		if !a.inStock(data) {
			return
		}

		switch data.OrderStatus {
		case events.OrderReleased:
			a.ShipOrder(event, data)

			//case events.OrderDelivered:
			//a.DeliverShipment(event, data)

		}
	}
}

func (a *Warehouse) inStock(data *events.OrderData) bool {
	for _, stock := range a.Cache.GetWarehouseOffers(a.provider) {
		if stock.Customer == data.Customer {
			for _, have := range stock.Offer {
				if data.Offer == have {
					log.Println("Have stock:", data.Offer, "for", data.Customer)
					return true
				}
			}
		}
	}
	log.Println("Not stocked:", data.Offer, "for", data.Customer)
	return false
}

func (a *Warehouse) ShipOrder(from cloudevents.Event, order *events.OrderData) {
	if a.Role != WarehouseRole {
		return
	}

	event := cloudevents.NewEvent(cloudevents.VersionV03)
	event.SetType(events.TransferActionType)
	event.SetSource(a.provider)
	event.SetSubject(uuid.New().String()) // TODO?
	event.SetExtension(events.ExtCause, from.ID())

	data := events.TransferActionData{
		ActionStatus: events.ActionStatusPotential,
		FromLocation: a.provider,
		ToLocation:   order.Customer,
		Offer:        order.Offer,
	}

	log.Println("Shipping", data.Offer, "with route", data.ToLocation, "-->", data.FromLocation)

	if err := event.SetData(data); err != nil {
		log.Fatalf("failed to set data, %s", err.Error())
	}

	if _, err := a.Client.Send(context.Background(), event); err != nil {
		log.Fatalf("failed to send: %v", err)
	}
}

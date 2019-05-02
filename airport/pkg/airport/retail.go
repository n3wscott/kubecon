package airport

import (
	"context"
	"github.com/cloudevents/sdk-go"
	"github.com/google/uuid"
	"github.com/n3wscott/kubecon/airport/pkg/events"
	"log"
	"time"
)

type Retail struct {
	ConnectedRole

	provider string
	name     string
}

func (a *Retail) Receive(event cloudevents.Event) {
	//fmt.Printf("CloudEvent:\n%s", event)
	//
	//fmt.Printf("----------------------------\n")

	switch event.Type() {
	case events.OrderType:
		switch event.Source() {
		case events.PassengerSource:
			a.HandleOrder(event)
		case a.provider:
			a.HandleStock(event)
		case events.ControllerSource:
			log.Println("Controller ordered something... TODO.")
		}

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

func (a *Retail) Connect() {
	if a.Role != BaristaRole {
		return
	}

	a.Cache.SetProductCount(a.provider, events.SmallProduct, events.ShipmentCount)
	a.Cache.SetProductCount(a.provider, events.MediumProduct, events.ShipmentCount)
	a.Cache.SetProductCount(a.provider, events.LargeProduct, events.ShipmentCount)

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

func (a *Retail) HandleOrder(event cloudevents.Event) {
	order := &events.OrderData{}
	if err := event.DataAs(order); err != nil {
		log.Printf("failed to get order as data, %v", err)
	}
	if order.Provider == a.provider {
		// it is for us! yay messaging!
		switch order.OrderStatus {
		case events.OrderReleased:
			a.DeliverOrder(event.ID(), order)
		}
	}
}

var served int

func (a *Retail) DeliverOrder(cause string, order *events.OrderData) {
	if a.Role != BaristaRole {
		return
	}

	slept := time.Duration(0)
	for i := 0; true; i++ {
		// give up?
		if slept.Seconds() > 10 {
			log.Println("Giving up on ", order.Offer, ", ", cause)
			return
		}

		count := a.Cache.GetProductCount(a.provider, order.Offer)
		if count > 0 {
			count = a.Cache.AdjustProductCount(a.provider, order.Offer, -1)
			break
		}
		nap := 250 * time.Millisecond
		time.Sleep(nap)
		slept += nap
	}
	served++ // DEBUG
	log.Println("Serving a", order.Offer, ". #", served)

	event := cloudevents.NewEvent(cloudevents.VersionV03)
	event.SetType(events.OrderType)
	event.SetSource(a.provider)
	event.SetExtension(events.ExtCause, cause)
	event.SetSubject(order.Customer)

	data := events.OrderData{
		Provider:    a.provider,
		OrderStatus: events.OrderDelivered,
		Customer:    order.Customer,
		Offer:       order.Offer,
	}

	if err := event.SetData(data); err != nil {
		log.Fatalf("failed to set data, %s", err.Error())
	}

	if _, err := a.Client.Send(context.Background(), event); err != nil {
		log.Fatalf("failed to send: %v", err)
	}
}

func (a *Retail) HandleStock(event cloudevents.Event) {
	if a.Role != InventoryRole {
		return
	}

	order := &events.OrderData{}
	if err := event.DataAs(order); err != nil {
		log.Printf("failed to get order as data, %v", err)
	}
	if order.Provider == a.provider {
		// it is for us! yay messaging!
		switch order.OrderStatus {
		case events.OrderDelivered:
			a.UpdateOfferLevel(event.ID(), order.Offer)
		}
	}
}

func (a *Retail) UpdateOfferLevel(cause string, offer events.Product) {
	if a.Role != InventoryRole {
		return
	}

	event := cloudevents.NewEvent(cloudevents.VersionV03)
	event.SetType(events.OfferType)
	event.SetSource(a.provider)
	event.SetSubject(string(offer))
	event.SetExtension(events.ExtCause, cause)

	count := a.Cache.GetProductCount(a.provider, offer)
	data := events.OfferData{
		InventoryLevel: count,
		Offer:          offer,
	}

	log.Println("Inventory of", data.Offer, "at", data.InventoryLevel)

	if err := event.SetData(data); err != nil {
		log.Fatalf("failed to set data, %s", err.Error())
	}

	if _, err := a.Client.Send(context.Background(), event); err != nil {
		log.Fatalf("failed to send: %v", err)
	}

	if data.InventoryLevel == 0 {
		a.OrderMore(offer)
	}
}

func (a *Retail) OrderMore(offer events.Product) {
	if a.Role != InventoryRole {
		return
	}

	event := cloudevents.NewEvent(cloudevents.VersionV03)
	event.SetType(events.OrderType)
	event.SetSource(a.provider)
	event.SetSubject(uuid.New().String()) // subject is random uuid?

	data := events.OrderData{
		OrderStatus: events.OrderReleased,
		Customer:    a.provider,
		Offer:       offer,
	}

	log.Println("Asking for more", data.Offer)

	if err := event.SetData(data); err != nil {
		log.Fatalf("failed to set data, %s", err.Error())
	}

	if _, err := a.Client.Send(context.Background(), event); err != nil {
		log.Fatalf("failed to send: %v", err)
	}
}

func (a *Retail) ShipmentArrived(from cloudevents.Event, shipment *events.TransferActionData) {
	if a.Role != InventoryRole {
		return
	}

	log.Println("More", shipment.Offer, "arrived for", shipment.ToLocation, "!")

	if shipment.ToLocation == a.provider {
		a.Cache.AdjustProductCount(shipment.ToLocation, shipment.Offer, events.ShipmentCount)
	}
}

// TODO: for carrier
//func (a *Retail) ShipmentArrived(from cloudevents.Event, shipment *events.TransferActionData) {
//	if a.Role != InventoryRole {
//		return
//	}
//
//	log.Println("More", shipment.Offer, "arrived!")
//
//	event := cloudevents.NewEvent(cloudevents.VersionV03)
//	event.SetType(events.TransferActionType)
//	event.SetSource(a.provider)
//	event.SetSubject(from.Subject()) // subject is random uuid?
//
//	data := events.TransferActionData{
//		ActionStatus: events.ActionStatusCompleted,
//		ToLocation:   shipment.ToLocation,
//		FromLocation: shipment.FromLocation,
//		Offer:        shipment.Offer,
//	}
//
//	log.Println("Accepting shipment of ", data.Offer)
//
//	if err := event.SetData(data); err != nil {
//		log.Fatalf("failed to set data, %s", err.Error())
//	}
//
//	if _, err := a.Client.Send(context.Background(), event); err != nil {
//		log.Fatalf("failed to send: %v", err)
//	}
//}

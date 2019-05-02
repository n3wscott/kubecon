package airport

import (
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/n3wscott/kubecon/airport/pkg/events"
)

type Director struct {
	ConnectedRole

	providers []string
}

func (a *Director) Receive(event cloudevents.Event) {

	// we only listen to events from the controller.

	if event.Source() == events.ControllerSource {
		fmt.Printf("----------------------------\n")
		switch event.Type() {
		case events.ResetType:
			fmt.Println("Resetting the cache.")
			if err := a.Cache.Reset(); err != nil {
				fmt.Println("failed to reset the cache,", err)
			}

		case events.DisconnectType:
			//  subject: Retailer.kn

		case events.ProductOfferType:
			data := make(events.CustomerOfferData, 0)
			if err := event.DataAs(&data); err != nil {
				fmt.Printf("failed to get customer offer data, %s", err)
				return
			}
			fmt.Println("For warehouse", event.Subject())

			for _, stock := range data {
				for _, inv := range stock.Offer {
					fmt.Println("Stock", inv, "for", stock.Customer)
				}
			}
			// Store this in the cache.
			a.Cache.SetWarehouseOffers(event.Subject(), data)

		case events.TransferActionType:
			// ignore.

		case events.CarrierOfferType:
			data := make(events.CarrierOfferData, 0)
			if err := event.DataAs(&data); err != nil {
				fmt.Printf("failed to get carrier offer data, %s", err)
				return
			}
			fmt.Println("For truck", event.Subject())
			for _, move := range data {
				fmt.Println("Move", move.FromLocation, "->", move.ToLocation)
			}
			// Store this in the cache.
			a.Cache.SetCarrierRoute(event.Subject(), data)

		default:
			fmt.Printf("Unhandled CloudEvent:\n%s", event)
		}

	}

}

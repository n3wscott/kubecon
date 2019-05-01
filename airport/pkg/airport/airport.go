package airport

import (
	"context"
	"github.com/cloudevents/sdk-go"
	"github.com/n3wscott/kubecon/airport/pkg/events"
)

const (
	BaristaRole   = "barista"
	InventoryRole = "stocker"
	WarehouseRole = "warehouse"
	TruckRole     = "truck"
)

type ConnectedRole struct {
	Role   string
	Client cloudevents.Client
}

type Airport struct {
	ConnectedRole

	retail    *Retail
	warehouse *Warehouse
	truck     *Truck
}

func NewKnAirport(client cloudevents.Client, role string) *Airport {
	a := &Airport{
		ConnectedRole: ConnectedRole{
			Client: client,
			Role:   role,
		},
	}
	return a
}

func (a *Airport) Start(ctx context.Context) error {
	switch a.Role {
	case BaristaRole, InventoryRole:
		a.retail = &Retail{
			ConnectedRole: a.ConnectedRole,
			name:          "Knative Coffee",
			provider:      events.RetailerPrefix + "kn",
		}
		a.retail.Connect()

	case WarehouseRole:
		a.warehouse = &Warehouse{
			ConnectedRole: a.ConnectedRole,
			name:          "Knative Warehouse",
			provider:      events.SupplierPrefix + "kn",
		}
		a.warehouse.Connect()

	case TruckRole:
		a.truck = &Truck{
			ConnectedRole: a.ConnectedRole,
			name:          "Knative Trucking",
			provider:      events.CarrierPrefix + "kn",
		}
		a.truck.Connect()
	}

	return a.Client.StartReceiver(ctx, a.Receive)
}

func (a *Airport) Receive(event cloudevents.Event) {
	//fmt.Printf("CloudEvent:\n%s", event)
	//
	//fmt.Printf("----------------------------\n")

	if a.retail != nil {
		a.retail.Receive(event)
	}

}

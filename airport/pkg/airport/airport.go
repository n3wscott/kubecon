package airport

import (
	"context"
	"github.com/cloudevents/sdk-go"
	"github.com/n3wscott/kubecon/airport/pkg/cache"
	"github.com/n3wscott/kubecon/airport/pkg/events"
)

const (
	BaristaRole   = "barista"
	InventoryRole = "stocker"
	WarehouseRole = "warehouse"
	TruckRole     = "truck"
	DirectorRole  = "director"

	name = "knative"
)

type ConnectedRole struct {
	Role              string
	Client            cloudevents.Client
	Cache             cache.Cache
	SinkAccessKeyName string
	SinkAccessKey     string
}

type Airport struct {
	ConnectedRole

	retail    *Retail
	warehouse *Warehouse
	truck     *Truck
	director  *Director
}

func NewKnAirport(client cloudevents.Client, store cache.Cache, role, sinkAccessKeyName, sinkAccessKey string) *Airport {
	a := &Airport{
		ConnectedRole: ConnectedRole{
			Client:            client,
			Role:              role,
			Cache:             store,
			SinkAccessKeyName: sinkAccessKeyName,
			SinkAccessKey:     sinkAccessKey,
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
			provider:      events.RetailerPrefix + name,
		}
		a.retail.Connect()

	case WarehouseRole:
		a.warehouse = &Warehouse{
			ConnectedRole: a.ConnectedRole,
			name:          "Knative Warehouse",
			provider:      events.SupplierPrefix + name,
		}
		a.warehouse.Connect()

	case TruckRole:
		a.truck = &Truck{
			ConnectedRole: a.ConnectedRole,
			name:          "Knative Trucking",
			provider:      events.CarrierPrefix + name,
		}
		a.truck.Connect()

	case DirectorRole:
		a.director = &Director{
			ConnectedRole: a.ConnectedRole,
			providers: []string{
				events.RetailerPrefix + name,
				events.SupplierPrefix + name,
				events.CarrierPrefix + name,
			},
		}
	}

	return a.Client.StartReceiver(ctx, a.Receive)
}

func (a *Airport) Receive(event cloudevents.Event) {
	if a.retail != nil {
		go a.retail.Receive(event)
	}

	if a.warehouse != nil {
		go a.warehouse.Receive(event)
	}

	if a.truck != nil {
		go a.truck.Receive(event)
	}

	if a.director != nil {
		go a.director.Receive(event)
	}

}

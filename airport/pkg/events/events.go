package events

const (
	KnLogo         = "https://raw.githubusercontent.com/knative/docs/master/docs/images/logo/rgb/knative-logo-rgb.png"
	RetailerPrefix = "Retailer."
	SupplierPrefix = "Supplier."
	CarrierPrefix  = "Carrier."

	// goes on type
	OfferPrefix = "Offer."

	// goes on subject
	SupplierSubject = "Supplier."

	ControllerSource = "Controller" // the main controller emits these for type:Offer.Product type:Reset type:Disconnect type:Offer.Service.Transport
	PassengerSource  = "Passenger"  // the user emits these for type:Order

	ConnectionType = "Connection"
	ResetType      = "Reset"
	DisconnectType = "Disconnect"

	OrderType          = "Order"
	OrderReleasedType  = "Order.OrderStatus.OrderReleased"
	OrderDeliveredType = "Order.OrderStatus.OrderDelivered"

	TransferActionType         = "TransferAction"
	TransferOrderReleasedType  = "TransferAction.ActionStatus.PotentalActionStatus"
	TransferOrderAcceptedType  = "TransferAction.ActionStatus.ActiveActionStatus"
	TransferOrderArrivedType   = "TransferAction.ActionStatus.ArrivedActionStatus"
	TransferOrderCompletedType = "TransferAction.ActionStatus.CompletedActionStatus"

	OfferType          = "Offer"
	InventoryLevelType = "Offer.InventoryLevel"
	ProductOfferType   = "Offer.Product"
	CarrierOfferType   = "Offer.Service.Transport"

	ShipmentCount = 2

	// OrderStatus can be
	OrderReleased  = "OrderReleased"
	OrderDelivered = "OrderDelivered"

	// ActionStatus can be
	ActionStatusCompleted = "CompletedActionStatus"
	ActionStatusArrived   = "ArrivedActionStatus"
	ActionStatusActive    = "ActiveActionStatus"
	ActionStatusPotential = "PotentialActionStatus"

	// Extensions
	ExtCause = "cause"
)

type Product string

const (
	SmallProduct  Product = "small"
	MediumProduct Product = "medium"
	LargeProduct  Product = "large"
)

// connection prefix goes on system and source.

type ConnectionData struct {
	System       string `json:"system,omitempty"`
	Organization string `json:"organization,omitempty"`
	Logo         string `json:"logo,omitempty"`
}

type OfferData struct {
	InventoryLevel int     `json:"inventoryLevel,omitempty"`
	Offer          Product `json:"offer,omitempty"`
}

type ProductOffer struct {
	Customer string    `json:"customer,omitempty"`
	Offer    []Product `json:"offer,omitempty"`
}

type CustomerOfferData []ProductOffer

type CarrierOffer struct {
	ToLocation   string `json:"toLocation,omitempty"`
	FromLocation string `json:"fromLocation,omitempty"`
}

type CarrierOfferData []CarrierOffer

type OrderData struct {
	Provider    string  `json:"provider,omitempty"`
	OrderStatus string  `json:"orderStatus,omitempty"`
	Customer    string  `json:"customer,omitempty"`
	Offer       Product `json:"offer,omitempty"`
}

type TransferActionData struct {
	ToLocation   string  `json:"toLocation,omitempty"`
	FromLocation string  `json:"fromLocation,omitempty"`
	ActionStatus string  `json:"actionStatus,omitempty"`
	Offer        Product `json:"offer,omitempty"`
}

package events

const (
	KnLogo         = "https://raw.githubusercontent.com/knative/docs/master/docs/images/logo/rgb/knative-logo-rgb.png"
	RetailerPrefix = "Retailer."
	SupplierPrefix = "Supplier."
	CarrierPrefix  = "Carrier."

	ConnectionType = "Connection"
)

// connection prefix goes on system and source.

type ConnectionData struct {
	System       string `json:"system,omitempty"`
	Organization string `json:"organization,omitempty"`
	Logo         string `json:"logo,omitempty"`
}

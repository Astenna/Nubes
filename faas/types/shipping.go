package types

type ShippingState string

const (
	InPreparation   ShippingState = "InPreparation"
	InTransit       ShippingState = "InTransit"
	PickupAvailavle ShippingState = "PickupAvailable"
	Delivered       ShippingState = "Delivered"
)

type Shipping struct {
	Id            string
	Address       string
	State         ShippingState
	isInitialized bool
}

func (s Shipping) GetTypeName() string {
	return "Shipping"
}
func (receiver *Shipping) Init() {
	receiver.isInitialized = true
}

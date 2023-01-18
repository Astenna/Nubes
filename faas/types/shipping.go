package types

type ShippingState string

const (
	InPreparation   ShippingState = "InPreparation"
	InTransit       ShippingState = "InTransit"
	PickupAvailavle ShippingState = "PickupAvailable"
	Delivered       ShippingState = "Delivered"
)

type Shipping struct {
	Id      string
	Address string
	State   ShippingState
}

func (s Shipping) GetTypeName() string {
	return "Shipping"
}

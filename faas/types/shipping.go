package types

import "github.com/Astenna/Nubes/lib"

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

func NewShipping(shipping Shipping) (Shipping, error) {
	out, _libError := lib.Insert(shipping)
	if _libError != nil {
		return *new(Shipping), _libError
	}
	shipping.Id = out
	return shipping, nil
}

func ReNewShipping(id string) Shipping {
	shipping := new(Shipping)
	shipping.Id = id
	return *shipping
}

func (s Shipping) GetTypeName() string {
	return "Shipping"
}

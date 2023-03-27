package types

import "github.com/Astenna/Nubes/lib"

type ShippingState string

const (
	InPreparation	ShippingState	= "InPreparation"
	InTransit	ShippingState	= "InTransit"
	PickupAvailavle	ShippingState	= "PickupAvailable"
	Delivered	ShippingState	= "Delivered"
)

type Shipping struct {
	Id		string
	Address		string
	State		ShippingState
	isInitialized	bool
	invocationDepth	int
}

func (s Shipping) GetTypeName() string {
	return "Shipping"
}
func (receiver *Shipping) Init() {
	receiver.isInitialized = true
}
func (receiver *Shipping) saveChangesIfInitialized() error {
	if receiver.isInitialized && receiver.invocationDepth == 1 {
		_libError := lib.Upsert(receiver, receiver.Id)
		if _libError != nil {
			return _libError
		}
	}
	return nil
}

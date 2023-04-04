package types

import (
	"time"

	"github.com/Astenna/Nubes/lib"
)

type ShippingState string

const (
	InPreparation   ShippingState = "InPreparation"
	InTransit       ShippingState = "InTransit"
	PickupAvailavle ShippingState = "PickupAvailable"
	Delivered       ShippingState = "Delivered"
)

type Shipping struct {
	Id              string
	Address         string
	State           ShippingState
	CreationDate    time.Time
	isInitialized   bool
	invocationDepth int
}

func (s Shipping) GetTypeName() string {
	return "Shipping"
}

// ExportShipping is an example of custom export implementation
// note that, after preparation of an object
// the method contains an invocation of lib.Export.
// Moreover, the methods signature follows the required converntion:
// func Export<type-name>(param <input-type>) (string, error)
// where <input-type> is an arbitrary type chosen according to the needs.
func ExportShipping(addr string) (string, error) {
	newShipping := Shipping{
		Address:      addr,
		State:        ShippingState(InPreparation),
		CreationDate: time.Now(),
	}
	exported, err := lib.Export[Shipping](newShipping)
	return exported.Id, err
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

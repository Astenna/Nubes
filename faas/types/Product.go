package types

import (
	"errors"

	"github.com/Astenna/Nubes/lib"
)

type Product struct {
	Id                string
	Name              string
	QuantityAvailable int
	SoldBy            lib.FaasReference[Shop]
	Price             float32
}

func (Product) GetTypeName() string {
	return "Product"
}

func (p *Product) DecreaseAvailabilityBy(decreaseNum int) error {
	if p.QuantityAvailable-decreaseNum < 0 {
		return errors.New("not enough quantity available")
	}

	p.QuantityAvailable = p.QuantityAvailable - decreaseNum
	return nil
}

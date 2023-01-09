package types

import (
	"errors"

	"github.com/Astenna/Nubes/lib"
)

type Product struct {
	Id                string
	Name              string
	QuantityAvailable int
	SoldBy            lib.Reference[Shop]
	Price             float32
}

func (Product) GetTypeName() string {
	return "Product"
}

func (p *Product) DecreaseAvailabilityBy(decreaseNum int) error {
	p, _libError := lib.Get[Product](p.Id)
	if _libError != nil {
		return _libError
	}
	if p.QuantityAvailable-decreaseNum < 0 {
		return errors.New("not enough quantity available")
	}

	p.QuantityAvailable = p.QuantityAvailable - decreaseNum
	_libError = lib.Upsert(p, p.Id)
	if _libError != nil {
		return _libError
	}
	return nil
}

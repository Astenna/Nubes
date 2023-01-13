package types

import (
	"errors"

	"github.com/Astenna/Nubes/lib"
)

type Product struct {
	Id			string
	Name			string
	QuantityAvailable	int
	SoldBy			lib.Reference[Shop]
	Price			float32
}

func NewProduct(product Product) (Product, error) {
	out, _libError := lib.Insert(product)
	if _libError != nil {
		return *new(Product), _libError
	}
	product.Id = out
	return product, nil
}

func ReNewProduct(id string) (Product, error) {
	product := new(Product)
	product.Id = id
	return *product, nil
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

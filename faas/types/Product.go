package types

import "github.com/Astenna/Nubes/lib"

type Product struct {
	Id                int
	Name              string
	QuantityAvailable int
	SoldBy            lib.Reference[Shop]
	Price             float32
}

func (Product) GetTypeName() string {
	return "Product"
}

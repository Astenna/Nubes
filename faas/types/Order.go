package types

import "github.com/Astenna/Nubes/lib"

type Order struct {
	Id       string
	Buyer    lib.Reference[User]
	Products []OrderedProduct
	Shipping lib.Reference[Shipping]
}

type OrderedProduct struct {
	Product  lib.Reference[Product]
	Quantity int
}

func (Order) GetTypeName() string {
	return "Order"
}

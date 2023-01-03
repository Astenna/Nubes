package types

import "github.com/Astenna/Nubes/lib"

type Order struct {
	Id       string
	Buyer    lib.FaasReference[User]
	Products []OrderedProduct
	Shipping lib.FaasReference[Shipping]
}

type OrderedProduct struct {
	Product  lib.FaasReference[Product]
	Quantity int
}

func (Order) GetTypeName() string {
	return "Order"
}

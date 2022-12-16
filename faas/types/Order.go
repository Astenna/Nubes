package types

import "github.com/Astenna/Nubes/lib"

type Order struct {
	Id       int
	Buyer    lib.Reference[User]
	Products []OrderedProduct
}

type OrderedProduct struct {
	Product  lib.Reference[Product]
	Quantity int
}

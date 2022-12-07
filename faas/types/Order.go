package types

import "github.com/Astenna/Thesis_PoC/faas_lib"

type Order struct {
	Id       int
	Buyer    faas_lib.Reference[User]
	Products []OrderedProduct
}

type OrderedProduct struct {
	Product  faas_lib.Reference[Product]
	Quantity int
}

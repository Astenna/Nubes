package types

type Order struct {
	Id       int
	Buyer    User
	Products []OrderedProduct
}

type OrderedProduct struct {
	Product  Product
	Quantity int
}

package types

type Order struct {
	id       int
	buyer    User
	products []OrderedProduct
}

type OrderedProduct struct {
	product  Product
	quantity int
}

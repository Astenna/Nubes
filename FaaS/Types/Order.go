package Types

type order struct {
	id       int
	buyer    user
	products []orderedProduct
}

type orderedProduct struct {
	product  product
	quantity int
}

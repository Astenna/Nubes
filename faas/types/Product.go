package types

type Product struct {
	id                int
	name              string
	quantityAvailable int
	soldBy            Shop
	price             float32
}

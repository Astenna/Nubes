package types

import "github.com/Astenna/Thesis_PoC/faas_lib"

type Product struct {
	Id                int
	Name              string
	QuantityAvailable int
	SoldBy            faas_lib.Reference[Shop]
	Price             float32
}

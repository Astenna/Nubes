package client_lib

import (
	"github.com/Astenna/Nubes/lib"
)

type OrderedProduct struct {
	Product lib.Reference[ProductStub]

	Quantity int
}

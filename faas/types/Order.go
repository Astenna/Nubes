package types

import (
	"errors"

	"github.com/Astenna/Nubes/lib"
)

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

func NewOrder(order Order) (Order, error) {

	for _, orderedProduct := range order.Products {
		product, err := orderedProduct.Product.Get()
		if err != nil {
			return *new(Order), errors.New("item " + orderedProduct.Product.Id + " not available")
		}
		product.DecreaseAvailabilityBy(orderedProduct.Quantity)
	}

	buyer, err := order.Buyer.Get()
	if err != nil {
		return *new(Order), errors.New("unable to retrieve user's address for shipping")
	}
	shipping, err := NewShipping(Shipping{
		State:   InPreparation,
		Address: buyer.Address,
	})
	if err != nil {
		return *new(Order), errors.New("failed to create shipping for the order: " + err.Error())
	}

	order.Shipping = *lib.NewReference[Shipping](shipping.Id)
	out, _libError := lib.Insert(order)
	if _libError != nil {
		return *new(Order), _libError
	}
	order.Id = out
	return order, nil
}

func ReNewOrder(id string) Order {
	order := new(Order)
	order.Id = id
	return *order
}

func (o Order) GetTypeName() string {
	return "Order"
}

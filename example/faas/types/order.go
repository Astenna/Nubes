package types

import (
	"errors"

	"github.com/Astenna/Nubes/lib"
)

type Order struct {
	Id            string
	Products      []OrderedProduct
	Buyer         lib.Reference[User]
	Shipping      lib.Reference[Shipping]
	isInitialized bool
}

type OrderedProduct struct {
	Product  lib.Reference[Product]
	Quantity int
}

func NewOrder(order Order) (Order, error) {

	for _, orderedProduct := range order.Products {
		product, err := orderedProduct.Product.Get()
		if err != nil {
			return *new(Order), errors.New("item " + orderedProduct.Product.Id() + " not available")
		}
		product.DecreaseAvailabilityBy(orderedProduct.Quantity)
	}

	buyer, err := order.Buyer.Get()
	if err != nil {
		return *new(Order), errors.New("unable to retrieve user's address for shipping")
	}
	shipping, err := lib.Export[Shipping](Shipping{
		State:   InPreparation,
		Address: buyer.Address,
	})
	if err != nil {
		return *new(Order), errors.New("failed to create shipping for the order: " + err.Error())
	}

	order.Shipping = lib.Reference[Shipping](shipping.Id)
	return order, nil
}

func (o Order) GetTypeName() string {
	return "Order"
}
func (receiver *Order) Init() {
	receiver.isInitialized = true
}

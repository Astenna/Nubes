package types

import (
	"errors"

	"github.com/Astenna/Nubes/lib"
)

type Order struct {
	Id              string
	Products        []OrderedProduct
	Buyer           lib.Reference[User]
	Shipping        lib.Reference[Shipping]
	isInitialized   bool
	invocationDepth int
}

type OrderedProduct struct {
	Product  lib.Reference[Product]
	Quantity int
}

func ExportOrder(order Order) (string, error) {

	for _, orderedProduct := range order.Products {
		product, err := orderedProduct.Product.Get()
		if err != nil {
			return "", errors.New("item " + orderedProduct.Product.Id() + " not available")
		}
		product.DecreaseAvailabilityBy(orderedProduct.Quantity)
	}

	buyer, err := order.Buyer.Get()
	if err != nil {
		return "", errors.New("unable to retrieve user's address for shipping")
	}
	shipping, err := lib.Export[Shipping](Shipping{
		State:   InPreparation,
		Address: buyer.AddressText,
	})
	if err != nil {
		return "", errors.New("failed to create shipping for the order: " + err.Error())
	}

	order.Shipping = lib.Reference[Shipping](shipping.Id)
	exportedOrder, err := lib.Export[Order](order)
	return exportedOrder.Id, err
}

func (o Order) GetTypeName() string {
	return "Order"
}
func (receiver *Order) Init() {
	receiver.isInitialized = true
}
func (receiver *Order) saveChangesIfInitialized() error {
	if receiver.isInitialized && receiver.invocationDepth == 1 {
		_libError := lib.Upsert(receiver, receiver.Id)
		if _libError != nil {
			return _libError
		}
	}
	return nil
}

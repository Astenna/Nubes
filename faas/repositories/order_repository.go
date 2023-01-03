package repositories

import (
	"errors"

	"github.com/Astenna/Nubes/faas/types"
	"github.com/Astenna/Nubes/lib"
)

func CreateOrder(order types.Order) (string, error) {
	// CHECK & DECREASE ORDERED PRODUCTS AVAILABILITY
	for _, orderedProduct := range order.Products {
		err := orderedProduct.Product.Get().DecreaseAvailabilityBy(orderedProduct.Quantity)
		if err != nil {
			return "", errors.New("item " + orderedProduct.Product.Id + " not available")
		}
		// OR: address the consistency issues by making smaller order
		// without not available item, return warning in the error (?)
		// https://yourbasic.org/golang/delete-element-slice/
	}

	// CREATE SHIPPING
	newShipping := types.Shipping{
		State:   types.InPreparation,
		Address: order.Buyer.Get().Address,
	}
	newShippingId, err := lib.Insert(newShipping)
	if err != nil {
		return "", errors.New("failed to create shipping for the order: " + err.Error())
	}

	// CREATE ORDER
	order.Shipping = *lib.NewReference[types.Shipping](newShippingId)
	return lib.Insert(order)
}

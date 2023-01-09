package repositories

import (
	"errors"

	"github.com/Astenna/Nubes/faas/types"
	"github.com/Astenna/Nubes/lib"
)

func CreateOrder(order types.Order) (string, error) {
	// CHECK & DECREASE ORDERED PRODUCTS AVAILABILITY
	for _, orderedProduct := range order.Products {
		product, err := orderedProduct.Product.Get()
		if err != nil {
			return "", errors.New("item " + orderedProduct.Product.Id + " not available")
		}
		product.DecreaseAvailabilityBy(orderedProduct.Quantity)
		// OR: address the consistency issues by making smaller order
		// without not available item, return warning in the error (?)
		// https://yourbasic.org/golang/delete-element-slice/
	}

	// CREATE SHIPPING
	buyer, err := order.Buyer.Get()
	if err != nil {
		return "", errors.New("unable to retrieve user's address for shipping")
	}
	newShipping := types.Shipping{
		State:   types.InPreparation,
		Address: buyer.Address,
	}
	newShippingId, err := lib.Insert(newShipping)
	if err != nil {
		return "", errors.New("failed to create shipping for the order: " + err.Error())
	}

	// CREATE ORDER
	order.Shipping = *lib.NewReference[types.Shipping](newShippingId)
	return lib.Insert(order)
}

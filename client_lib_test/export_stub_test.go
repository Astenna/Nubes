package client_lib_test

import (
	"testing"

	clib "github.com/Astenna/Nubes/client_lib"
	"github.com/stretchr/testify/require"

	"github.com/google/uuid"
)

func TestClib(t *testing.T) {
	// Arrange
	existingOrderId := "d192eeda-e709-4415-bbe2-cb91a4968962"
	newId := uuid.New().String()
	newUser := clib.UserStub{
		Email:     newId,
		FirstName: "Kinga",
		LastName:  "Marek",
		Password:  "password",
		Orders:    clib.ReferenceList[clib.OrderStub]{existingOrderId},
	}
	newOrdersSet := []string{"i'm invalid", existingOrderId}

	// Act
	exportedUser, err := clib.ExportUser(newUser)
	require.Equal(t, err, nil, "error occurred in ExportUser invocation", err)

	loadTheSameUser, err := clib.LoadUser(newId)
	require.Equal(t, err, nil, "error occurred in LoadUser invocation", err)

	err = loadTheSameUser.SetOrders(newOrdersSet)
	require.Equal(t, err, nil, "error occurred in SetOrders invocation", err)

	retrievedOrders, err := exportedUser.GetOrdersIds()
	require.Equal(t, err, nil, "error occurred in GetOrders invocation", err)

	// Assert
	require.Equal(t, len(newOrdersSet), len(retrievedOrders), "number of orders is not equal, expected equal")
	require.Equal(t, newOrdersSet[0], retrievedOrders[0], "expected the same orders id, found diffrent")
	require.Equal(t, newOrdersSet[1], retrievedOrders[1], "expected the same orders id, found diffrent")
}

func TestReferenceNavigationListOneToMany(t *testing.T) {
	// Arrange
	newShop := clib.ShopStub{
		Name: "TestReferenceNavigationListShop",
	}
	newProduct := clib.ProductStub{
		Name:              "TestReferenceNavigationListProduct",
		QuantityAvailable: 400.0,
		Price:             3.5,
	}

	// Act
	exportedShop, err := clib.ExportShop(newShop)
	require.Equal(t, nil, err, "error occurred in ExportShop invocation", err)
	// newProduct is sold by newShop
	newProduct.SoldBy = *clib.NewReference[clib.ShopStub](exportedShop.GetId())
	exportedProduct, err := clib.ExportProduct(newProduct)
	require.Equal(t, nil, err, "error occurred in ExportProduct invocation", err)
	// retrieve newProduct ID from newShop
	productsIds, err := exportedShop.GetProductsIds()
	require.Equal(t, nil, err, "error occurred in GetProductsIds invocation", err)
	products, err := exportedShop.GetProducts()
	require.Equal(t, nil, err, "error occurred in GetProducts invocation", err)

	// Assert
	require.Equal(t, 1, len(products), "expected number of products is 1, found %d", len(products))
	require.Equal(t, exportedProduct.GetId(), productsIds[0], "expected product id to be equal to the exported one, found %s", productsIds[0])
	require.Equal(t, exportedProduct.GetId(), products[0].GetId(), "expected product id to be equal to the exported one, found %s", products[0].GetId())
}

func TestReferenceNavigationListManyToMany(t *testing.T) {
	// Arrange
	newUserId := uuid.New().String()
	newUser := clib.UserStub{
		Email:     newUserId,
		FirstName: "TestReferenceNavigationList",
		LastName:  "TestManyToMany",
	}
	newShop := clib.ShopStub{
		Name: "ShopTestManyToManyRelationship",
	}

	// Act
	exportedShop, err := clib.ExportShop(newShop)
	require.Equal(t, nil, err, "error occurred in ExportShop invocation", err)

	exportedUser, err := clib.ExportUser(newUser)
	require.Equal(t, nil, err, "error occurred in ExportUser invocation", err)

	// TODO
	// Assert
}

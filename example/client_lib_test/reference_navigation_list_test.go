package client_lib_test

import (
	"testing"

	clib "github.com/Astenna/Nubes/example/client_lib"
	"github.com/Astenna/Nubes/example/faas/types"
	"github.com/Astenna/Nubes/lib"
	"github.com/stretchr/testify/require"

	"github.com/google/uuid"
)

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
	newProduct.SoldBy = exportedShop.AsReference()
	exportedProduct, err := clib.ExportProduct(newProduct)
	require.Equal(t, nil, err, "error occurred in ExportProduct invocation", err)
	// retrieve newProduct ID from newShop
	productsIds, err := exportedShop.Products.GetIds()
	require.Equal(t, nil, err, "error occurred in GetProductsIds invocation", err)
	products, err := exportedShop.Products.Get()
	require.Equal(t, nil, err, "error occurred in GetProducts invocation", err)

	// Assert
	require.Equal(t, 1, len(products), "expected number of products is 1, found %d", len(products))
	require.Equal(t, exportedProduct.GetId(), productsIds[0], "expected product id to be equal to the exported one, found %s", productsIds[0])
	require.Equal(t, exportedProduct.GetId(), products[0].GetId(), "expected product id to be equal to the exported one, found %s", products[0].GetId())
}

func TestReferenceNavigationListManyToManyByPartiotionKey(t *testing.T) {
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
	exportedUser, err := clib.ExportUser(newUser)
	require.Equal(t, nil, err, "error occurred in ExportUser invocation", err)

	exportedShop, err := clib.ExportShop(newShop)
	require.Equal(t, nil, err, "error occurred in ExportShop invocation", err)

	err = exportedShop.Owners.AddToManyToMany(newUserId)
	require.Equal(t, nil, err, "error occurred in AddOwners invocation", err)

	ownedShops, err := exportedUser.Shops.GetIds()
	require.Equal(t, nil, err, "error occurred in AddOwners invocation", err)

	// Assert
	require.Equal(t, 1, len(ownedShops), "expected number of ownedShops is 1, found %d", len(ownedShops))
	require.Equal(t, exportedShop.GetId(), ownedShops[0], "expected id of ownedShop to be equal to previously aded one, but found %s", ownedShops[0])
}

func TestReferenceNavigationListManyToManyByWithIndex(t *testing.T) {
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

	err = exportedUser.Shops.AddToManyToMany(exportedShop.GetId())
	require.Equal(t, nil, err, "error occurred in AddOwners invocation", err)

	shopOwners, err := exportedShop.Owners.GetIds()
	require.Equal(t, nil, err, "error occurred in AddOwners invocation", err)

	// Assert
	require.Equal(t, 1, len(shopOwners), "expected number of ownedShops is 1, found %d", len(shopOwners))
	require.Equal(t, newUserId, shopOwners[0], "expected id of ownedShop to be equal to previously aded one, but found %s", shopOwners[0])
}

func TestDeleteFromManyToManyRelationship(t *testing.T) {
	// Arrange
	newUserId := uuid.New().String()
	newUser := clib.UserStub{
		Email:     newUserId,
		FirstName: "TestReferenceNavigationList",
		LastName:  "TestManyToMany",
	}
	newShop := types.Shop{
		Name: "ShopTestManyToManyRelationship",
	}

	// Act
	exportedShop, err := lib.Export[types.Shop](newShop)
	require.Equal(t, nil, err, "error occurred in Export[types.Shop] invocation", err)

	exportedUser, err := lib.Export[types.User](newUser)
	require.Equal(t, nil, err, "error occurred in Export[types.User]", err)

	err = exportedUser.Shops.AddToManyToMany(exportedShop.Id)
	require.Equal(t, nil, err, "error occurred in exportedShop.Owners.Add", err)

	err = exportedUser.Shops.DeleteBatchFromManyToMany([]string{exportedShop.Id})
	require.Equal(t, nil, err, "error occurred in exportedUser.Shops.DeleteBatchFromManyToMany", err)
	owners, err := exportedShop.GetOwners()

	// Assert
	require.Equal(t, nil, err, "error occurred in  exportedShop.GetOwners()", err)
	require.Zero(t, len(owners))
}

package faas_lib_test

import (
	"testing"

	"github.com/Astenna/Nubes/example/faas/types"
	"github.com/Astenna/Nubes/lib"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestAddToManyToManyRelationshipRetrieveByPK(t *testing.T) {
	// Arrange
	newUserId := uuid.New().String()
	newUser := types.User{
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

	err = exportedShop.Owners.AddToManyToMany(newUserId)
	require.Equal(t, nil, err, "error occurred in exportedShop.Owners.Add", err)

	shopId, err := exportedUser.GetShops()
	require.Equal(t, nil, err, "error occurred in exportedUser.Shops.GetIds()", err)

	// Assert
	require.Equal(t, 1, len(shopId), "expected number of shop ids is 1, found", len(shopId))
	require.Equal(t, exportedShop.Id, shopId[0], "expected returned shopId to be equal previously exported one, but found", shopId[0])
}

func TestAddToManyToManyRelationshipRetrieveByIndex(t *testing.T) {
	// Arrange
	newUserId := uuid.New().String()
	newUser := types.User{
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

	ownersId, err := exportedShop.GetOwners()
	require.Equal(t, nil, err, "error occurred in exportedUser.Shops.GetIds()", err)

	// Assert
	require.Equal(t, 1, len(ownersId), "expected number of shop ids is 1, found", len(ownersId))
	require.Equal(t, newUserId, ownersId[0], "expected returned ownerId to be equal previously exported one, but found", ownersId[0])
}

func TestGetLoaded(t *testing.T) {
	// Arrange
	newProduct1 := types.Product{
		Name:              "GetLoadedTest1",
		QuantityAvailable: 525,
		Price:             10.99,
	}
	newProduct2 := types.Product{
		Name:              "GetLoadedTest2",
		QuantityAvailable: 525,
		Price:             10.99,
	}
	newShop := types.Shop{
		Name: "ShopTestGetLoaded",
	}

	// Act
	exportedShop, err := lib.Export[types.Shop](newShop)
	require.Equal(t, nil, err, "error occurred in Export[types.Shop] invocation", err)

	newProduct1.SoldBy = lib.Reference[types.Shop](exportedShop.Id)
	exportedProduct1, err := lib.Export[types.Product](newProduct1)
	require.Equal(t, nil, err, "error occurred in Export[types.Product]", err)

	newProduct2.SoldBy = lib.Reference[types.Shop](exportedShop.Id)
	exportedProduct2, err := lib.Export[types.Product](newProduct2)
	require.Equal(t, nil, err, "error occurred in Export[types.Product]", err)

	productsList, err := exportedShop.Products.Get()
	require.Equal(t, nil, err, "error occurred in exportedShop.Products.GetLoaded()", err)

	// Assert
	require.Equal(t, 2, len(productsList), "expected number of shop's products is 2, found", len(productsList))

	productsList0Quantity, err := productsList[0].GetQuantityAvailable()
	require.Equal(t, nil, err, "error occurred in  productsList[0].GetQuantityAvailable()", err)
	exportedProduct1Quantity, err := exportedProduct1.GetQuantityAvailable()
	require.Equal(t, nil, err, "error occurred in exportedProduct1.GetQuantityAvailable()", err)

	productsList1Quantity, err := productsList[1].GetQuantityAvailable()
	require.Equal(t, nil, err, "error occurred in  productsList[1].GetQuantityAvailable()", err)
	exportedProduct2Quantity, err := exportedProduct2.GetQuantityAvailable()
	require.Equal(t, nil, err, "error occurred in exportedProduct2.GetQuantityAvailable()", err)

	if productsList[0].Id == exportedProduct1.Id {
		require.Equal(t, productsList0Quantity, exportedProduct1Quantity)
		// require.Equal(t, productsList[0].GetName(), exportedProduct1.GetName())
		require.Equal(t, productsList1Quantity, exportedProduct2Quantity)
		// require.Equal(t, productsList[1].GetName(), exportedProduct2.GetName())
	} else {
		require.Equal(t, productsList1Quantity, exportedProduct1Quantity)
		require.Equal(t, productsList0Quantity, exportedProduct2Quantity)
	}
}

func TestGetIds(t *testing.T) {
	// Arrange
	newProduct1 := types.Product{
		Name:              "GetLoadedTest1",
		QuantityAvailable: 525,
		Price:             10.99,
	}
	newProduct2 := types.Product{
		Name:              "GetLoadedTest2",
		QuantityAvailable: 525,
		Price:             10.99,
	}
	newShop := types.Shop{
		Name: "ShopTestGetIds",
	}

	// Act
	exportedShop, err := lib.Export[types.Shop](newShop)
	require.Equal(t, nil, err, "error occurred in Export[types.Shop] invocation", err)

	newProduct1.SoldBy = lib.Reference[types.Shop](exportedShop.Id)
	exportedProduct1, err := lib.Export[types.Product](newProduct1)
	require.Equal(t, nil, err, "error occurred in Export[types.Product]", err)

	newProduct2.SoldBy = lib.Reference[types.Shop](exportedShop.Id)
	exportedProduct2, err := lib.Export[types.Product](newProduct2)
	require.Equal(t, nil, err, "error occurred in Export[types.Product]", err)

	productIds, err := exportedShop.Products.GetIds()
	require.Equal(t, nil, err, "error occurred in exportedShop.Products.GetLoaded()", err)

	// Assert
	require.Equal(t, 2, len(productIds), "expected number of shop's products is 2, found", len(productIds))
	require.ElementsMatch(t, []string{exportedProduct1.Id, exportedProduct2.Id}, productIds, "expected product Ids to match, found different values")
}

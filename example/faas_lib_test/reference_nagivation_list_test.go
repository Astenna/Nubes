package faas_lib_test

import (
	"testing"

	"github.com/Astenna/Nubes/faas/types"
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

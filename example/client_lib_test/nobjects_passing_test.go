package client_lib_test

import (
	"testing"

	clib "github.com/Astenna/Nubes/example/client_lib"
	"github.com/stretchr/testify/require"

	"github.com/google/uuid"
)

func TestNobjectAsReturnParam(t *testing.T) {

	// Arrange
	newId1 := uuid.New().String()
	newId2 := uuid.New().String()
	coordinates1 := clib.Coordinates{Longitude: 10, Latitude: 10}
	coordinates2 := clib.Coordinates{Longitude: 15, Latitude: 15}
	newShop := clib.ShopStub{
		Name: "TestReferenceNavigationListShop",
	}
	newUser1 := clib.UserStub{
		Email:              newId1,
		FirstName:          "Kinga1",
		LastName:           "Marek1",
		Password:           "password1",
		AddressCoordinates: coordinates1,
	}
	newUser2 := clib.UserStub{
		Email:              newId2,
		FirstName:          "Kinga2",
		LastName:           "Marek2",
		Password:           "password2",
		AddressCoordinates: coordinates2,
	}

	// Act
	exportedShop, err := clib.ExportShop(newShop)
	require.Equal(t, nil, err, "error occurred in ExportShop invocation", err)

	_, err = clib.ExportUser(newUser1)
	require.Equal(t, nil, err, "error occurred in ExportUser invocation", err)

	_, err = clib.ExportUser(newUser2)
	require.Equal(t, nil, err, "error occurred in ExportUser invocation", err)

	err = exportedShop.Owners.AddToManyToMany(newId1)
	require.Equal(t, nil, err, "error occurred in AddToManyToMany invocation", err)

	err = exportedShop.Owners.AddToManyToMany(newId2)
	require.Equal(t, nil, err, "error occurred in AddToManyToMany invocation", err)

	nearestRef, err1 := exportedShop.GetNearestOwnerReference(clib.Coordinates{Longitude: 11, Latitude: 11})
	nearestStub, err2 := exportedShop.GetNearestOwnerCopy(clib.Coordinates{Longitude: 11, Latitude: 11})

	// Assert
	require.Equal(t, nil, err1, "error occurred in GetNearestOwnerReference invocation", err)
	require.Equal(t, nil, err2, "error occurred in GetNearestOwnerCopy invocation", err)
	require.Equal(t, nearestRef.Id(), nearestStub.Email)
}

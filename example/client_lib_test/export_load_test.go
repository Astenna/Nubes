package client_lib_test

import (
	"testing"

	clib "github.com/Astenna/Nubes/example/client_lib"
	"github.com/stretchr/testify/require"

	"github.com/google/uuid"
)

func TestLoadAndExport(t *testing.T) {
	// Arrange
	existingOrderId := "d192eeda-e709-4415-bbe2-cb91a4968962"
	newId := uuid.New().String()
	newUser := clib.UserStub{
		Email:     newId,
		FirstName: "Kinga",
		LastName:  "Marek",
		Password:  "password",
		Orders:    append(clib.OrderReferenceList(1), existingOrderId),
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

package faas_lib_test

import (
	"testing"

	"github.com/Astenna/Nubes/example/faas/types"
	"github.com/Astenna/Nubes/lib"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestDelete(t *testing.T) {
	// Arrange
	product := types.Product{
		Name:              "Product1",
		QuantityAvailable: 100,
		Price:             88.88,
	}
	// Act
	exportedProduct, err := lib.Export[types.Product](product)
	require.Equal(t, err, nil)
	err = lib.Delete[types.Product](exportedProduct.Id)
	// Assert
	require.Equal(t, err, nil, "error should be null")
}

func TestCustomDelete(t *testing.T) {
	// Arrange
	newEmail := uuid.New().String()
	password := "password"
	newUser := types.User{Email: newEmail, Password: password, FirstName: "Kinga", LastName: "Marek"}
	// Act
	_, err := lib.Export[types.User](newUser)
	require.Equal(t, err, nil)
	err = types.DeleteUser(types.DeleteParam{Password: password, Email: newEmail})
	// Assert
	require.Equal(t, err, nil)
}

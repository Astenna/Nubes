package client_lib_test

import (
	"testing"

	clib "github.com/Astenna/Nubes/example/client_lib"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestDelete(t *testing.T) {
	// Arrange
	product := clib.ProductStub{
		Name:              "Product1",
		QuantityAvailable: 100,
		Price:             88.88,
	}
	// Act
	exportedProduct, err := clib.ExportProduct(product)
	require.Equal(t, err, nil)
	err = clib.DeleteProduct(exportedProduct.GetId())
	// Assert
	require.Equal(t, err, nil, "error should be null")
}

func TestCustomDelete(t *testing.T) {
	// Arrange
	newEmail := uuid.New().String()
	password := "password"
	newUser := clib.UserStub{Email: newEmail, Password: password, FirstName: "Kinga", LastName: "Marek"}
	// Act
	_, err := clib.ExportUser(newUser)
	require.Equal(t, err, nil)
	err = clib.DeleteUser(clib.DeleteParam{Password: password, Email: newEmail})
	// Assert
	require.Equal(t, err, nil)
}

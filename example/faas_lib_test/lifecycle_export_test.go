package faas_lib_test

import (
	"testing"

	"github.com/Astenna/Nubes/example/faas/types"
	"github.com/Astenna/Nubes/lib"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// EXPORT
// - fail if alredy exists
// - success if not existing
// - changes made through get/set are saved (retrieve from another instance)
// - changes made through state changing method exist

var existingUserId string

func init() {
	existingUserId = "crazy_integration_test@email.com"
	lib.Export[types.User](types.User{Email: existingUserId})
}

func TestExportFailsIfIdAlreadyExists(t *testing.T) {
	// Arrange
	// Act
	result, err := lib.Export[types.User](types.User{Email: existingUserId})
	// Assert
	if err == nil || result != nil {
		t.Error("export should file to recreate existing user id")
	}
}

func TestExportCreatesNewInstance(t *testing.T) {
	// Arrange
	notExistingId := uuid.New().String()
	newUser := types.User{Email: notExistingId, Password: "password", FirstName: "Kinga", LastName: "Marek"}
	// Act
	_, err := lib.Export[types.User](newUser)
	// Assert
	require.Equal(t, err, nil, "error should be null")
}

func TestCustomExport(t *testing.T) {
	// Arrange
	addr := "Address for TestCustomExportDefinition"
	// Act
	id, err := types.ExportShipping(addr)
	require.Equal(t, err, nil, "error should be null")
	_, err = lib.Load[types.Shipping](id)
	// Assert
	require.Equal(t, err, nil, "error should be null")
}

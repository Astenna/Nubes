package faas_lib_test

import (
	"testing"

	"github.com/Astenna/Nubes/faas/types"
	"github.com/Astenna/Nubes/lib"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// EXPORT
// - fail if alredy exists
// - success if not existing
// - changes made through get/set are saved (retrieve from another instance)
// - changes made through state changing method exist

var existingUserId string

func init() {
	existingUserId = "crazy_integration_test@emai.com"
	_, err := lib.Load[types.User](existingUserId)
	if err != nil {
		panic("initialization error")
	}
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
	result, err := lib.Export[types.User](newUser)
	// Assert
	assert.Equal(t, err, nil, "error should be null")
	assert.Equal(t, newUser.Email, result.Email, "IDs should match")
}

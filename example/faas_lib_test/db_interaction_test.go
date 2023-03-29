package faas_lib_test

import (
	"testing"

	"github.com/Astenna/Nubes/example/faas/types"
	"github.com/Astenna/Nubes/lib"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestIsInstanceAlreadyCreatedReturnsFalseIfItemNotFound(t *testing.T) {
	// Arrange
	input := lib.IsInstanceAlreadyCreatedParam{
		Id:       uuid.New().String(),
		TypeName: (*(new(types.User))).GetTypeName(),
	}
	// Act
	result, err := lib.IsInstanceAlreadyCreated(input)
	// Assert
	require.Equal(t, err, nil, "error occurred in IsInstanceAlreadyCreated invocation", err)
	require.Equal(t, result, false, "expected false, returned true")
}

func TestIsInstanceAlreadyCreatedReturnsTrueIfItemExists(t *testing.T) {
	// Arrange
	input := lib.IsInstanceAlreadyCreatedParam{
		Id:       existingUserId,
		TypeName: (*(new(types.User))).GetTypeName(),
	}
	// Act
	result, err := lib.IsInstanceAlreadyCreated(input)
	// Assert
	require.Equal(t, err, nil, "error occurred in IsInstanceAlreadyCreated invocation", err)
	require.Equal(t, result, true, "expected true, returned false")
}

func TestAreInstancesAlreadyCreatedReturnsNotFoundErrorIfInstanceDoesNotExist(t *testing.T) {
	// Arrange
	input := lib.LoadBatchParam{
		Ids:      []string{uuid.NewString(), uuid.NewString()},
		TypeName: (*(new(types.User))).GetTypeName(),
	}
	// Act
	err := lib.AreInstancesAlreadyCreated(input)
	// Assert
	require.IsType(t, err, lib.NotFoundError{})
	require.Equal(t, err.(lib.NotFoundError).TypeName, input.TypeName)
	require.ElementsMatch(t, err.(lib.NotFoundError).Ids, input.Ids)
}

func TestAreInstancesAlreadyCreatedReturnsIdsOfNotExistingInstances(t *testing.T) {
	// Arrange
	existingId := uuid.New().String()
	newUser := types.User{Email: existingId, Password: "password", FirstName: "Kinga", LastName: "Marek"}
	input := lib.LoadBatchParam{
		Ids:      []string{uuid.NewString(), uuid.NewString(), existingId},
		TypeName: (*(new(types.User))).GetTypeName(),
	}

	// Act
	_, err := lib.Export[types.User](newUser)
	require.Equal(t, err, nil, "error should be null")
	err = lib.AreInstancesAlreadyCreated(input)

	// Assert
	require.IsType(t, lib.NotFoundError{}, err)
	require.Equal(t, input.TypeName, err.(lib.NotFoundError).TypeName)
	require.ElementsMatch(t, err.(lib.NotFoundError).Ids, input.Ids[0:2])
	require.NotContains(t, err.(lib.NotFoundError).Ids, []string{existingId})
}

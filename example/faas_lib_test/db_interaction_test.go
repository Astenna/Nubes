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

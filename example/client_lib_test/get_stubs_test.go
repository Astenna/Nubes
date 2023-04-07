package client_lib_test

import (
	"testing"

	clib "github.com/Astenna/Nubes/example/client_lib"
	"github.com/stretchr/testify/require"

	"github.com/google/uuid"
)

func TestReferenceListGetStubs(t *testing.T) {
	// Arrange
	newId1 := uuid.New().String()
	newId2 := uuid.New().String()
	newUser1 := clib.UserStub{
		Email:     newId1,
		FirstName: "Kinga1",
		LastName:  "Marek1",
		Password:  "password1",
	}
	newUser2 := clib.UserStub{
		Email:     newId2,
		FirstName: "Kinga2",
		LastName:  "Marek2",
		Password:  "password2",
	}
	ids := []string{newId1, newId2}

	// Act
	_, err := clib.ExportUser(newUser1)
	require.Equal(t, nil, err, "error occurred in ExportUser invocation", err)

	_, err = clib.ExportUser(newUser2)
	require.Equal(t, nil, err, "error occurred in ExportUser invocation", err)

	stubs, err := clib.GetStubs[clib.UserStub](ids)

	// Assert
	require.Equal(t, nil, err, "error occurred in clib.GetStubs[clib.UserStub](ids) invocation", err)
	if stubs[0].Email == newId1 {
		require.EqualValues(t, newUser1, stubs[0])
		require.EqualValues(t, newUser2, stubs[1])
	} else {
		require.EqualValues(t, newUser1, stubs[1])
		require.EqualValues(t, newUser2, stubs[0])
	}
}

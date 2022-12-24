package faas

import (
	"errors"

	"github.com/Astenna/Nubes/faas/types"
	lib "github.com/Astenna/Nubes/lib"
)

func CreateUser(user types.User) (string, error) {
	if user.LastName != "" && user.FirstName != "" {
		newId, err := lib.Insert(&user)

		if err != nil {
			return "", errors.New("failed to create user")
		}
		return newId, nil
	}

	return "", errors.New("the fields FirstName and LastName can not be empty")
}

func DeleteUser(id string) error {
	err := lib.Delete[types.User](id)
	return err
}

func GetUser(id string) (*types.User, error) {
	return lib.Get[types.User](id)
}

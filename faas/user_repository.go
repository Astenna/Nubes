package faas

import (
	"errors"

	"github.com/Astenna/Nubes/faas/types"
	lib "github.com/Astenna/Nubes/lib"
)

func CreateUser(user types.User) error {
	if user.LastName != "" && user.FirstName != "" {
		err := lib.Insert(&user)

		if err != nil {
			return errors.New("failed to create user")
		}
		return nil
	}

	return errors.New("the fields FirstName and LastName can not be empty")
}

func DeleteUser(id string) error {
	err := lib.Delete[types.User](id)
	return err
}

func GetUser(id string) (*types.User, error) {
	return lib.Get[types.User](id)
}

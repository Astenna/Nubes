package faas

import (
	"errors"
	"github.com/Astenna/Thesis_PoC/faas/types"
	lib "github.com/Astenna/Thesis_PoC/faas_lib"
)

func CreateUser(user types.User) error {
	if user.LastName != "" && user.FirstName != "" {
		err := lib.Create(&user)

		if err != nil {
			return errors.New("failed to create user")
		}
		return nil
	}

	return errors.New("the fields FirstName and LastName can not be empty")
}

func DeleteUser(id int) error {
	err := lib.Delete[types.User](id)
	return err
}

func GetUser(id int) (*types.User, error) {
	return lib.Get[types.User](id)
}

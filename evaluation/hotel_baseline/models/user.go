package models

import (
	"errors"

	"github.com/Astenna/Nubes/evaluation/hotel_baseline/db"
)

type User struct {
	Email     string
	FirstName string
	LastName  string
	Password  string
}

func Login(email, password string) error {
	user, err := db.GetSingleItemByPartitonKey[User](db.UserTable, "Email", email)

	if err != nil {
		return err
	}

	if user.Password == password {
		return nil
	}

	return errors.New("incorrect password")
}

func DeleteUser(email, password string) error {
	user, err := db.GetSingleItemByPartitonKey[User](db.UserTable, "Email", email)
	if err != nil {
		return err
	}

	if user.Password == password {
		return db.DeleteUser(email)
	}

	return errors.New("invalid password")
}

func RegisterUser(user User) error {
	return db.Insert(user, db.UserTable)
}

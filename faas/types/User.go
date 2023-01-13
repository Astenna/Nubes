package types

import (
	"github.com/Astenna/Nubes/lib"
)

type User struct {
	FirstName	string
	LastName	string
	Email		string	`dynamodbav:"Id" nubes:"readonly"`
	Password	string	`nubes:"readonly"`
	Address		string
}

func NewUser(user User) (User, error) {
	out, _libError := lib.Insert(user)
	if _libError != nil {
		return *new(User), _libError
	}
	user.Email = out
	return user, nil
}

func ReNewUser(id string) User {
	user := new(User)
	user.Email = id
	return *user
}

func DeleteUser(id string) {

}

func (User) GetTypeName() string {
	return "User"
}

func (u User) GetId() string {
	return u.Email
}

func (u User) VerifyPassword(password string) (bool, error) {
	tempReceiverName, _libError := lib.Get[User](u.Email)
	u = *tempReceiverName
	if _libError != nil {
		return *new(bool), _libError
	}
	if u.Password == password {
		return true, nil
	}
	_libError = lib.Upsert(u, u.Email)
	if _libError != nil {
		return *new(bool), _libError
	}
	return false, nil
}

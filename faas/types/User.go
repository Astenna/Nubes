package types

import "github.com/Astenna/Nubes/lib"

type User struct {
	FirstName	string
	LastName	string
	Email		string	`dynamodbav:"Id" nubes:"readonly"`
	Password	string	`nubes:"readonly"`
	Address		string
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

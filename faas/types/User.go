package types

import (
	"github.com/Astenna/Nubes/lib"
)

type User struct {
	FirstName string
	LastName  string
	Email     string `dynamodbav:"Id" nubes:"id,readonly"`
	Password  string `nubes:"readonly"`
	Address   string
	Shops     lib.ReferenceNavigationList[Shop] `nubes:"hasMany-Owners"`
}

func DeleteUser(id string) error {
	_libError := lib.Delete[User](id)
	if _libError != nil {
		return _libError
	}
	return nil
}

func (User) GetTypeName() string {
	return "User"
}

func (User) New() error {
	return nil
}

func (u User) GetId() string {
	return u.Email
}

func (u *User) Init() error {
	u.Shops = *lib.NewReferenceNavigationList[Shop](u.Email, u.GetTypeName(), "", true)
	return nil
}

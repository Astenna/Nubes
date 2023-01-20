package types

import (
	"github.com/Astenna/Nubes/lib"
)

type User struct {
	FirstName     string
	LastName      string
	Email         string `dynamodbav:"Id" nubes:"id,readonly"`
	Password      string `nubes:"readonly"`
	Address       string
	Shops         lib.ReferenceNavigationList[Shop] `nubes:"hasMany-Owners" dynamodbav:"-"`
	isInitialized bool
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

func (u User) GetId() string {
	return u.Email
}

func (receiver *User) Init() {
	receiver.isInitialized = true
	receiver.Shops = *lib.NewReferenceNavigationList[Shop](receiver.Email, receiver.GetTypeName(), "", true)
}

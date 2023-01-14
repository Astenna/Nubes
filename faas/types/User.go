package types

import (
	"github.com/Astenna/Nubes/lib"
)

type User struct {
	FirstName string
	LastName  string
	Email     string `dynamodbav:"Id" nubes:"readonly"`
	Password  string `nubes:"readonly"`
	Address   string
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

func (u User) GetAddress() (string, error) {
	fieldValue, _libError := lib.GetField(lib.HandlerParameters{Id: u.Email, Parameter: lib.GetFieldParam{TypeName: "User", FieldName: "Address"}})
	if _libError != nil {
		return *new(string), _libError
	}
	u.Address = fieldValue.(string)
	return u.Address, nil
}

func (u *User) SetAddress(adr string) error {
	u.Address = adr
	_libError := lib.SetField(lib.HandlerParameters{Id: u.Email, Parameter: lib.SetFieldParam{TypeName: "User", FieldName: "Address", Value: u.Address}})
	if _libError != nil {
		return _libError
	}
	return nil
}

func (u User) VerifyPassword(password string) (bool, error) {
	tempReceiverName, _libError := lib.Get[User](u.Email)
	if _libError != nil {
		return *new(bool), _libError
	}
	u = *tempReceiverName

	if u.Password == password {
		return true, nil
	}
	_libError = lib.Upsert(u, u.Email)
	if _libError != nil {
		return *new(bool), _libError
	}
	return false, nil
}

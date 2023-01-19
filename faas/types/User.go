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
	Shops     lib.ReferenceList[Shop]
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
	return "sjsjs"
}

// func (u User) GetAddress() (string, error) {
// 	fieldValue, _libError := lib.GetField(u.Email, lib.GetFieldParam{TypeName: "User", FieldName: "Address"})
// 	if _libError != nil {
// 		return *new(string), _libError
// 	}
// 	u.Address = fieldValue.(string)
// 	return u.Address, nil
// }

// func (u *User) SetAddress(adr string) error {
// 	u.Address = adr
// 	_libError := lib.SetField(u.Email, lib.SetFieldParam{TypeName: "User", FieldName: "Address", Value: u.Address})
// 	if _libError != nil {
// 		return _libError
// 	}
// 	return nil
// }

// func (u User) VerifyPassword(password string) (bool, error) {
// 	tempReceiverName, _libError := lib.GetObjectState[User](u.Email)
// 	if _libError != nil {
// 		return *new(bool), _libError
// 	}
// 	u = *tempReceiverName
// 	if u.Password == password {
// 		return true, nil
// 	}
// 	_libError = lib.Upsert(u, u.Email)
// 	if _libError != nil {
// 		return *new(bool), _libError
// 	}
// 	return false, nil
// }

package types

import (
	"fmt"

	"github.com/Astenna/Nubes/lib"
)

type Coordinates struct {
	Longitude float64
	Latitude  float64
}

type User struct {
	FirstName          string
	LastName           string
	Email              string `nubes:"id,readonly" dynamodbav:"Id"`
	Password           string `nubes:"readonly"`
	AddressText        string
	AddressCoordinates Coordinates
	Shops              lib.ReferenceNavigationList[Shop] `nubes:"hasMany-Owners" dynamodbav:"-"`
	Orders             lib.ReferenceList[Order]
	isInitialized      bool
	invocationDepth    int
}

type DeleteParam struct {
	Email    string
	Password string
}

// DeleteUser is an example of custom delete implementation
// that uses input parameter type different than in the default delete.
// Note that, the invocation of lib.Delete must be added inside the function.
func DeleteUser(param DeleteParam) error {
	userToBeDeleted, err := lib.Load[User](param.Email)
	if err != nil {
		return err
	}

	passwordOk, err := userToBeDeleted.VerifyPassword(param.Password)
	if err != nil {
		return err
	}
	if passwordOk {
		return lib.Delete[User](param.Email)
	}

	return fmt.Errorf("invalid password")
}

func (User) GetTypeName() string {
	return "User"
}

func (u User) GetId() string {
	return u.Email
}

func (u *User) SetLastName(lastName string) error {
	u.LastName = lastName
	if u.isInitialized {
		_libError := lib.SetField(lib.SetFieldParam{Id: u.Email, TypeName: "User", FieldName: "LastName", Value: u.LastName})
		if _libError != nil {
			return _libError
		}
	}
	return nil
}

func (u *User) GetLastName() (string, error) {
	if u.isInitialized {
		fieldValue, _libError := lib.GetFieldOfType[string](lib.GetStateParam{Id: u.Email, TypeName: "User", FieldName: "LastName"})
		if _libError != nil {
			return *new(string), _libError
		}
		u.LastName = fieldValue
	}
	return u.LastName, nil
}

func (u User) GetShops() ([]string, error) {
	if !u.isInitialized {
		return nil, fmt.Errorf(`fields of type ReferenceNavigationList can be used only after instance initialization. 
			Use lib.Load or lib.Export from the Nubes library to create initialized instances`)
	}
	return u.Shops.GetIds()
}

func (u User) VerifyPassword(password string) (bool, error) {
	u.invocationDepth++
	if u.isInitialized && u.invocationDepth == 1 {
		_libError := lib.GetStub(u.Email, &u)
		if _libError != nil {
			u.invocationDepth--
			return *new(bool), _libError
		}
	}

	if u.Password == password {
		_libUpsertError := u.saveChangesIfInitialized()
		u.invocationDepth--
		return true, _libUpsertError
	}
	_libUpsertError := u.saveChangesIfInitialized()
	u.invocationDepth--
	return false, _libUpsertError
}

func (receiver *User) Init() {
	receiver.isInitialized = true
	receiver.Shops = *lib.NewReferenceNavigationList[Shop](receiver.Email, receiver.GetTypeName(), "", true)
}
func (receiver *User) saveChangesIfInitialized() error {
	if receiver.isInitialized && receiver.invocationDepth == 1 {
		_libError := lib.Upsert(receiver, receiver.Email)
		if _libError != nil {
			return _libError
		}
	}
	return nil
}

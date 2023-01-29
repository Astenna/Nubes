package types

import "github.com/Astenna/Nubes/lib"

type Account struct {
	Nickname      string
	Email         string `nubes:"id,readonly" dynamodbav:"Id"`
	Password      string `nubes:"readonly"`
	isInitialized bool
}

func (Account) GetTypeName() string {
	return "Account"
}

func (u Account) GetId() string {
	return u.Email
}

func (u Account) VerifyPassword(password string) (bool, error) {
	if u.isInitialized {
		tempReceiverName, _libError := lib.GetObjectState[Account](u.Email)
		if _libError != nil {
			return *new(bool), _libError
		}
		u = *tempReceiverName
		u.Init()
	}

	if u.Password == password {
		return true, nil
	}
	if u.isInitialized {
		_libError := lib.Upsert(u, u.Email)
		if _libError != nil {
			return *new(bool), _libError
		}
	}
	return false, nil
}

func (receiver *Account) Init() {
	receiver.isInitialized = true
}

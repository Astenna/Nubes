package types

import "github.com/Astenna/Nubes/lib"

type Account struct {
	Nickname        string
	Email           string `nubes:"id,readonly" dynamodbav:"Id"`
	Password        string `nubes:"readonly"`
	isInitialized   bool
	invocationDepth int
}

func (Account) GetTypeName() string {
	return "Account"
}

func (u Account) GetId() string {
	return u.Email
}

func (u Account) VerifyPassword(password string) (bool, error) {
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
	if u.isInitialized {
		_libError := lib.Upsert(u, u.Email)
		if _libError != nil {
			u.invocationDepth--
			return *new(bool), _libError
		}
	}
	_libUpsertError := u.saveChangesIfInitialized()
	u.invocationDepth--
	return false, _libUpsertError
}
func (receiver *Account) Init() {
	receiver.isInitialized = true
}
func (receiver *Account) saveChangesIfInitialized() error {
	if receiver.isInitialized && receiver.invocationDepth == 1 {
		_libError := lib.Upsert(receiver, receiver.Email)
		if _libError != nil {
			return _libError
		}
	}
	return nil
}

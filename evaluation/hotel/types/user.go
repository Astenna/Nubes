package types

import "github.com/Astenna/Nubes/lib"

type User struct {
	FirstName       string
	LastName        string
	Email           string `nubes:"id,readonly" dynamodbav:"Id"`
	Password        string
	Reservations    lib.ReferenceNavigationList[Reservation] `nubes:"hasOne-User" dynamodbav:"-"`
	isInitialized   bool
	invocationDepth int
}

func (o User) GetTypeName() string {
	return "User"
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
func (receiver User) GetId() string {
	return receiver.Email
}
func (receiver *User) Init() {
	receiver.isInitialized = true
	receiver.Reservations = *lib.NewReferenceNavigationList[Reservation](receiver.Email, receiver.GetTypeName(), "User", false)
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

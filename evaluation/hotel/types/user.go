package types

import "github.com/Astenna/Nubes/lib"

type User struct {
	FirstName    string
	LastName     string
	Email        string `nubes:"id,readonly" dynamodbav:"Id"`
	Password     string
	Reservations lib.ReferenceNavigationList[Reservation] `nubes:"hasOne-User"`
}

func (o User) GetTypeName() string {
	return "User"
}

func (u User) VerifyPassword(password string) (bool, error) {

	if u.Password == password {
		return true, nil
	}
	return false, nil
}

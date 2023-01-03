package types

type User struct {
	FirstName string
	LastName  string
	Email     string `dynamodbav:"Id" nubes:"readonly"`
	Password  string `nubes:"readonly"`
	Address   string
}

func (User) GetTypeName() string {
	return "User"
}

func (u User) GetId() string {
	return u.Email
}

func (u User) VerifyPassword(password string) (bool, error) {
	if u.Password == password {
		return true, nil
	}
	return false, nil
}

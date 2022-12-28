package types

type User struct {
	FirstName string
	LastName  string
	Email     string `dynamodbav:"Id" nubes:"readonly"`
	Password  string `nubes:"readonly"`
}

func (User) GetTypeName() string {
	return "User"
}

func (u User) GetId() string {
	return u.Email
}

func (u User) StateChange(test string) (User, error) {
	return *new(User), nil
}

package types

type User struct {
	Id        string
	FirstName string
	LastName  string
	Email     string
}

func (User) GetTypeName() string {
	return "User"
}

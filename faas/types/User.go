package types

type User struct {
	Id        int
	FirstName string
	LastName  string
	Email     string
}

func (User) GetTypeName() string {
	return "User"
}

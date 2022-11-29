package types

type User struct {
	Id        int
	FirstName string
	LastName  string
}

func (User) GetTableName() string {
	return "User"
}

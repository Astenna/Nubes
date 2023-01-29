package models

type Account struct {
	Nickname string
	Email    string `dynamodbav:"Id"`
	Password string
}

type LoginParams struct {
	Email    string
	Password string
}

func (a Account) GetId() string {
	return a.Email
}

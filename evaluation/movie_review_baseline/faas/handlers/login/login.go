package main

import (
	"strings"

	"github.com/Astenna/Nubes/movie_review_baseline/faas/db"
	"github.com/Astenna/Nubes/movie_review_baseline/faas/models"
	"github.com/aws/aws-lambda-go/lambda"
)

func LoginHandler(input models.LoginParams) (bool, error) {
	account, err := db.GetById[models.Account](input.Email, "Account")
	if err != nil {
		return false, err
	}

	// never store passwords in plain text!
	// this is just for demonstration purposes
	return strings.Compare(account.Password, input.Password) == 0, nil
}

func main() {
	lambda.Start(LoginHandler)
}

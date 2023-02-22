package main

import (
	"fmt"

	"github.com/Astenna/Nubes/evaluation/movie_review_baseline/faas/db"
	"github.com/Astenna/Nubes/evaluation/movie_review_baseline/faas/models"
	"github.com/aws/aws-lambda-go/lambda"
)

func RegisterHandler(input models.Account) (string, error) {
	if input.Email == "" && input.Password == "" && input.Nickname == "" {
		return "", fmt.Errorf("missing email, password or nickane")
	}
	return db.Insert(input, "Account")
}

func main() {
	lambda.Start(RegisterHandler)
}

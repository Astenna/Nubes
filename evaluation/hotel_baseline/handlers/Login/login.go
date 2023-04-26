package main

import (
	"github.com/Astenna/Nubes/evaluation/hotel_baseline/models"
	"github.com/aws/aws-lambda-go/lambda"
)

type LoginParams struct {
	Email    string
	Password string
}

func Login(login LoginParams) error {
	return models.Login(login.Email, login.Password)
}

func main() {
	lambda.Start(Login)
}

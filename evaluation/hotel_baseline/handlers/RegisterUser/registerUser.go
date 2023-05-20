package main

import (
	"github.com/Astenna/Nubes/evaluation/hotel_baseline/models"
	"github.com/aws/aws-lambda-go/lambda"
)

func RegisterUser(user models.User) error {
	return models.RegisterUser(user)
}

func main() {
	lambda.Start(RegisterUser)
}

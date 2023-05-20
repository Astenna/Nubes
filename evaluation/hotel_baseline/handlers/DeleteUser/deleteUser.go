package main

import (
	"github.com/Astenna/Nubes/evaluation/hotel_baseline/models"
	"github.com/aws/aws-lambda-go/lambda"
)

type DeleteParam struct {
	Email    string
	Password string
}

func DeleteUser(param DeleteParam) error {
	return models.DeleteUser(param.Email, param.Password)
}

func main() {
	lambda.Start(DeleteUser)
}

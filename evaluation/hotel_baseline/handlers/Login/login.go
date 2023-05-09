package main

import (
	"log"

	"github.com/Astenna/Nubes/evaluation/hotel_baseline/models"
	"github.com/aws/aws-lambda-go/lambda"
)

type LoginParams struct {
	Email    string
	Password string
}

func Login(login LoginParams) error {
	log.Printf("INPUT: %+v", login)
	out := models.Login(login.Email, login.Password)
	if out != nil {
		log.Printf("OUT: %+v", out)

	} else {
		log.Printf("OUT: SUCCESS")
	}
	return out
}

func main() {
	lambda.Start(Login)
}

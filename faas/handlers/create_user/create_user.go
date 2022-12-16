package main

import (
	"github.com/Astenna/Nubes/faas"
	"github.com/Astenna/Nubes/faas/types"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler1(user types.User) error {
	if err := faas.CreateUser(user); err != nil {
		return err
	}
	return nil
}

func main() {
	lambda.Start(Handler1)
}

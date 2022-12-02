package main

import (
	"github.com/Astenna/Thesis_PoC/faas"
	"github.com/Astenna/Thesis_PoC/faas/types"
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

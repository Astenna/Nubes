package main

import (
	"github.com/Astenna/Thesis_PoC/faas"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler2(id string) error {
	if err := faas.DeleteUser(id); err != nil {
		return err
	}
	return nil
}

func main() {
	lambda.Start(Handler2)
}

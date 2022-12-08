package main

import (
	"github.com/Astenna/Thesis_PoC/faas"
	"github.com/Astenna/Thesis_PoC/faas/types"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler3(id string) (types.User, error) {
	user, err := faas.GetUser(id)
	return *user, err
}

func main() {
	lambda.Start(Handler3)
}

package main

import (
	"github.com/aws/aws-lambda-go/lambda"
)

func CreateShop(id int, name string) error {
	// retrieve by id from DB
	// change state on the object from DB
	// save to DB
	return nil
}

func main() {
	lambda.Start(CreateShop)
}

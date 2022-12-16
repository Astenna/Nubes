package main

import (
	"github.com/Astenna/Nubes/faas"
	"github.com/Astenna/Nubes/faas/types"
	"github.com/aws/aws-lambda-go/lambda"
)

func CreateShop(shop types.Shop) error {
	if err := faas.CreateShop(shop); err != nil {
		return err
	}
	return nil
}

func main() {
	lambda.Start(CreateShop)
}

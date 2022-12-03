package main

import (
	"github.com/Astenna/Thesis_PoC/faas"
	"github.com/Astenna/Thesis_PoC/faas/types"
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

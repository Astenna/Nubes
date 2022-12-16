package main

import (
	"github.com/Astenna/Nubes/faas"
	"github.com/Astenna/Nubes/faas/types"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler5(id string) (types.Shop, error) {
	shop, err := faas.GetShop(id)
	return *shop, err
}

func main() {
	lambda.Start(Handler5)
}

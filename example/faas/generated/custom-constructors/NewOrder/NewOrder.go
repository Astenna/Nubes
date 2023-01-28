package main

import (
	org "github.com/Astenna/Nubes/example/faas/types"
	"github.com/aws/aws-lambda-go/lambda"
)

func NewOrderHandler(input org.Order) (org.Order, error) {
	result, _err := org.NewOrder(input)
	if _err != nil {
		return result, _err
	}
	return result, _err
}

func main() {
	lambda.Start(NewOrderHandler)
}

package main

import (
	org "github.com/Astenna/Nubes/example/faas/types"
	"github.com/Astenna/Nubes/lib"
	"github.com/aws/aws-lambda-go/lambda"
)

func DecreaseAvailabilityByHandler(input lib.HandlerParameters) error {
	instance := new(org.Product)
	instance.Id = input.Id
	instance.Init()

	_err := instance.DecreaseAvailabilityBy(input.Parameter.(int))
	if _err != nil {
		return _err
	}

	return _err
}

func main() {
	lambda.Start(DecreaseAvailabilityByHandler)
}

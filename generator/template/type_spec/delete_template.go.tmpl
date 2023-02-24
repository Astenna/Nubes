package main

import (
	"fmt"

	lib "github.com/Astenna/Nubes/lib"
	"github.com/aws/aws-lambda-go/lambda"
)

func DeleteHandler(input lib.HandlerParameters) error {
	if input.Id == "" {
		return fmt.Errorf("missing Id in HandlerParameters")
	}
	if input.TypeName == "" {
		return fmt.Errorf("missing TypeName in HandlerParameters")
	}

	err := lib.DeleteWithTypeNameAsArg(input.Id, input.TypeName)

	if err != nil {
		return fmt.Errorf("failed to delete type %s with id: %s. Error %w", input.TypeName, input.Id, err)
	}

	return nil
}

func main() {
	lambda.Start(DeleteHandler)
}

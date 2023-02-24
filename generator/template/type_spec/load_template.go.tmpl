package main

import (
	"fmt"

	lib "github.com/Astenna/Nubes/lib"
	"github.com/aws/aws-lambda-go/lambda"
)

func LoadHandler(input lib.HandlerParameters) error {
	if input.Id == "" {
		return fmt.Errorf("missing Id in HandlerParameters")
	}

	if input.TypeName == "" {
		return fmt.Errorf("missing TypeName in HandlerParameters")
	}

	libInput := lib.IsInstanceAlreadyCreatedParam{
		Id:       input.Id,
		TypeName: input.TypeName,
	}

	exists, err := lib.IsInstanceAlreadyCreated(libInput)

	if err != nil {
		return fmt.Errorf("error occurred while checking if %s with id %s exist. Error %w", input.Id, input.TypeName, err)
	}

	if !exists {
		return fmt.Errorf("instance of type %s with id %s does not exist", input.Id, input.TypeName)
	}

	return nil
}

func main() {
	lambda.Start(LoadHandler)
}

package client_lib

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
)

func NewDiscount() (*DiscountStub, error) {

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("NewDiscount")})
	if _err != nil {
		return nil, _err
	}
	if out.FunctionError != nil {
		return nil, fmt.Errorf("lambda function invocation failed. Error: %s", string(out.Payload))
	}

	result := new(DiscountStub)
	_err = json.Unmarshal(out.Payload, result)
	if _err != nil {
		return nil, _err
	}

	return result, nil
}

package client_lib

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
)

func NewOrder(input OrderStub) (*OrderStub, error) {
	jsonParam, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("NewOrder"), Payload: jsonParam})
	if _err != nil {
		return nil, _err
	}
	if out.FunctionError != nil {
		return nil, fmt.Errorf("lambda function invocation failed. Error: %s", string(out.Payload))
	}

	result := new(OrderStub)
	_err = json.Unmarshal(out.Payload, result)
	if _err != nil {
		return nil, _err
	}

	return result, nil
}

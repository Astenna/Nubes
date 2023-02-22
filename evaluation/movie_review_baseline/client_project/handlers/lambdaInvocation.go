package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

var sess = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))

var LambdaClient = lambda.New(sess)

func invokeLambdaToGetSingleItem[T any](input any, functionName string) (*T, error) {
	jsonParam, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String(functionName), Payload: jsonParam})
	if _err != nil {
		return nil, _err
	}
	if out.FunctionError != nil {
		return nil, fmt.Errorf("lambda function designed to verify if instance exists failed. Error: %s", string(out.Payload))
	}

	result := new(T)
	err = json.Unmarshal(out.Payload, result)
	if err != nil {
		return nil, err
	}
	return result, err
}

func invokeLambdaToGetList[T any](input any, functionName string) ([]T, error) {
	jsonParam, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String(functionName), Payload: jsonParam})
	if _err != nil {
		return nil, _err
	}
	if out.FunctionError != nil {
		return nil, fmt.Errorf("lambda function designed to verify if instance exists failed. Error: %s", string(out.Payload))
	}

	result := new([]T)
	err = json.Unmarshal(out.Payload, result)
	if err != nil {
		return nil, err
	}
	return *result, err
}

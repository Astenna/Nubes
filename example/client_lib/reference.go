package client_lib


import (
	"encoding/json"
	"fmt"
	"errors"

	"github.com/Astenna/Nubes/lib"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
)

// REFERENCE

type Reference[T lib.Nobject] string

func NewReference[T lib.Nobject](id string) *Reference[T] {
	result := Reference[T](id)
	return &result
}

func (r Reference[T]) Id() string {
	return string(r)
}

func (r Reference[T]) Get() (*T, error) {
	newInstance := new(T)

	params := lib.LoadBatchParam{
		Ids:      []string{string(r)},
		TypeName: (*newInstance).GetTypeName(),
	}
	jsonParam, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("Load"), Payload: jsonParam})
	if _err != nil {
		return nil, _err
	}
	if out.FunctionError != nil {
		return nil, fmt.Errorf("lambda function designed to verify if instance exists failed. Error: %s", string(out.Payload))
	}

	casted := any(newInstance)
	setIdInterf, _ := casted.(setId)
	setIdInterf.setId(string(r))
	setIdInterf.init()
	return newInstance, nil
}

func GetStub[T lib.Nobject](id string) (T, error) {
	if id == "" {
		return *new(T), errors.New("missing id")
	}

	result := new(T)
	params := lib.GetStateParam{
		Id:       id,
		TypeName: (*result).GetTypeName(),
		GetStub:  true,
	}
	jsonParam, err := json.Marshal(params)
	if err != nil {
		return *new(T), err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("GetState"), Payload: jsonParam})
	if _err != nil {
		return *new(T), _err
	}
	if out.FunctionError != nil {
		return *new(T), fmt.Errorf(string(out.Payload[:]))
	}

	err = json.Unmarshal(out.Payload, result)
	if err != nil {
		return *new(T), err
	}
	return *result, err
}

func GetStubs[T lib.Nobject](ids []string) ([]T, error) {
	if len(ids) < 1 {
		return nil, nil
	}

	params := lib.GetBatchParam{
		Ids:      ids,
		TypeName: (*new(T)).GetTypeName(),
	}
	jsonParam, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("GetBatch"), Payload: jsonParam})
	if _err != nil {
		return nil, _err
	}
	if out.FunctionError != nil {
		return nil, fmt.Errorf("lambda function designed retrieve the objects' states failed. Error: %s", string(out.Payload))
	}

	stubs := make([]T, len(ids))
	err = json.Unmarshal(out.Payload, &stubs)
	if err != nil {
		return nil, err
	}

	return stubs, err
}

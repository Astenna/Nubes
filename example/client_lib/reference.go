package client_lib


import (
	"encoding/json"
	"errors"
	"fmt"

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

func (r Reference[T]) GetStub() (interface{}, error) {
	if string(r) == "" {
		return *new(ProductStub), errors.New("id of the type not set, use  LoadProduct or ExportProduct to create new instance of the type")
	}

	instance := new(T)
	var factory map[string]interface{} = map[string]interface{}{
		"User": UserStub{},
	}

	params := lib.GetStateParam{
		Id:       string(r),
		TypeName: (*instance).GetTypeName(),
		GetStub:  true,
	}
	jsonParam, err := json.Marshal(params)
	if err != nil {
		return *new(ProductStub), err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("GetState"), Payload: jsonParam})
	if _err != nil {
		return *new(ProductStub), _err
	}
	if out.FunctionError != nil {
		return *new(ProductStub), fmt.Errorf(string(out.Payload[:]))
	}

	result := factory[params.TypeName]
	err = json.Unmarshal(out.Payload, &result)
	if err != nil {
		return *new(ProductStub), err
	}
	return result, err
}
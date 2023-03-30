package client_lib


import (
	"encoding/json"
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
package client_lib

import (
	"encoding/json"
	"fmt"

	"github.com/Astenna/Nubes/lib"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
)

type ReferenceList[T lib.Nobject] []string

func NewReferenceList[T lib.Nobject](ids []string) *ReferenceList[T] {
	result := ReferenceList[T](ids)
	return &result
}

func (r ReferenceList[T]) Ids() []string {
	return []string(r)
}

func (r ReferenceList[T]) Get() ([]T, error) {
	return loadBatch[T](r.Ids())
}

func loadBatch[T lib.Nobject](ids []string) ([]T, error) {
	if len(ids) == 0 {
		return []T{}, nil
	}

	params := lib.LoadBatchParam{
		Ids:      ids,
		TypeName: (*new(T)).GetTypeName(),
	}
	jsonParam, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	out, err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("Load"), Payload: jsonParam})
	if out.FunctionError != nil {
		return nil, fmt.Errorf("lambda function designed to verify if instance exists failed. Error: %s", string(out.Payload))
	}
	foundIds := ids
	if err != nil {
		if notFound, casted := err.(lib.NotFoundError); casted {
			foundIds = difference(ids, notFound.Ids)
		} else {
			return nil, err
		}
	}

	result := make([]T, len(foundIds))
	for i, id := range foundIds {
		newInstance := new(T)
		casted := any(newInstance)
		setIdInterf, _ := casted.(setId)
		setIdInterf.setId(id)
		setIdInterf.init()
		result[i] = *newInstance
	}
	return result, err
}

func difference(a, b []string) []string {
	set_b := make(map[string]struct{}, len(b))
	for _, x := range b {
		set_b[x] = struct{}{}
	}

	var diff []string
	for _, x := range a {
		if _, found := set_b[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}
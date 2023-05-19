package client_lib

import (
	"encoding/json"
	"fmt"

	"github.com/Astenna/Nubes/lib"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
)

type setId interface {
	setId(id string)
	init()
}

type referenceNavigationList[T lib.Nobject, Stub any] struct {
	param lib.ReferenceNavigationListParam
}

func newReferenceNavigationList[T lib.Nobject, Stub any](param lib.ReferenceNavigationListParam) *referenceNavigationList[T, Stub] {
	r := new(referenceNavigationList[T, Stub])
	r.param = param
	return r
}

func (r referenceNavigationList[T, Stub]) GetIds() ([]string, error) {
	jsonParam, err := json.Marshal(r.param)
	if err != nil {
		return nil, err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("ReferenceGetIds"), Payload: jsonParam})
	if _err != nil {
		return nil, _err
	}
	if out.FunctionError != nil {
		return nil, fmt.Errorf(string(out.Payload[:]))
	}

	var result []string
	err = json.Unmarshal(out.Payload, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r referenceNavigationList[T, Stub]) Get() ([]T, error) {
	jsonParam, err := json.Marshal(r.param)
	if err != nil {
		return nil, err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("ReferenceGet"), Payload: jsonParam})
	if out.FunctionError != nil {
		return nil, fmt.Errorf("lambda function designed to verify if instance exists failed. Error: %s", string(out.Payload))
	}

	var notFoundError lib.NotFoundError
	if _err != nil {
		casted := false
		if notFoundError, casted = err.(lib.NotFoundError); !casted {
			return nil, err
		}
	}

	var foundIds []string
	err = json.Unmarshal(out.Payload, &foundIds)
	if err != nil {
		return nil, err
	}

	if len(notFoundError.Ids) > 0 {
		foundIds = difference(foundIds, notFoundError.Ids)
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

func (r referenceNavigationList[T, Stub]) GetStubs() ([]Stub, error) {
	jsonParam, err := json.Marshal(r.param)
	if err != nil {
		return nil, err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("ReferenceGetStubs"), Payload: jsonParam})
	if _err != nil {
		return nil, _err
	}
	if out.FunctionError != nil {
		return nil, fmt.Errorf("lambda function designed to the retrieve objects' states failed. Error: %s", string(out.Payload))
	}

	var stubs []Stub
	err = json.Unmarshal(out.Payload, &stubs)
	if err != nil {
		return nil, err
	}

	return stubs, err
}

func (r referenceNavigationList[T, Stub]) AddToManyToMany(newId string) error {
	if newId == "" {
		return fmt.Errorf("missing id")
	}

	if r.param.IsManyToMany {
		params := lib.AddToManyToManyParam{
			RefNavListParam: r.param,
			NewId:           newId,
		}

		jsonParam, err := json.Marshal(params)
		if err != nil {
			return err
		}

		out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("ReferenceAddToManyToMany"), Payload: jsonParam})
		if _err != nil {
			return _err
		}
		if out.FunctionError != nil {
			return fmt.Errorf(string(out.Payload[:]))
		}

		return nil
	}

	return fmt.Errorf("can not add elements to ReferenceNavigationList of OneToMany relationship")
}

func (r referenceNavigationList[T, Stub]) DeleteFromManyToMany(ids []string) error {
	if len(ids) == 0 {
		return fmt.Errorf("missing ids to delete")
	}

	if r.param.IsManyToMany {
		jsonParam, err := json.Marshal(lib.DeleteFromManyToManyParam{
			RefNavListParam: r.param,
			IdsToDelete:     ids,
		})
		if err != nil {
			return err
		}

		out, _err := LambdaClient.Invoke(&lambda.InvokeInput{
			FunctionName: aws.String("ReferenceDeleteFromManyToMany"),
			Payload:      jsonParam,
		})
		if _err != nil {
			return _err
		}
		if out.FunctionError != nil {
			return fmt.Errorf(string(out.Payload[:]))
		}

		return nil
	}

	return fmt.Errorf("can not delete elements in ReferenceNavigationList of OneToMany relationship")
}

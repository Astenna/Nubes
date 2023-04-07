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
	setup lib.ReferenceNavigationListSetup[T]
}

func newReferenceNavigationList[T lib.Nobject, Stub any](ownerId, ownerTypeName, referringFieldName string, isManyToMany bool) *referenceNavigationList[T, Stub] {
	r := new(referenceNavigationList[T, Stub])
	r.setup = lib.NewReferenceNavigationListSetup[T](ownerId, ownerTypeName, referringFieldName, isManyToMany)
	return r
}

func (r referenceNavigationList[T, Stub]) GetIds() ([]string, error) {

	if r.setup.UsesIndex {
		out, err := r.getByIndex()
		return out, err
	}

	if r.setup.IsManyToMany && !r.setup.UsesIndex {
		out, err := r.getSortKeysByPartitionKey()
		return out, err
	}

	return nil, fmt.Errorf("invalid initialization of ReferenceNavigationList")
}

func (r referenceNavigationList[T, Stub]) Get() ([]T, error) {
	ids, err := r.GetIds()
	if err != nil {
		return nil, err
	}

	return loadBatch[T](ids)
}

func (r referenceNavigationList[T, Stub]) GetStubs() ([]Stub, error) {
	ids, err := r.GetIds()
	if err != nil {
		return nil, err
	}

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
		return nil, fmt.Errorf("lambda function designed to the objects' states failed. Error: %s", string(out.Payload))
	}

	stubs := make([]Stub, len(ids))
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

	if r.setup.IsManyToMany {
		typeName := (*new(T)).GetTypeName()
		params := lib.AddToManyToManyParam{
			TypeName:                     typeName,
			NewId:                        newId,
			InsertToManyToManyTableParam: r.setup.GetInsertToManyToManyTableParam(newId),
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

	if r.setup.IsManyToMany {
		params := r.setup.GetDeleteFromManyToManyParam(ids)

		jsonParam, err := json.Marshal(params)
		if err != nil {
			return err
		}

		out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("ReferenceDelteFromManyToMany"), Payload: jsonParam})
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

func (r referenceNavigationList[T, Stub]) getByIndex() ([]string, error) {

	jsonParam, err := json.Marshal(r.setup.GetQueryByIndexParam())
	if err != nil {
		return nil, err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("ReferenceGetByIndex"), Payload: jsonParam})
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

func (r referenceNavigationList[T, Stub]) getSortKeysByPartitionKey() ([]string, error) {

	param, err := r.setup.GetQueryByPartitionKeyParam()
	if err != nil {
		return nil, err
	}
	jsonParam, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("ReferenceGetSortKeysByPartitionKey"), Payload: jsonParam})
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
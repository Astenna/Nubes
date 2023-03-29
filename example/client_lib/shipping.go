package client_lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/Astenna/Nubes/lib"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
)

type shipping struct {
	id string
}

// ALL THE CODE BELOW IS GENERATED ONLY FOR NOBJECTS TYPES
func (shipping) GetTypeName() string {
	return "Shipping"
}

// LOAD AND EXPORT

func LoadShipping(id string) (*shipping, error) {
	newInstance := new(shipping)

	params := lib.LoadBatchParam{
		Ids:      []string{id},
		TypeName: newInstance.GetTypeName(),
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

	newInstance.id = id
	newInstance.init()
	return newInstance, nil
}

func loadShippingWithoutCheckIfExists(id string) *shipping {
	newInstance := new(shipping)
	newInstance.id = id
	return newInstance
}

// setId interface for initilization in ReferenceNavigationList
func (u *shipping) setId(id string) {
	u.id = id
}

func (r *shipping) init() {

}

func ExportShipping(input ShippingStub) (*shipping, error) {
	newInstance := new(shipping)

	params := lib.HandlerParameters{
		TypeName:  newInstance.GetTypeName(),
		Parameter: input,
	}
	jsonParam, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("Export"), Payload: jsonParam})
	if _err != nil {
		return nil, _err
	}
	if out.FunctionError != nil {
		return nil, fmt.Errorf("lambda function designed to verify if instance exists failed. Error: %s", string(out.Payload[:]))
	}

	newInstance.id, err = strconv.Unquote(string(out.Payload[:]))
	newInstance.init()
	return newInstance, err
}

// DELETE

func DeleteShipping(id string) error {
	newInstance := new(shipping)

	params := lib.HandlerParameters{
		Id:       id,
		TypeName: newInstance.GetTypeName(),
	}
	jsonParam, err := json.Marshal(params)
	if err != nil {
		return err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("Delete"), Payload: jsonParam})
	if _err != nil {
		return _err
	}
	if out.FunctionError != nil {
		return fmt.Errorf("lambda function failed. Error: %s", string(out.Payload))
	}

	return nil
}

// GETID

func (s shipping) GetId() string {
	return s.id
}

// REFERENCE

func (s shipping) Reference() Reference[shipping] {
	return *NewReference[shipping](s.GetId())
}

// GETTERS AND SETTERS

func (s shipping) GetAddress() (string, error) {
	if s.id == "" {
		return *new(string), errors.New("id of the type not set, use  LoadShipping or ExportShipping to create new instance of the type")
	}

	params := lib.GetStateParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "Address",
	}
	jsonParam, err := json.Marshal(params)
	if err != nil {
		return *new(string), err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("GetState"), Payload: jsonParam})
	if _err != nil {
		return *new(string), _err
	}
	if out.FunctionError != nil {
		return *new(string), fmt.Errorf(string(out.Payload[:]))
	}

	result := new(string)
	err = json.Unmarshal(out.Payload, result)
	if err != nil {
		return *new(string), err
	}
	return *result, err

}

func (s shipping) SetAddress(newValue string) error {
	if s.id == "" {
		return errors.New("id of the type not set, use LoadShipping or ExportShipping to create new instance of the type")
	}

	params := lib.SetFieldParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "Address",
		Value:     newValue,
	}
	jsonParam, err := json.Marshal(params)
	if err != nil {
		return err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("SetField"), Payload: jsonParam})
	if _err != nil {
		return _err
	}
	if out.FunctionError != nil {
		return fmt.Errorf(string(out.Payload[:]))
	}
	return nil
}

func (s shipping) GetState() (ShippingState, error) {
	if s.id == "" {
		return *new(ShippingState), errors.New("id of the type not set, use  LoadShipping or ExportShipping to create new instance of the type")
	}

	params := lib.GetStateParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "State",
	}
	jsonParam, err := json.Marshal(params)
	if err != nil {
		return *new(ShippingState), err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("GetState"), Payload: jsonParam})
	if _err != nil {
		return *new(ShippingState), _err
	}
	if out.FunctionError != nil {
		return *new(ShippingState), fmt.Errorf(string(out.Payload[:]))
	}

	result := new(ShippingState)
	err = json.Unmarshal(out.Payload, result)
	if err != nil {
		return *new(ShippingState), err
	}
	return *result, err

}

func (s shipping) SetState(newValue ShippingState) error {
	if s.id == "" {
		return errors.New("id of the type not set, use LoadShipping or ExportShipping to create new instance of the type")
	}

	params := lib.SetFieldParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "State",
		Value:     newValue,
	}
	jsonParam, err := json.Marshal(params)
	if err != nil {
		return err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("SetField"), Payload: jsonParam})
	if _err != nil {
		return _err
	}
	if out.FunctionError != nil {
		return fmt.Errorf(string(out.Payload[:]))
	}
	return nil
}

// (STATE-CHANGING) METHODS

func (r shipping) GetStub() (ShippingStub, error) {
	if r.id == "" {
		return *new(ShippingStub), errors.New("id of the type not set, use  LoadShipping or ExportShipping to create new instance of the type")
	}

	params := lib.GetStateParam{
		Id:       r.GetId(),
		TypeName: r.GetTypeName(),
		GetStub:  true,
	}
	jsonParam, err := json.Marshal(params)
	if err != nil {
		return *new(ShippingStub), err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("GetState"), Payload: jsonParam})
	if _err != nil {
		return *new(ShippingStub), _err
	}
	if out.FunctionError != nil {
		return *new(ShippingStub), fmt.Errorf(string(out.Payload[:]))
	}

	result := new(ShippingStub)
	err = json.Unmarshal(out.Payload, result)
	if err != nil {
		return *new(ShippingStub), err
	}
	return *result, err
}

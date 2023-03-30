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

type order struct {
	id string
}

// ALL THE CODE BELOW IS GENERATED ONLY FOR NOBJECTS TYPES
func (order) GetTypeName() string {
	return "Order"
}

// LOAD AND EXPORT

func LoadOrder(id string) (*order, error) {
	newInstance := new(order)

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

func loadOrderWithoutCheckIfExists(id string) *order {
	newInstance := new(order)
	newInstance.id = id
	return newInstance
}

// setId interface for initilization in ReferenceNavigationList
func (u *order) setId(id string) {
	u.id = id
}

func (r *order) init() {

}

func ExportOrder(input OrderStub) (*order, error) {
	newInstance := new(order)

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

func DeleteOrder(id string) error {
	newInstance := new(order)

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

func (s order) GetId() string {
	return s.id
}

// REFERENCE

func (s order) AsReference() Reference[order] {
	return *NewReference[order](s.GetId())
}

// GETTERS AND SETTERS

func (s order) GetProducts() ([]OrderedProduct, error) {
	if s.id == "" {
		return *new([]OrderedProduct), errors.New("id of the type not set, use  LoadOrder or ExportOrder to create new instance of the type")
	}

	params := lib.GetStateParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "Products",
	}
	jsonParam, err := json.Marshal(params)
	if err != nil {
		return *new([]OrderedProduct), err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("GetState"), Payload: jsonParam})
	if _err != nil {
		return *new([]OrderedProduct), _err
	}
	if out.FunctionError != nil {
		return *new([]OrderedProduct), fmt.Errorf(string(out.Payload[:]))
	}

	result := new([]OrderedProduct)
	err = json.Unmarshal(out.Payload, result)
	if err != nil {
		return *new([]OrderedProduct), err
	}
	return *result, err

}

func (s order) SetProducts(newValue []OrderedProduct) error {
	if s.id == "" {
		return errors.New("id of the type not set, use LoadOrder or ExportOrder to create new instance of the type")
	}

	params := lib.SetFieldParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "Products",
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

func (s order) GetBuyer() (user, error) {
	if s.id == "" {
		return *new(user), errors.New("id of the type not set, use  LoadOrder or ExportOrder to create new instance of the type")
	}

	params := lib.GetStateParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "Buyer",
	}
	jsonParam, err := json.Marshal(params)
	if err != nil {
		return *new(user), err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("GetState"), Payload: jsonParam})
	if _err != nil {
		return *new(user), _err
	}
	if out.FunctionError != nil {
		return *new(user), fmt.Errorf(string(out.Payload[:]))
	}

	result := new(lib.Reference[user])
	err = json.Unmarshal(out.Payload, result)
	if err != nil {
		return *new(user), err
	}
	var referenceResult = loadUserWithoutCheckIfExists(result.Id())
	return *referenceResult, err

}

func (s order) GetBuyerId() (string, error) {
	if s.id == "" {
		return "", errors.New("id of the type not set, use  LoadOrder or ExportOrder to create new instance of the type")
	}

	params := lib.GetStateParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "Buyer",
	}
	jsonParam, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("GetState"), Payload: jsonParam})
	if _err != nil {
		return "", _err
	}
	if out.FunctionError != nil {
		return "", fmt.Errorf(string(out.Payload[:]))
	}

	result := new(lib.Reference[user])
	err = json.Unmarshal(out.Payload, result)
	if err != nil {
		return "", err
	}

	return result.Id(), err
}

func (s order) SetBuyer(newValue string) error {
	if s.id == "" {
		return errors.New("id of the type not set, use LoadOrder or ExportOrder to create new instance of the type")
	}

	params := lib.SetFieldParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "Buyer",
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

func (s order) GetShipping() (shipping, error) {
	if s.id == "" {
		return *new(shipping), errors.New("id of the type not set, use  LoadOrder or ExportOrder to create new instance of the type")
	}

	params := lib.GetStateParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "Shipping",
	}
	jsonParam, err := json.Marshal(params)
	if err != nil {
		return *new(shipping), err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("GetState"), Payload: jsonParam})
	if _err != nil {
		return *new(shipping), _err
	}
	if out.FunctionError != nil {
		return *new(shipping), fmt.Errorf(string(out.Payload[:]))
	}

	result := new(lib.Reference[shipping])
	err = json.Unmarshal(out.Payload, result)
	if err != nil {
		return *new(shipping), err
	}
	var referenceResult = loadShippingWithoutCheckIfExists(result.Id())
	return *referenceResult, err

}

func (s order) GetShippingId() (string, error) {
	if s.id == "" {
		return "", errors.New("id of the type not set, use  LoadOrder or ExportOrder to create new instance of the type")
	}

	params := lib.GetStateParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "Shipping",
	}
	jsonParam, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("GetState"), Payload: jsonParam})
	if _err != nil {
		return "", _err
	}
	if out.FunctionError != nil {
		return "", fmt.Errorf(string(out.Payload[:]))
	}

	result := new(lib.Reference[shipping])
	err = json.Unmarshal(out.Payload, result)
	if err != nil {
		return "", err
	}

	return result.Id(), err
}

func (s order) SetShipping(newValue string) error {
	if s.id == "" {
		return errors.New("id of the type not set, use LoadOrder or ExportOrder to create new instance of the type")
	}

	params := lib.SetFieldParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "Shipping",
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

func (r order) GetStub() (OrderStub, error) {
	if r.id == "" {
		return *new(OrderStub), errors.New("id of the type not set, use  LoadOrder or ExportOrder to create new instance of the type")
	}

	params := lib.GetStateParam{
		Id:       r.GetId(),
		TypeName: r.GetTypeName(),
		GetStub:  true,
	}
	jsonParam, err := json.Marshal(params)
	if err != nil {
		return *new(OrderStub), err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("GetState"), Payload: jsonParam})
	if _err != nil {
		return *new(OrderStub), _err
	}
	if out.FunctionError != nil {
		return *new(OrderStub), fmt.Errorf(string(out.Payload[:]))
	}

	result := new(OrderStub)
	err = json.Unmarshal(out.Payload, result)
	if err != nil {
		return *new(OrderStub), err
	}
	return *result, err
}

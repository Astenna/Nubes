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

type product struct {
	id string
}

// ALL THE CODE BELOW IS GENERATED ONLY FOR NOBJECTS TYPES
func (product) GetTypeName() string {
	return "Product"
}

// LOAD AND EXPORT

func LoadProduct(id string) (*product, error) {
	newInstance := new(product)

	params := lib.HandlerParameters{
		Id:       id,
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

func loadProductWithoutCheckIfExists(id string) *product {
	newInstance := new(product)
	newInstance.id = id
	return newInstance
}

// setId interface for initilization in ReferenceNavigationList
func (u *product) setId(id string) {
	u.id = id
}

func (r *product) init() {

}

func ExportProduct(input ProductStub) (*product, error) {
	newInstance := new(product)

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

func DeleteProduct(id string) error {
	newInstance := new(product)

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

func (s product) GetId() string {
	return s.id
}

// GETTERS AND SETTERS

func (s product) GetName() (string, error) {
	if s.id == "" {
		return *new(string), errors.New("id of the type not set, use  LoadProduct or ExportProduct to create new instance of the type")
	}

	params := lib.GetStateParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "Name",
	}
	jsonParam, err := json.Marshal(params)
	if err != nil {
		return *new(string), err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("GetField"), Payload: jsonParam})
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

func (s product) SetName(newValue string) error {
	if s.id == "" {
		return errors.New("id of the type not set, use LoadProduct or ExportProduct to create new instance of the type")
	}

	params := lib.SetFieldParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "Name",
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

func (s product) GetQuantityAvailable() (int, error) {
	if s.id == "" {
		return *new(int), errors.New("id of the type not set, use  LoadProduct or ExportProduct to create new instance of the type")
	}

	params := lib.GetStateParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "QuantityAvailable",
	}
	jsonParam, err := json.Marshal(params)
	if err != nil {
		return *new(int), err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("GetField"), Payload: jsonParam})
	if _err != nil {
		return *new(int), _err
	}
	if out.FunctionError != nil {
		return *new(int), fmt.Errorf(string(out.Payload[:]))
	}

	result := new(int)
	err = json.Unmarshal(out.Payload, result)
	if err != nil {
		return *new(int), err
	}
	return *result, err

}

func (s product) SetQuantityAvailable(newValue int) error {
	if s.id == "" {
		return errors.New("id of the type not set, use LoadProduct or ExportProduct to create new instance of the type")
	}

	params := lib.SetFieldParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "QuantityAvailable",
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

func (s product) GetSoldBy() (shop, error) {
	if s.id == "" {
		return *new(shop), errors.New("id of the type not set, use  LoadProduct or ExportProduct to create new instance of the type")
	}

	params := lib.GetStateParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "SoldBy",
	}
	jsonParam, err := json.Marshal(params)
	if err != nil {
		return *new(shop), err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("GetField"), Payload: jsonParam})
	if _err != nil {
		return *new(shop), _err
	}
	if out.FunctionError != nil {
		return *new(shop), fmt.Errorf(string(out.Payload[:]))
	}

	result := new(lib.Reference[shop])
	err = json.Unmarshal(out.Payload, result)
	if err != nil {
		return *new(shop), err
	}
	var referenceResult = loadShopWithoutCheckIfExists(result.Id())
	return *referenceResult, err

}

func (s product) GetSoldById() (string, error) {
	if s.id == "" {
		return "", errors.New("id of the type not set, use  LoadProduct or ExportProduct to create new instance of the type")
	}

	params := lib.GetStateParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "SoldBy",
	}
	jsonParam, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("GetField"), Payload: jsonParam})
	if _err != nil {
		return "", _err
	}
	if out.FunctionError != nil {
		return "", fmt.Errorf(string(out.Payload[:]))
	}

	result := new(lib.Reference[shop])
	err = json.Unmarshal(out.Payload, result)
	if err != nil {
		return "", err
	}

	return result.Id(), err
}

func (s product) SetSoldBy(newValue string) error {
	if s.id == "" {
		return errors.New("id of the type not set, use LoadProduct or ExportProduct to create new instance of the type")
	}

	params := lib.SetFieldParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "SoldBy",
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

// REFENCE LIST GETTER - returns ids
func (s product) GetDiscountIds() ([]string, error) {
	if s.id == "" {
		return nil, errors.New("id of the type not set, use LoadProduct or ExportProduct to create new instance of the type")
	}

	params := lib.GetStateParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "Discount",
	}
	jsonParam, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("GetField"), Payload: jsonParam})
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

// REFERENCE LIST GETTER - returns initialized isntances
func (s product) GetDiscount() ([]discount, error) {
	if s.id == "" {
		return nil, errors.New("id of the type not set, use LoadProduct or ExportProduct to create new instance of the type")
	}

	params := lib.GetStateParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "Discount",
	}
	jsonParam, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("GetField"), Payload: jsonParam})
	if _err != nil {
		return nil, _err
	}
	if out.FunctionError != nil {
		return nil, fmt.Errorf(string(out.Payload[:]))
	}

	var ids []string
	err = json.Unmarshal(out.Payload, &ids)
	if err != nil {
		return nil, err
	}

	result := make([]discount, len(ids))
	for index, id := range ids {
		instance := loadDiscountWithoutCheckIfExists(id)
		result[index] = *instance
	}

	return result, nil
}
func (s product) SetDiscount(ids []string) error {
	if s.id == "" {
		return errors.New("id of the type not set, use LoadProduct or ExportProduct to create new instance of the type")
	}

	params := lib.SetFieldParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "Discount",
		Value:     ids,
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

func (s product) GetPrice() (float64, error) {
	if s.id == "" {
		return *new(float64), errors.New("id of the type not set, use  LoadProduct or ExportProduct to create new instance of the type")
	}

	params := lib.GetStateParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "Price",
	}
	jsonParam, err := json.Marshal(params)
	if err != nil {
		return *new(float64), err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("GetField"), Payload: jsonParam})
	if _err != nil {
		return *new(float64), _err
	}
	if out.FunctionError != nil {
		return *new(float64), fmt.Errorf(string(out.Payload[:]))
	}

	result := new(float64)
	err = json.Unmarshal(out.Payload, result)
	if err != nil {
		return *new(float64), err
	}
	return *result, err

}

func (s product) SetPrice(newValue float64) error {
	if s.id == "" {
		return errors.New("id of the type not set, use LoadProduct or ExportProduct to create new instance of the type")
	}

	params := lib.SetFieldParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "Price",
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

func (p product) DecreaseAvailabilityBy(input int) error {
	if p.id == "" {
		return errors.New("id of the type not set, use  LoadProduct or ExportProduct to create new instance of the type")
	}

	params := new(lib.HandlerParameters)
	params.Id = p.id
	params.Parameter = input

	jsonParam, err := json.Marshal(params)
	if err != nil {
		return err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("ProductDecreaseAvailabilityBy"), Payload: jsonParam})
	if _err != nil {
		return _err
	}
	if out.FunctionError != nil {
		return fmt.Errorf(string(out.Payload[:]))
	}

	return _err
}

func (r product) GetStub() (ProductStub, error) {
	if r.id == "" {
		return *new(ProductStub), errors.New("id of the type not set, use  LoadProduct or ExportProduct to create new instance of the type")
	}

	params := lib.GetStateParam{
		Id:       r.GetId(),
		TypeName: r.GetTypeName(),
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

	result := new(ProductStub)
	err = json.Unmarshal(out.Payload, result)
	if err != nil {
		return *new(ProductStub), err
	}
	return *result, err
}

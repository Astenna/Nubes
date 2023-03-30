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

type user struct {
	id string

	Shops referenceNavigationList[shop, ShopStub]
}

// ALL THE CODE BELOW IS GENERATED ONLY FOR NOBJECTS TYPES
func (user) GetTypeName() string {
	return "User"
}

// LOAD AND EXPORT

func LoadUser(id string) (*user, error) {
	newInstance := new(user)

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

func loadUserWithoutCheckIfExists(id string) *user {
	newInstance := new(user)
	newInstance.id = id
	return newInstance
}

// setId interface for initilization in ReferenceNavigationList
func (u *user) setId(id string) {
	u.id = id
}

func (r *user) init() {

	r.Shops = *newReferenceNavigationList[shop, ShopStub](r.id, r.GetTypeName(), "", true)

}

func ExportUser(input UserStub) (*user, error) {
	newInstance := new(user)

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

func DeleteUser(id string) error {
	newInstance := new(user)

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

func (s user) GetId() string {
	return s.id
}

// REFERENCE

func (s user) AsReference() Reference[user] {
	return *NewReference[user](s.GetId())
}

// GETTERS AND SETTERS

func (s user) GetFirstName() (string, error) {
	if s.id == "" {
		return *new(string), errors.New("id of the type not set, use  LoadUser or ExportUser to create new instance of the type")
	}

	params := lib.GetStateParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "FirstName",
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

func (s user) SetFirstName(newValue string) error {
	if s.id == "" {
		return errors.New("id of the type not set, use LoadUser or ExportUser to create new instance of the type")
	}

	params := lib.SetFieldParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "FirstName",
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

func (s user) GetLastName() (string, error) {
	if s.id == "" {
		return *new(string), errors.New("id of the type not set, use  LoadUser or ExportUser to create new instance of the type")
	}

	params := lib.GetStateParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "LastName",
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

func (s user) SetLastName(newValue string) error {
	if s.id == "" {
		return errors.New("id of the type not set, use LoadUser or ExportUser to create new instance of the type")
	}

	params := lib.SetFieldParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "LastName",
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

func (s user) GetEmail() (string, error) {
	if s.id == "" {
		return *new(string), errors.New("id of the type not set, use  LoadUser or ExportUser to create new instance of the type")
	}

	params := lib.GetStateParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "Email",
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

func (s user) GetPassword() (string, error) {
	if s.id == "" {
		return *new(string), errors.New("id of the type not set, use  LoadUser or ExportUser to create new instance of the type")
	}

	params := lib.GetStateParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "Password",
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

func (s user) GetAddress() (string, error) {
	if s.id == "" {
		return *new(string), errors.New("id of the type not set, use  LoadUser or ExportUser to create new instance of the type")
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

func (s user) SetAddress(newValue string) error {
	if s.id == "" {
		return errors.New("id of the type not set, use LoadUser or ExportUser to create new instance of the type")
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

// REFENCE LIST GETTER - returns ids
func (s user) GetOrdersIds() ([]string, error) {
	if s.id == "" {
		return nil, errors.New("id of the type not set, use LoadUser or ExportUser to create new instance of the type")
	}

	params := lib.GetStateParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "Orders",
	}
	jsonParam, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("GetState"), Payload: jsonParam})
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
func (s user) GetOrders() ([]order, error) {
	if s.id == "" {
		return nil, errors.New("id of the type not set, use LoadUser or ExportUser to create new instance of the type")
	}

	params := lib.GetStateParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "Orders",
	}
	jsonParam, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("GetState"), Payload: jsonParam})
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

	result := make([]order, len(ids))
	for index, id := range ids {
		instance := loadOrderWithoutCheckIfExists(id)
		result[index] = *instance
	}

	return result, nil
}
func (s user) SetOrders(ids []string) error {
	if s.id == "" {
		return errors.New("id of the type not set, use LoadUser or ExportUser to create new instance of the type")
	}

	params := lib.SetFieldParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "Orders",
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

// (STATE-CHANGING) METHODS

func (u user) VerifyPassword(input string) (bool, error) {
	if u.id == "" {
		return *new(bool), errors.New("id of the type not set, use  LoadUser or ExportUser to create new instance of the type")
	}

	params := new(lib.HandlerParameters)
	params.Id = u.id
	params.Parameter = input

	jsonParam, err := json.Marshal(params)
	if err != nil {
		return *new(bool), err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("UserVerifyPassword"), Payload: jsonParam})
	if _err != nil {
		return *new(bool), _err
	}
	if out.FunctionError != nil {
		return *new(bool), fmt.Errorf(string(out.Payload[:]))
	}

	result := new(bool)

	_err = json.Unmarshal(out.Payload, result)
	if _err != nil {
		return *new(bool), err
	}

	return *result, _err
}

func (r user) GetStub() (UserStub, error) {
	if r.id == "" {
		return *new(UserStub), errors.New("id of the type not set, use  LoadUser or ExportUser to create new instance of the type")
	}

	params := lib.GetStateParam{
		Id:       r.GetId(),
		TypeName: r.GetTypeName(),
		GetStub:  true,
	}
	jsonParam, err := json.Marshal(params)
	if err != nil {
		return *new(UserStub), err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("GetState"), Payload: jsonParam})
	if _err != nil {
		return *new(UserStub), _err
	}
	if out.FunctionError != nil {
		return *new(UserStub), fmt.Errorf(string(out.Payload[:]))
	}

	result := new(UserStub)
	err = json.Unmarshal(out.Payload, result)
	if err != nil {
		return *new(UserStub), err
	}
	return *result, err
}

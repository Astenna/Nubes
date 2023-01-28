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

type account struct {
	id string
}

// ALL THE CODE BELOW IS GENERATED ONLY FOR NOBJECTS TYPES
func (account) GetTypeName() string {
	return "Account"
}

// LOAD AND EXPORT

func LoadAccount(id string) (*account, error) {
	newInstance := new(account)

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

func loadAccountWithoutCheckIfExists(id string) *account {
	newInstance := new(account)
	newInstance.id = id
	return newInstance
}

// setId interface for initilization in ReferenceNavigationList
func (u *account) setId(id string) {
	u.id = id
}

func (r *account) init() {

}

func ExportAccount(input AccountStub) (*account, error) {
	newInstance := new(account)

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

func DeleteAccount(id string) error {
	newInstance := new(account)

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

func (s account) GetId() string {
	return s.id
}

// GETTERS AND SETTERS

func (s account) GetNickname() (string, error) {
	if s.id == "" {
		return *new(string), errors.New("id of the type not set, use  LoadAccount or ExportAccount to create new instance of the type")
	}

	params := lib.GetStateParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "Nickname",
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

func (s account) SetNickname(newValue string) error {
	if s.id == "" {
		return errors.New("id of the type not set, use LoadAccount or ExportAccount to create new instance of the type")
	}

	params := lib.SetFieldParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "Nickname",
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

func (s account) GetEmail() (string, error) {
	if s.id == "" {
		return *new(string), errors.New("id of the type not set, use  LoadAccount or ExportAccount to create new instance of the type")
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

func (s account) GetPassword() (string, error) {
	if s.id == "" {
		return *new(string), errors.New("id of the type not set, use  LoadAccount or ExportAccount to create new instance of the type")
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

// (STATE-CHANGING) METHODS

func (u account) VerifyPassword(input string) (bool, error) {
	if u.id == "" {
		return *new(bool), errors.New("id of the type not set, use  LoadAccount or ExportAccount to create new instance of the type")
	}

	params := new(lib.HandlerParameters)
	params.Id = u.id
	params.Parameter = input

	jsonParam, err := json.Marshal(params)
	if err != nil {
		return *new(bool), err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("AccountVerifyPassword"), Payload: jsonParam})
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

func (r account) GetStub() (AccountStub, error) {
	if r.id == "" {
		return *new(AccountStub), errors.New("id of the type not set, use  LoadAccount or ExportAccount to create new instance of the type")
	}

	params := lib.GetStateParam{
		Id:       r.GetId(),
		TypeName: r.GetTypeName(),
		GetStub:  true,
	}
	jsonParam, err := json.Marshal(params)
	if err != nil {
		return *new(AccountStub), err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("GetState"), Payload: jsonParam})
	if _err != nil {
		return *new(AccountStub), _err
	}
	if out.FunctionError != nil {
		return *new(AccountStub), fmt.Errorf(string(out.Payload[:]))
	}

	result := new(AccountStub)
	err = json.Unmarshal(out.Payload, result)
	if err != nil {
		return *new(AccountStub), err
	}
	return *result, err
}

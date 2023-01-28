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

type movie struct {
	id string

	Reviews referenceNavigationList[review, ReviewStub]
}

// ALL THE CODE BELOW IS GENERATED ONLY FOR NOBJECTS TYPES
func (movie) GetTypeName() string {
	return "Movie"
}

// LOAD AND EXPORT

func LoadMovie(id string) (*movie, error) {
	newInstance := new(movie)

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

func loadMovieWithoutCheckIfExists(id string) *movie {
	newInstance := new(movie)
	newInstance.id = id
	return newInstance
}

// setId interface for initilization in ReferenceNavigationList
func (u *movie) setId(id string) {
	u.id = id
}

func (r *movie) init() {

	r.Reviews = *newReferenceNavigationList[review, ReviewStub](r.id, r.GetTypeName(), "Movie", false)

}

func ExportMovie(input MovieStub) (*movie, error) {
	newInstance := new(movie)

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

func DeleteMovie(id string) error {
	newInstance := new(movie)

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

func (s movie) GetId() string {
	return s.id
}

// GETTERS AND SETTERS

func (s movie) GetTitle() (string, error) {
	if s.id == "" {
		return *new(string), errors.New("id of the type not set, use  LoadMovie or ExportMovie to create new instance of the type")
	}

	params := lib.GetStateParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "Title",
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

func (s movie) SetTitle(newValue string) error {
	if s.id == "" {
		return errors.New("id of the type not set, use LoadMovie or ExportMovie to create new instance of the type")
	}

	params := lib.SetFieldParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "Title",
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

func (s movie) GetProductionYear() (int, error) {
	if s.id == "" {
		return *new(int), errors.New("id of the type not set, use  LoadMovie or ExportMovie to create new instance of the type")
	}

	params := lib.GetStateParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "ProductionYear",
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

func (s movie) SetProductionYear(newValue int) error {
	if s.id == "" {
		return errors.New("id of the type not set, use LoadMovie or ExportMovie to create new instance of the type")
	}

	params := lib.SetFieldParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "ProductionYear",
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

func (s movie) GetCategory() (category, error) {
	if s.id == "" {
		return *new(category), errors.New("id of the type not set, use  LoadMovie or ExportMovie to create new instance of the type")
	}

	params := lib.GetStateParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "Category",
	}
	jsonParam, err := json.Marshal(params)
	if err != nil {
		return *new(category), err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("GetField"), Payload: jsonParam})
	if _err != nil {
		return *new(category), _err
	}
	if out.FunctionError != nil {
		return *new(category), fmt.Errorf(string(out.Payload[:]))
	}

	result := new(lib.Reference[category])
	err = json.Unmarshal(out.Payload, result)
	if err != nil {
		return *new(category), err
	}
	var referenceResult = loadCategoryWithoutCheckIfExists(result.Id())
	return *referenceResult, err

}

func (s movie) GetCategoryId() (string, error) {
	if s.id == "" {
		return "", errors.New("id of the type not set, use  LoadMovie or ExportMovie to create new instance of the type")
	}

	params := lib.GetStateParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "Category",
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

	result := new(lib.Reference[category])
	err = json.Unmarshal(out.Payload, result)
	if err != nil {
		return "", err
	}

	return result.Id(), err
}

func (s movie) SetCategory(newValue string) error {
	if s.id == "" {
		return errors.New("id of the type not set, use LoadMovie or ExportMovie to create new instance of the type")
	}

	params := lib.SetFieldParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "Category",
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

func (r movie) GetStub() (MovieStub, error) {
	if r.id == "" {
		return *new(MovieStub), errors.New("id of the type not set, use  LoadMovie or ExportMovie to create new instance of the type")
	}

	params := lib.GetStateParam{
		Id:       r.GetId(),
		TypeName: r.GetTypeName(),
		GetStub:  true,
	}
	jsonParam, err := json.Marshal(params)
	if err != nil {
		return *new(MovieStub), err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("GetState"), Payload: jsonParam})
	if _err != nil {
		return *new(MovieStub), _err
	}
	if out.FunctionError != nil {
		return *new(MovieStub), fmt.Errorf(string(out.Payload[:]))
	}

	result := new(MovieStub)
	err = json.Unmarshal(out.Payload, result)
	if err != nil {
		return *new(MovieStub), err
	}
	return *result, err
}

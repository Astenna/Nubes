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

type review struct {
	id string
}

// ALL THE CODE BELOW IS GENERATED ONLY FOR NOBJECTS TYPES
func (review) GetTypeName() string {
	return "Review"
}

// LOAD AND EXPORT

func LoadReview(id string) (*review, error) {
	newInstance := new(review)

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

func loadReviewWithoutCheckIfExists(id string) *review {
	newInstance := new(review)
	newInstance.id = id
	return newInstance
}

// setId interface for initilization in ReferenceNavigationList
func (u *review) setId(id string) {
	u.id = id
}

func (r *review) init() {

}

func ExportReview(input ReviewStub) (*review, error) {
	newInstance := new(review)

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

func DeleteReview(id string) error {
	newInstance := new(review)

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

func (s review) GetId() string {
	return s.id
}

// GETTERS AND SETTERS

func (s review) GetRating() (int, error) {
	if s.id == "" {
		return *new(int), errors.New("id of the type not set, use  LoadReview or ExportReview to create new instance of the type")
	}

	params := lib.GetStateParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "Rating",
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

func (s review) SetRating(newValue int) error {
	if s.id == "" {
		return errors.New("id of the type not set, use LoadReview or ExportReview to create new instance of the type")
	}

	params := lib.SetFieldParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "Rating",
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

func (s review) GetMovie() (movie, error) {
	if s.id == "" {
		return *new(movie), errors.New("id of the type not set, use  LoadReview or ExportReview to create new instance of the type")
	}

	params := lib.GetStateParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "Movie",
	}
	jsonParam, err := json.Marshal(params)
	if err != nil {
		return *new(movie), err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("GetField"), Payload: jsonParam})
	if _err != nil {
		return *new(movie), _err
	}
	if out.FunctionError != nil {
		return *new(movie), fmt.Errorf(string(out.Payload[:]))
	}

	result := new(lib.Reference[movie])
	err = json.Unmarshal(out.Payload, result)
	if err != nil {
		return *new(movie), err
	}
	var referenceResult = loadMovieWithoutCheckIfExists(result.Id())
	return *referenceResult, err

}

func (s review) GetMovieId() (string, error) {
	if s.id == "" {
		return "", errors.New("id of the type not set, use  LoadReview or ExportReview to create new instance of the type")
	}

	params := lib.GetStateParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "Movie",
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

	result := new(lib.Reference[movie])
	err = json.Unmarshal(out.Payload, result)
	if err != nil {
		return "", err
	}

	return result.Id(), err
}

func (s review) SetMovie(newValue string) error {
	if s.id == "" {
		return errors.New("id of the type not set, use LoadReview or ExportReview to create new instance of the type")
	}

	params := lib.SetFieldParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "Movie",
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

func (s review) GetReviewer() (account, error) {
	if s.id == "" {
		return *new(account), errors.New("id of the type not set, use  LoadReview or ExportReview to create new instance of the type")
	}

	params := lib.GetStateParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "Reviewer",
	}
	jsonParam, err := json.Marshal(params)
	if err != nil {
		return *new(account), err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("GetField"), Payload: jsonParam})
	if _err != nil {
		return *new(account), _err
	}
	if out.FunctionError != nil {
		return *new(account), fmt.Errorf(string(out.Payload[:]))
	}

	result := new(lib.Reference[account])
	err = json.Unmarshal(out.Payload, result)
	if err != nil {
		return *new(account), err
	}
	var referenceResult = loadAccountWithoutCheckIfExists(result.Id())
	return *referenceResult, err

}

func (s review) GetReviewerId() (string, error) {
	if s.id == "" {
		return "", errors.New("id of the type not set, use  LoadReview or ExportReview to create new instance of the type")
	}

	params := lib.GetStateParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "Reviewer",
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

	result := new(lib.Reference[account])
	err = json.Unmarshal(out.Payload, result)
	if err != nil {
		return "", err
	}

	return result.Id(), err
}

func (s review) SetReviewer(newValue string) error {
	if s.id == "" {
		return errors.New("id of the type not set, use LoadReview or ExportReview to create new instance of the type")
	}

	params := lib.SetFieldParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "Reviewer",
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

func (s review) GetText() (string, error) {
	if s.id == "" {
		return *new(string), errors.New("id of the type not set, use  LoadReview or ExportReview to create new instance of the type")
	}

	params := lib.GetStateParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "Text",
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

func (s review) SetText(newValue string) error {
	if s.id == "" {
		return errors.New("id of the type not set, use LoadReview or ExportReview to create new instance of the type")
	}

	params := lib.SetFieldParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "Text",
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

func (s review) GetDownvotedBy() (map[string]struct{}, error) {
	if s.id == "" {
		return *new(map[string]struct{}), errors.New("id of the type not set, use  LoadReview or ExportReview to create new instance of the type")
	}

	params := lib.GetStateParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "DownvotedBy",
	}
	jsonParam, err := json.Marshal(params)
	if err != nil {
		return *new(map[string]struct{}), err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("GetField"), Payload: jsonParam})
	if _err != nil {
		return *new(map[string]struct{}), _err
	}
	if out.FunctionError != nil {
		return *new(map[string]struct{}), fmt.Errorf(string(out.Payload[:]))
	}

	result := new(map[string]struct{})
	err = json.Unmarshal(out.Payload, result)
	if err != nil {
		return *new(map[string]struct{}), err
	}
	return *result, err

}

func (s review) GetUpvotedBy() (map[string]struct{}, error) {
	if s.id == "" {
		return *new(map[string]struct{}), errors.New("id of the type not set, use  LoadReview or ExportReview to create new instance of the type")
	}

	params := lib.GetStateParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "UpvotedBy",
	}
	jsonParam, err := json.Marshal(params)
	if err != nil {
		return *new(map[string]struct{}), err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("GetField"), Payload: jsonParam})
	if _err != nil {
		return *new(map[string]struct{}), _err
	}
	if out.FunctionError != nil {
		return *new(map[string]struct{}), fmt.Errorf(string(out.Payload[:]))
	}

	result := new(map[string]struct{})
	err = json.Unmarshal(out.Payload, result)
	if err != nil {
		return *new(map[string]struct{}), err
	}
	return *result, err

}

func (s review) GetMapField() (map[string]string, error) {
	if s.id == "" {
		return *new(map[string]string), errors.New("id of the type not set, use  LoadReview or ExportReview to create new instance of the type")
	}

	params := lib.GetStateParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "MapField",
	}
	jsonParam, err := json.Marshal(params)
	if err != nil {
		return *new(map[string]string), err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("GetField"), Payload: jsonParam})
	if _err != nil {
		return *new(map[string]string), _err
	}
	if out.FunctionError != nil {
		return *new(map[string]string), fmt.Errorf(string(out.Payload[:]))
	}

	result := new(map[string]string)
	err = json.Unmarshal(out.Payload, result)
	if err != nil {
		return *new(map[string]string), err
	}
	return *result, err

}

func (s review) SetMapField(newValue map[string]string) error {
	if s.id == "" {
		return errors.New("id of the type not set, use LoadReview or ExportReview to create new instance of the type")
	}

	params := lib.SetFieldParam{
		Id:        s.GetId(),
		TypeName:  s.GetTypeName(),
		FieldName: "MapField",
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

func (m review) Downvote(input account) (int, error) {
	if m.id == "" {
		return *new(int), errors.New("id of the type not set, use  LoadReview or ExportReview to create new instance of the type")
	}

	params := new(lib.HandlerParameters)
	params.Id = m.id
	params.Parameter = input

	jsonParam, err := json.Marshal(params)
	if err != nil {
		return *new(int), err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("ReviewDownvote"), Payload: jsonParam})
	if _err != nil {
		return *new(int), _err
	}
	if out.FunctionError != nil {
		return *new(int), fmt.Errorf(string(out.Payload[:]))
	}

	result := new(int)

	_err = json.Unmarshal(out.Payload, result)
	if _err != nil {
		return *new(int), err
	}

	return *result, _err
}

func (m review) Upvote(input account) (int, error) {
	if m.id == "" {
		return *new(int), errors.New("id of the type not set, use  LoadReview or ExportReview to create new instance of the type")
	}

	params := new(lib.HandlerParameters)
	params.Id = m.id
	params.Parameter = input

	jsonParam, err := json.Marshal(params)
	if err != nil {
		return *new(int), err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("ReviewUpvote"), Payload: jsonParam})
	if _err != nil {
		return *new(int), _err
	}
	if out.FunctionError != nil {
		return *new(int), fmt.Errorf(string(out.Payload[:]))
	}

	result := new(int)

	_err = json.Unmarshal(out.Payload, result)
	if _err != nil {
		return *new(int), err
	}

	return *result, _err
}

func (r review) GetStub() (ReviewStub, error) {
	if r.id == "" {
		return *new(ReviewStub), errors.New("id of the type not set, use  LoadReview or ExportReview to create new instance of the type")
	}

	params := lib.GetStateParam{
		Id:       r.GetId(),
		TypeName: r.GetTypeName(),
		GetStub:  true,
	}
	jsonParam, err := json.Marshal(params)
	if err != nil {
		return *new(ReviewStub), err
	}

	out, _err := LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("GetState"), Payload: jsonParam})
	if _err != nil {
		return *new(ReviewStub), _err
	}
	if out.FunctionError != nil {
		return *new(ReviewStub), fmt.Errorf(string(out.Payload[:]))
	}

	result := new(ReviewStub)
	err = json.Unmarshal(out.Payload, result)
	if err != nil {
		return *new(ReviewStub), err
	}
	return *result, err
}

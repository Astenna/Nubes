package lib

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/google/uuid"
)

func Insert(objToInsert Nobject) (string, error) {
	var attributeVals, err = dynamodbattribute.MarshalMap(objToInsert)
	if err != nil {
		return "", err
	}

	var newId string
	if custom, ok := objToInsert.(CustomId); ok {
		if newId = custom.GetId(); newId == "" {
			return "", errors.New("id field empty. It must be set when using non-default id field")
		}
	} else {
		newId = uuid.New().String()
		attr := attributeVals["Id"].SetS(newId)
		// without this, dynamodb throws error because more than
		// one of the supported datatypes is set to not nil
		attr.NULL = nil
	}

	input := &dynamodb.PutItemInput{
		Item:      attributeVals,
		TableName: aws.String(objToInsert.GetTypeName()),
	}

	_, err = DBClient.PutItem(input)
	if err != nil {
		return "", err
	}

	return newId, nil
}

func Upsert(objToInsert Nobject, id string) error {
	var attributeVals, err = dynamodbattribute.MarshalMap(objToInsert)
	if err != nil {
		return err
	}

	attr := attributeVals["Id"].SetS(id)
	// without this, dynamodb throws error because more than
	// one of the supported datatypes is set to not nil
	attr.NULL = nil

	input := &dynamodb.PutItemInput{
		Item:      attributeVals,
		TableName: aws.String(objToInsert.GetTypeName()),
	}

	_, err = DBClient.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

func Delete[T Nobject](id string) error {
	if id == "" {
		return fmt.Errorf("missing id of object to delete")
	}

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String((*new(T)).GetTypeName()),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(id),
			},
		},
	}

	_, err := DBClient.DeleteItem(input)
	return err
}

func Get[T Nobject](id string, projections ...string) (*T, error) {
	if id == "" {
		return nil, fmt.Errorf("missing id of object to get")
	}

	var projectionExpr expression.Expression
	if len(projections) > 0 {
		expr, err := expression.NewBuilder().WithProjection(getProjection(projections)).Build()
		projectionExpr = expr

		if err != nil {
			return nil, fmt.Errorf("error occurred when creating projection: %w", err)
		}
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String((*new(T)).GetTypeName()),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(id),
			},
		},
		ExpressionAttributeNames: projectionExpr.Names(),
		ProjectionExpression:     projectionExpr.Projection(),
	}

	item, err := DBClient.GetItem(input)
	if err != nil {
		return nil, err
	}

	var parsedItem = new(T)
	// in order not to leave the id empty
	// in case the projection was used, but ID was not requested
	if item.Item != nil {
		item.Item["Id"] = &dynamodb.AttributeValue{
			S: aws.String(id),
		}

		err = dynamodbattribute.UnmarshalMap(item.Item, parsedItem)
		return parsedItem, err
	}

	return nil, err
}

func Update[T Nobject](values aws.JSONValue) error {
	if len(values) == 0 {
		return fmt.Errorf("no values specified for update")
	}
	if values["id"] == "" {
		return fmt.Errorf("missing id of object to update")
	}

	update := expression.UpdateBuilder{}
	for k, v := range values {
		update = update.Set(expression.Name(k), expression.Value(v))
	}
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return fmt.Errorf("error occurred when building dynamodb update expression %w", err)
	}

	_, err = DBClient.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String((*new(T)).GetTypeName()),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(values["Id"].(string)),
			},
		},
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})

	return err
}

type GetFieldParam struct {
	FieldName string
	TypeName  string
}

func GetField(param HandlerParameters) (interface{}, error) {
	getFielParam, castErr := param.Parameter.(GetFieldParam)
	if castErr {
		return *new(interface{}), fmt.Errorf("missing GetFieldParam")
	}
	if param.Id == "" {
		return *new(interface{}), fmt.Errorf("missing id of object's field  to get")
	}
	if getFielParam.FieldName == "" {
		return *new(interface{}), fmt.Errorf("missing field name of object's field to get")
	}
	if getFielParam.TypeName == "" {
		return *new(interface{}), fmt.Errorf("missing type name of object's field to get")
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(getFielParam.TypeName),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(param.Id),
			},
		},
		ProjectionExpression: &getFielParam.FieldName,
	}

	item, err := DBClient.GetItem(input)
	if err != nil {
		return *new(interface{}), err
	}

	if item.Item != nil {
		var parsedItem interface{}
		err = dynamodbattribute.Unmarshal(item.Item[getFielParam.FieldName], &parsedItem)
		return parsedItem, err

	}

	return nil, err
}

func getProjection(names []string) expression.ProjectionBuilder {
	if len(names) == 0 {
		return *new(expression.ProjectionBuilder)
	}

	var builder expression.ProjectionBuilder
	for _, name := range names {
		builder = builder.AddNames(expression.Name(name))
	}
	return builder
}

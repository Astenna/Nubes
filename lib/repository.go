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

func Get[T Nobject](id string) (*T, error) {
	if id == "" {
		return nil, fmt.Errorf("missing id of object to get")
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String((*new(T)).GetTypeName()),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(id),
			},
		},
	}

	item, err := DBClient.GetItem(input)
	if err != nil {
		return nil, err
	}

	var parsedItem = new(T)
	err = dynamodbattribute.UnmarshalMap(item.Item, parsedItem)
	return parsedItem, err
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

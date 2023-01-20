package lib

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

func Load[T Nobject](id string) (*T, error) {
	instance := *new(T)
	dbIdAttributeName := "Id"

	item, err := DBClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(instance.GetTypeName()),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(id),
			},
		},
		ProjectionExpression: &dbIdAttributeName,
	})

	if err != nil {
		return nil, err
	}

	var parsedItem = new(T)
	if item.Item != nil {
		err = dynamodbattribute.UnmarshalMap(item.Item, parsedItem)
		return parsedItem, err
	}

	return nil, nil
}

func Export[T Nobject](objToInsert Nobject) (*T, error) {
	var attributeVals, err = dynamodbattribute.MarshalMap(objToInsert)
	if err != nil {
		return new(T), err
	}

	var newId string
	if custom, ok := objToInsert.(CustomId); ok {
		if newId = custom.GetId(); newId == "" {
			return new(T), errors.New("id field empty. It must be set when using non-default id field")
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
		return new(T), err
	}

	// Unmarshal back the input to get the object with ID set
	// TODO: verify if it's more efficient to
	// Unmarshal or to use the reflection to set the ID
	var parsedItem = new(T)
	err = dynamodbattribute.UnmarshalMap(attributeVals, parsedItem)
	return parsedItem, err
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

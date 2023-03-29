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
	instance := new(T)
	instanceTypeName := (*instance).GetTypeName()
	dbIdAttributeName := "Id"

	item, err := DBClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(instanceTypeName),
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

	if item.Item != nil {
		err = dynamodbattribute.UnmarshalMap(item.Item, instance)
		invokeInitOnNobjectType(instance)
		return instance, err
	}

	return nil, NotFoundError{TypeName: instanceTypeName, Ids: []string{id}}
}

func LoadBatch[T Nobject](ids []string) ([]*T, error) {
	if ids == nil {
		return nil, fmt.Errorf("missing ids of objects to get")
	}

	keysToRetrieve := make([]map[string]*dynamodb.AttributeValue, len(ids))
	for i, id := range ids {
		keysToRetrieve[i] = map[string]*dynamodb.AttributeValue{"Id": {
			S: aws.String(id),
		}}
	}

	tableName := (*new(T)).GetTypeName()
	input := &dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			tableName: {
				Keys:                 keysToRetrieve,
				ProjectionExpression: aws.String("Id"),
			},
		},
	}

	items, err := DBClient.BatchGetItem(input)
	if err != nil {
		return nil, err
	}

	var parsedItem = new([]*T)
	if items.Responses[tableName] != nil {
		err = dynamodbattribute.UnmarshalListOfMaps(items.Responses[tableName], parsedItem)

		for _, item := range *parsedItem {
			invokeInitOnNobjectType(item)
		}
		return *parsedItem, err
	}

	return nil, err
}

func Export[T Nobject](objToInsert Nobject) (*T, error) {
	var attributeVals, err = dynamodbattribute.MarshalMap(objToInsert)
	if err != nil {
		return new(T), err
	}

	var newId string
	var conditionExpression *string
	if custom, ok := objToInsert.(CustomId); ok {
		if newId = custom.GetId(); newId == "" {
			return new(T), errors.New("id field empty. It must be set when using non-default id field")
		}
		newId = custom.GetId()
		conditionExpression = aws.String("attribute_not_exists(Id)")
	} else {
		newId = uuid.New().String()
		attr := attributeVals["Id"].SetS(newId)
		// without this, dynamodb throws error because more than
		// one of the supported datatypes is set to not nil
		attr.NULL = nil
	}

	input := &dynamodb.PutItemInput{
		Item:                attributeVals,
		TableName:           aws.String(objToInsert.GetTypeName()),
		ConditionExpression: conditionExpression,
	}

	_, err = DBClient.PutItem(input)
	if err != nil {
		if _, ok := err.(*dynamodb.ConditionalCheckFailedException); ok {
			return nil, fmt.Errorf("instance of %s with id: %s already exists. Use lib.Load(id) to work on existing instances", objToInsert.GetTypeName(), newId)
		}
		return nil, err
	}

	// Unmarshal back the input to get the object with ID set
	// TODO: verify if it's more efficient to
	// Unmarshal or to use the reflection to set the ID
	var parsedItem = new(T)
	err = dynamodbattribute.UnmarshalMap(attributeVals, parsedItem)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal object returned from DB")
	}

	invokeInitOnNobjectType(parsedItem)
	return parsedItem, err
}

func invokeInitOnNobjectType[T Nobject](item *T) {
	castedToInterface := any(item)
	if initInterface, ok := castedToInterface.(nobjectInit); ok {
		initInterface.Init()
	}
}

func Delete[T Nobject](id string) error {
	if id == "" {
		return fmt.Errorf("missing id of object to delete")
	}
	typeName := (*new(T)).GetTypeName()

	return DeleteWithTypeNameAsArg(id, typeName)
}

func DeleteWithTypeNameAsArg(id, typeName string) error {
	if id == "" {
		return fmt.Errorf("missing id of object to delete")
	}
	if typeName == "" {
		return fmt.Errorf("missing typeName of object to delete")
	}

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(typeName),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(id),
			},
		},
		ConditionExpression: aws.String("attribute_exists(Id)"),
	}

	_, err := DBClient.DeleteItem(input)
	if _, ok := err.(*dynamodb.ConditionalCheckFailedException); ok {
		return fmt.Errorf("delete failed. Instance of %s with id: %s not found", typeName, id)
	}
	return err
}

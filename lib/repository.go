package lib

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

func Insert(objToInsert Object) error {
	var attributeVals, err = dynamodbattribute.MarshalMap(objToInsert)
	if err != nil {
		return err
	}

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

func Delete[T Object](id string) error {
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

func Get[T Object](id string) (*T, error) {
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

func Update[T Object](values aws.JSONValue) error {
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

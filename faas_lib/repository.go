package faas_lib

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"strconv"
)

func Create(objToInsert Object) error {
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

func Delete[T Object](id int) error {

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String((*new(T)).GetTypeName()),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				N: aws.String(strconv.Itoa(id)),
			},
		},
	}

	_, err := DBClient.DeleteItem(input)
	return err
}

func Get[T Object](id int) (*T, error) {

	input := &dynamodb.GetItemInput{
		TableName: aws.String((*new(T)).GetTypeName()),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				N: aws.String(strconv.Itoa(id)),
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

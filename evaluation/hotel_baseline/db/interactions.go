package db

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

func GetSingleItemByPartitonKey[T any](tableName, keyAttribute, keyValue string) (T, error) {
	var result = new(T)
	if keyAttribute == "" {
		return *result, fmt.Errorf("missing keyAttribute of element to get")
	}
	if keyValue == "" {
		return *result, fmt.Errorf("missing keyValue of element to get")
	}
	if tableName == "" {
		return *result, fmt.Errorf("missing tablename of element to get")
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			keyAttribute: {
				S: aws.String(keyValue),
			},
		},
	}

	item, err := DbClient.GetItem(input)
	if err != nil {
		return *result, err
	}

	if item.Item != nil {
		err = dynamodbattribute.UnmarshalMap(item.Item, result)
		return *result, err
	}

	return *result, fmt.Errorf("element with key: %s not found in table", keyValue, tableName)
}

func GetItemsByPartitonKey[T any](tableName, partitionAttribute, partitionValue string) ([]T, error) {
	if partitionAttribute == "" {
		return nil, fmt.Errorf("missing partition key attribute name of elements to get")
	}
	if partitionValue == "" {
		return nil, fmt.Errorf("missing partition key value of elements to get")
	}
	if tableName == "" {
		return nil, fmt.Errorf("missing tablename of element to get")
	}

	keyCondition := expression.Key(partitionAttribute).Equal(expression.Value(partitionValue))
	expr, errExpression := expression.NewBuilder().
		WithKeyCondition(keyCondition).
		Build()
	if errExpression != nil {
		fmt.Println("error: creating dynamoDB expression ", errExpression)
		return nil, errExpression
	}
	input := &dynamodb.QueryInput{
		TableName:                 aws.String(tableName),
		ExpressionAttributeNames:  expr.Names(),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeValues: expr.Values(),
	}

	items, err := DbClient.Query(input)
	if err != nil {
		return nil, err
	}

	parsedItems := make([]T, len(items.Items))
	err = dynamodbattribute.UnmarshalListOfMaps(items.Items, parsedItems)
	return parsedItems, err
}

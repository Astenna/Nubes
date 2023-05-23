package db

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

func GetItemsByPartitonKeys[T any](tableName, partitionAttributeName string, partitionAttributeValues []string) ([]T, error) {
	if partitionAttributeName == "" {
		return nil, fmt.Errorf("missing partition key attribute name of elements to get")
	}
	if len(partitionAttributeValues) == 0 {
		return nil, fmt.Errorf("missing partitionAttributeValue of elements to get")
	}
	if tableName == "" {
		return nil, fmt.Errorf("missing tablename of elements to get")
	}

	keysToRetrieve := make([]map[string]*dynamodb.AttributeValue, len(partitionAttributeValues))
	for i, id := range partitionAttributeValues {
		keysToRetrieve[i] = map[string]*dynamodb.AttributeValue{partitionAttributeName: {
			S: aws.String(id),
		}}
	}

	input := &dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			tableName: {
				Keys: keysToRetrieve,
			},
		},
	}

	items, err := DbClient.BatchGetItem(input)
	if err != nil {
		return nil, err
	}

	var parsedItem = new([]T)
	if items.Responses[tableName] != nil {

		err = dynamodbattribute.UnmarshalListOfMaps(items.Responses[tableName], parsedItem)
		return *parsedItem, err
	}

	return nil, err
}

func GetHotelIdsInCityByIndex(cityName string) ([]string, error) {
	keyCondition := expression.Key("CityName").Equal(expression.Value(cityName))
	expr, errExpression := expression.NewBuilder().
		WithKeyCondition(keyCondition).
		Build()
	if errExpression != nil {
		fmt.Println("error: creating dynamoDB expression ", errExpression)
		return nil, errExpression
	}

	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String(HotelTable),
		IndexName:                 aws.String(HotelTableIndex),
		ExpressionAttributeNames:  expr.Names(),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeValues: expr.Values(),
	}

	items, err := DbClient.Query(queryInput)
	if err != nil {
		return nil, err
	}

	outputIds := make([]string, len(items.Items))
	for index, attr := range items.Items {
		outputIds[index] = *attr["HotelName"].S
	}
	return outputIds, nil
}

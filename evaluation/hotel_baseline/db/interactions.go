package db

import (
	"fmt"
	"time"

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
	if len(items.Items) > 0 {
		err = dynamodbattribute.UnmarshalListOfMaps(items.Items, &parsedItems)
	}
	return parsedItems, err
}

type GetItemBySortKey struct {
	PkName    string
	PkValue   string
	SkName    string
	SkValue   time.Time
	TableName string
}

func (g GetItemBySortKey) Verify() error {
	if g.PkName == "" {
		return fmt.Errorf("missing partition key name")
	}
	if g.PkValue == "" {
		return fmt.Errorf("missing partition key value")
	}
	if g.SkName == "" {
		return fmt.Errorf("missing sort key name")
	}
	if g.SkValue.IsZero() {
		return fmt.Errorf("missing sort key value")
	}
	if g.TableName == "" {
		return fmt.Errorf("missing table name value")
	}
	return nil
}

func GetItemBeforeSortKey[T any](param GetItemBySortKey) (T, error) {
	if err := param.Verify(); err != nil {
		return *new(T), err
	}

	pkKeyCondition := expression.Key(param.PkName).
		Equal(expression.Value(param.PkValue))
	skCondition := expression.KeyLessThanEqual(expression.Key(param.SkName), expression.Value(param.SkValue))
	pkAndSkCondition := expression.KeyAnd(pkKeyCondition, skCondition)

	expr, errExpression := expression.NewBuilder().
		WithKeyCondition(pkAndSkCondition).
		Build()
	if errExpression != nil {
		fmt.Println("error: creating dynamoDB expression ", errExpression)
		return *new(T), errExpression
	}

	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String(param.TableName),
		ExpressionAttributeNames:  expr.Names(),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeValues: expr.Values(),
		Limit:                     aws.Int64(1),
		ScanIndexForward:          aws.Bool(false),
	}

	items, err := DbClient.Query(queryInput)
	if err != nil {
		return *new(T), err
	}

	result := new(T)
	err = dynamodbattribute.UnmarshalMap(items.Items[0], result)
	return *result, err
}

func GetItemAfterSortKey[T any](param GetItemBySortKey) (T, error) {
	if err := param.Verify(); err != nil {
		return *new(T), err
	}

	pkKeyCondition := expression.Key(param.PkName).
		Equal(expression.Value(param.PkValue))
	skCondition := expression.KeyGreaterThanEqual(expression.Key(param.SkName), expression.Value(param.SkValue))
	pkAndSkCondition := expression.KeyAnd(pkKeyCondition, skCondition)

	expr, errExpression := expression.NewBuilder().
		WithKeyCondition(pkAndSkCondition).
		Build()
	if errExpression != nil {
		fmt.Println("error: creating dynamoDB expression ", errExpression)
		return *new(T), errExpression
	}

	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String(param.TableName),
		ExpressionAttributeNames:  expr.Names(),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeValues: expr.Values(),
		Limit:                     aws.Int64(1),
	}

	items, err := DbClient.Query(queryInput)
	if err != nil {
		return *new(T), err
	}

	result := new(T)
	err = dynamodbattribute.UnmarshalMap(items.Items[0], result)
	return *result, err
}

func Insert(toInsert any, tableName string) error {
	var attributeVals, err = dynamodbattribute.MarshalMap(toInsert)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      attributeVals,
		TableName: aws.String(tableName),
	}

	_, err = DbClient.PutItem(input)
	return err
}

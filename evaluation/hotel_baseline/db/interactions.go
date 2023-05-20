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

func GetItemBeforeSortKey[T any](param GetItemBySortKey) (*T, error) {
	if err := param.Verify(); err != nil {
		return nil, err
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
		return nil, errExpression
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
		return nil, err
	}

	if len(items.Items) > 0 {
		result := new(T)
		err = dynamodbattribute.UnmarshalMap(items.Items[0], result)
		return result, err
	}
	return nil, err
}

func GetItemAfterSortKey[T any](param GetItemBySortKey) (*T, error) {
	if err := param.Verify(); err != nil {
		return nil, err
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
		return nil, errExpression
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
		return nil, err
	}

	if len(items.Items) > 0 {
		result := new(T)
		err = dynamodbattribute.UnmarshalMap(items.Items[0], result)
		return result, err
	}
	return nil, err
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

func DeleteByPartitionKey(pkValue, pkName, tableName string) error {
	if pkValue == "" {
		return fmt.Errorf("missing pkValue of item to delete")
	}
	if pkName == "" {
		return fmt.Errorf("missing pkName of item to delete")
	}

	if tableName == "" {
		return fmt.Errorf("missing tableName of user to delete")
	}

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			pkName: {
				S: aws.String(pkValue),
			},
		},
		ConditionExpression: aws.String("attribute_exists(" + pkName + ")"),
	}

	_, err := DbClient.DeleteItem(input)
	if _, ok := err.(*dynamodb.ConditionalCheckFailedException); ok {
		return fmt.Errorf("delete failed. Instance of item with pk: %s not found in table %s", pkValue, tableName)
	}
	return err
}

func SetHotelField[T any](cityName, hotelName, fieldName string, fieldValue T) error {
	if fieldName == "" {
		return fmt.Errorf("missing fieldName")
	}
	if hotelName == "" {
		return fmt.Errorf("missing keyValue")
	}
	if cityName == "" {
		return fmt.Errorf("missing keyValue")
	}

	update := expression.UpdateBuilder{}
	update = update.Set(expression.Name(fieldName), expression.Value(fieldValue))
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return fmt.Errorf("error occurred when building dynamodb update expression %w", err)
	}

	_, err = DbClient.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(HotelTable),
		Key: map[string]*dynamodb.AttributeValue{
			"CityName": {
				S: aws.String(cityName),
			},
			"HotelName": {
				S: aws.String(hotelName),
			},
		},
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})

	return err
}

type UserReservationsJoinTableEntry struct {
	UserId          string
	CityHotelRoomId string
	DateIn          time.Time
}

func GetUserReservationsCompositeKeys(userEmail string) ([]CompositeKey, error) {
	if userEmail == "" {
		return nil, fmt.Errorf("missing userEmail")
	}

	keyCondition := expression.Key("UserId").Equal(expression.Value(userEmail))
	projection := expression.NamesList(expression.Name("CityHotelRoomId"), expression.Name("DateIn"))
	expr, errExpression := expression.NewBuilder().
		WithKeyCondition(keyCondition).
		WithProjection(projection).
		Build()
	if errExpression != nil {
		fmt.Println("error: creating dynamoDB expression ", errExpression)
		return nil, errExpression
	}
	input := &dynamodb.QueryInput{
		TableName:                 aws.String(UserResevationsJoinTable),
		ExpressionAttributeNames:  expr.Names(),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeValues: expr.Values(),
		ProjectionExpression:      expr.Projection(),
	}

	items, err := DbClient.Query(input)
	if err != nil {
		return nil, err
	}

	var compositeKeys = make([]CompositeKey, len(items.Items))
	for i, item := range items.Items {
		compositeKeys[i] = CompositeKey{
			PartitionKey: *item["CityHotelRoomId"].S,
			SortKey:      *item["DateIn"].S,
		}
	}
	return compositeKeys, nil
}

type CompositeKey struct {
	PartitionKey string
	SortKey      string
}

func GetBatchItemsUsingCompositeKeys[T any](keys []CompositeKey, tableName, pkAttributeName, skAttributeName string) ([]T, error) {

	keysToRetrieve := make([]map[string]*dynamodb.AttributeValue, len(keys))
	for i, key := range keys {
		keysToRetrieve[i] = map[string]*dynamodb.AttributeValue{
			pkAttributeName: {
				S: aws.String(key.PartitionKey),
			},
			skAttributeName: {
				S: aws.String(key.SortKey),
			},
		}
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

	var parsedItem []T
	if items.Responses[tableName] != nil {

		err = dynamodbattribute.UnmarshalListOfMaps(items.Responses[tableName], &parsedItem)
		return parsedItem, err
	}

	return nil, err
}

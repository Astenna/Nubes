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

func GetObjectState[T Nobject](id string) (*T, error) {
	if id == "" {
		return nil, fmt.Errorf("missing id of object to get")
	}

	typeName := (*new(T)).GetTypeName()
	input := &dynamodb.GetItemInput{
		TableName: aws.String(typeName),
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
	if item.Item != nil {
		err = dynamodbattribute.UnmarshalMap(item.Item, parsedItem)
		return parsedItem, err
	}

	return nil, fmt.Errorf("%s with id: %s not found", typeName, id)
}

func GetByIndex(param QueryByIndexParam) ([]string, error) {
	if err := param.Validate(); err != nil {
		return nil, err
	}

	keyCondition := expression.Key(param.KeyAttributeName).Equal(expression.Value(param.KeyAttributeValue))
	expr, errExpression := expression.NewBuilder().
		WithKeyCondition(keyCondition).
		WithProjection(getProjection([]string{param.OutputAttributeName})).
		Build()
	if errExpression != nil {
		fmt.Println("error: creating dynamoDB expression ", errExpression)
		return nil, errExpression
	}

	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String(param.TableName),
		IndexName:                 aws.String(param.IndexName),
		ExpressionAttributeNames:  expr.Names(),
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
		ExpressionAttributeValues: expr.Values(),
	}

	items, err := DBClient.Query(queryInput)
	if err != nil {
		return nil, err
	}

	outputIds := make([]string, len(items.Items))
	for index, attr := range items.Items {
		outputIds[index] = *attr[param.OutputAttributeName].S
	}
	return outputIds, nil
}

func GetSortKeysByPartitionKey(q QueryByPartitionKeyParam) ([]string, error) {
	if err := q.Validate(); err != nil {
		return nil, err
	}

	keyCondition := expression.Key(q.PartitionAttributeName).Equal(expression.Value(q.PatritionAttributeValue))
	expr, errExpression := expression.NewBuilder().
		WithKeyCondition(keyCondition).
		WithProjection(getProjection([]string{q.OutputAttributeName})).
		Build()
	if errExpression != nil {
		fmt.Println("error: creating dynamoDB expression ", errExpression)
		return nil, errExpression
	}
	input := &dynamodb.QueryInput{
		TableName:                 aws.String(q.TableName),
		ExpressionAttributeNames:  expr.Names(),
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
		ExpressionAttributeValues: expr.Values(),
	}

	items, err := DBClient.Query(input)
	if err != nil {
		return nil, err
	}

	outputIds := make([]string, len(items.Items))
	for index, attr := range items.Items {
		outputIds[index] = *attr[q.OutputAttributeName].S
	}

	return outputIds, err
}

func GetBatch[T Nobject](ids []string) (*[]T, error) {
	if ids == nil {
		return nil, fmt.Errorf("missing id of object to get")
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
				Keys: keysToRetrieve,
			},
		},
	}

	items, err := DBClient.BatchGetItem(input)
	if err != nil {
		return nil, err
	}

	var parsedItem = new([]T)
	if items.Responses[tableName] != nil {

		err = dynamodbattribute.UnmarshalListOfMaps(items.Responses[tableName], parsedItem)
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

func GetField(param GetFieldParam) (interface{}, error) {
	if err := param.Validate(); err != nil {
		return nil, err
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(param.TypeName),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(param.Id),
			},
		},
		ProjectionExpression: &param.FieldName,
	}

	item, err := DBClient.GetItem(input)
	if err != nil {
		return *new(interface{}), err
	}

	if item.Item != nil {
		var parsedItem interface{}
		err = dynamodbattribute.Unmarshal(item.Item[param.FieldName], &parsedItem)
		return parsedItem, err
	}

	return nil, err
}

func SetField(param SetFieldParam) error {
	if err := param.Validate(); err != nil {
		return err
	}

	update := expression.UpdateBuilder{}
	update = update.Set(expression.Name(param.FieldName), expression.Value(param.Value))
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return fmt.Errorf("error occurred when building dynamodb update expression %w", err)
	}

	_, err = DBClient.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(param.TypeName),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(param.Id),
			},
		},
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})

	return err
}

func IsInstanceAlreadyCreated(param IsInstanceAlreadyCreatedParam) (bool, error) {
	item, err := DBClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(param.TypeName),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(param.Id),
			},
		},
		ProjectionExpression: aws.String("Id"),
	})

	if err != nil {
		return false, err
	}

	if item.Item != nil {
		return true, nil
	}

	return false, nil
}

func isInstanceWithDefaultIdAlreadyCreated(instance Nobject) (bool, error) {
	var attributeVals, err = dynamodbattribute.MarshalMap(instance)
	if err != nil {
		return false, err
	}

	return *(attributeVals["Id"].S) == "", nil
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

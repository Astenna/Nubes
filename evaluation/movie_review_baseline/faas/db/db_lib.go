package db

import (
	"errors"
	"fmt"

	"github.com/Astenna/Nubes/evaluation/movie_review_baseline/faas/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/google/uuid"
)

type CustomId interface {
	GetId() string
}

func Insert(objToInsert any, tableName string) (string, error) {
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
		attr.NULL = nil
	}

	input := &dynamodb.PutItemInput{
		Item:      attributeVals,
		TableName: aws.String(tableName),
	}

	_, err = DBClient.PutItem(input)
	if err != nil {
		return "", err
	}
	return newId, nil
}

func Delete(id, tableName string) error {
	if id == "" {
		return fmt.Errorf("missing id of object to delete")
	}
	if tableName == "" {
		return fmt.Errorf("missing tableName of object to delete")
	}

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(id),
			},
		},
		ConditionExpression: aws.String("attribute_exists(Id)"),
	}

	_, err := DBClient.DeleteItem(input)
	if _, ok := err.(*dynamodb.ConditionalCheckFailedException); ok {
		return fmt.Errorf("delete failed. Instance of %s with id: %s not found", tableName, id)
	}
	return err
}

func GetById[T any](id, tablename string) (*T, error) {
	if id == "" {
		return nil, fmt.Errorf("missing id of object to get")
	}
	if tablename == "" {
		return nil, fmt.Errorf("missing tablename of object to get")
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(tablename),
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

	return nil, fmt.Errorf("%s with id: %s not found", tablename, id)
}

func GetMoviesByCategory(categoryName string) ([]models.CategoryListItem, error) {
	if categoryName == "" {
		return nil, fmt.Errorf("missing categoryName")
	}

	keyCondition := expression.Key("Category").Equal(expression.Value(categoryName))
	expr, errExpression := expression.NewBuilder().
		WithKeyCondition(keyCondition).
		Build()
	if errExpression != nil {
		fmt.Println("error: creating dynamoDB expression ", errExpression)
		return nil, errExpression
	}

	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String("Movie"),
		IndexName:                 aws.String("MovieCategory"),
		ExpressionAttributeNames:  expr.Names(),
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
		ExpressionAttributeValues: expr.Values(),
	}

	items, err := DBClient.Query(queryInput)
	if err != nil {
		return nil, err
	}

	result := make([]models.CategoryListItem, len(items.Items))
	for index, attr := range items.Items {
		result[index].Id = *attr["Id"].S
		// TODO: Title is not here since it is not contained in the index!
		result[index].Title = *attr["Title"].S
	}
	return result, nil
}

func GetMovieReviews(movieId string) ([]models.Review, error) {
	if movieId == "" {
		return nil, fmt.Errorf("missing movieId")
	}

	keyCondition := expression.Key(movieId).Equal(expression.Value("MovieId"))
	expr, errExpression := expression.NewBuilder().
		WithKeyCondition(keyCondition).
		Build()
	if errExpression != nil {
		fmt.Println("error: creating dynamoDB expression ", errExpression)
		return nil, errExpression
	}

	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String("Movie"),
		IndexName:                 aws.String("ReviewMovie"),
		ExpressionAttributeNames:  expr.Names(),
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
		ExpressionAttributeValues: expr.Values(),
	}

	items, err := DBClient.Query(queryInput)
	if err != nil {
		return nil, err
	}

	result := make([]models.Review, len(items.Items))
	err = dynamodbattribute.UnmarshalListOfMaps(items.Items, result)
	return result, nil
}

func Upsert(toUpsert any, id, tableName string) error {

	if id == "" {
		return fmt.Errorf("missing id")
	}
	if tableName == "" {
		return fmt.Errorf("missing tableName")
	}
	if toUpsert == nil {
		return fmt.Errorf("missing object toUpsert")
	}

	var attributeVals, err = dynamodbattribute.MarshalMap(toUpsert)
	if err != nil {
		return err
	}

	attr := attributeVals["Id"].SetS(id)
	// without this, dynamodb throws error because more than
	// one of the supported datatypes is set to not nil
	attr.NULL = nil

	input := &dynamodb.PutItemInput{
		Item:      attributeVals,
		TableName: aws.String(tableName),
	}

	_, err = DBClient.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

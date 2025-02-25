package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var _session = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))

var dbClient = dynamodb.New(_session)

func insert[T any](toBeInserted T, tableName string) {
	var attributeVals, err = dynamodbattribute.MarshalMap(toBeInserted)
	if err != nil {
		fmt.Println(err)
		return
	}

	input := &dynamodb.PutItemInput{
		Item:      attributeVals,
		TableName: aws.String(tableName),
	}

	_, err = dbClient.PutItem(input)
	if err != nil {
		fmt.Println(err)
		return
	}
}

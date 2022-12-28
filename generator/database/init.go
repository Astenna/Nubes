package database

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func CreateTypeTables(isNobjectTypeMap map[string]bool) {
	var _session = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	var dblient = dynamodb.New(_session)

	for typeName, isNobjectType := range isNobjectTypeMap {
		if isNobjectType {
			_, err := dblient.CreateTable(&dynamodb.CreateTableInput{
				BillingMode: aws.String("PROVISIONED"),
				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(1),
					WriteCapacityUnits: aws.Int64(1),
				},
				AttributeDefinitions: []*dynamodb.AttributeDefinition{
					{
						AttributeName: aws.String("Id"),
						AttributeType: aws.String("S"),
					},
				},
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("Id"),
						KeyType:       aws.String("HASH"),
					},
				},
				TableName: aws.String(typeName),
			})

			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

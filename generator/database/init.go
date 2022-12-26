package database

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func CreateTypeTables(nobjects []string) {
	var _session = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	var dblient = dynamodb.New(_session)

	for _, nobjectType := range nobjects {
		_, err := dblient.CreateTable(&dynamodb.CreateTableInput{
			BillingMode: aws.String("PROVISIONED"),
			ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(5),
				WriteCapacityUnits: aws.Int64(5),
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
			TableName: aws.String(nobjectType),
		})

		if err != nil {
			fmt.Println(err)
		}
	}
}

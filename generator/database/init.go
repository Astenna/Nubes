package database

import (
	"fmt"

	"github.com/Astenna/Nubes/generator/parser"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func CreateTypeTables(parsedPackage parser.ParsedPackage) {
	var _session = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	var dblient = dynamodb.New(_session)

	for typeName, isNobjectType := range parsedPackage.IsNobjectInOrginalPackage {
		if isNobjectType {
			createTableInput := &dynamodb.CreateTableInput{
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
			}

			if typeIndexes, ok := parsedPackage.TypeAttributesIndexes[typeName]; ok {
				for _, attributeName := range typeIndexes {
					createTableInput.GlobalSecondaryIndexes = []*dynamodb.GlobalSecondaryIndex{
						{
							IndexName: aws.String(typeName + attributeName),
							KeySchema: []*dynamodb.KeySchemaElement{
								{
									AttributeName: aws.String(attributeName),
									KeyType:       aws.String("HASH"),
								},
							},
							Projection: &dynamodb.Projection{
								ProjectionType: aws.String("KEYS_ONLY"),
							},
							ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
								ReadCapacityUnits:  aws.Int64(1),
								WriteCapacityUnits: aws.Int64(1),
							},
						},
					}
					createTableInput.AttributeDefinitions = append(createTableInput.AttributeDefinitions,
						&dynamodb.AttributeDefinition{
							AttributeName: aws.String(attributeName),
							AttributeType: aws.String("S"),
						},
					)
				}
			}

			_, err := dblient.CreateTable(createTableInput)

			if err != nil {
				if _, ok := err.(*dynamodb.ResourceInUseException); ok {
					fmt.Println("Table for type: ", typeName, " already created")
					continue
				}
				fmt.Println(err)
			}
		}
	}
}

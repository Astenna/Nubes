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

	tableCreated := map[string]struct{}{}
	for _, typeManyToManyRelationship := range parsedPackage.ManyToManyRelationships {
		for _, relationship := range typeManyToManyRelationship {

			if _, exists := tableCreated[relationship.TableName]; !exists {
				joinTable := &dynamodb.CreateTableInput{
					BillingMode: aws.String("PROVISIONED"),
					ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
						ReadCapacityUnits:  aws.Int64(1),
						WriteCapacityUnits: aws.Int64(1),
					},
					AttributeDefinitions: []*dynamodb.AttributeDefinition{
						{
							AttributeName: aws.String(relationship.PartionKeyName),
							AttributeType: aws.String("S"),
						},
						{
							AttributeName: aws.String(relationship.SortKeyName),
							AttributeType: aws.String("S"),
						},
					},
					KeySchema: []*dynamodb.KeySchemaElement{
						{
							AttributeName: aws.String(relationship.PartionKeyName),
							KeyType:       aws.String("HASH"),
						},
						{
							AttributeName: aws.String(relationship.SortKeyName),
							KeyType:       aws.String("RANGE"),
						},
					},
					TableName: aws.String(relationship.TableName),
					GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
						{
							IndexName: aws.String(relationship.TableName + "Reversed"),
							KeySchema: []*dynamodb.KeySchemaElement{
								{
									AttributeName: aws.String(relationship.SortKeyName),
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
					},
				}

				_, err := dblient.CreateTable(joinTable)

				if err != nil {
					if _, ok := err.(*dynamodb.ResourceInUseException); ok {
						fmt.Println("Join table for many-to-many relationship: ", relationship.TableName, " already created")
						continue
					}
					fmt.Println(err)
				}
				tableCreated[relationship.TableName] = struct{}{}
			}
		}
	}
}

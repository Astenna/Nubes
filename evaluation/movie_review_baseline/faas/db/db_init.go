package db

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type IndexDefinition struct {
	Column    string
	IndexName string
}

func InitializeTables() {
	tableNames := []string{"Account", "Review", "Movie"}
	tableToIndex := map[string]IndexDefinition{"Review": {Column: "MovieId", IndexName: "ReviewMovie"},
		"Movie": {Column: "Category", IndexName: "MovieCategory"}}

	for _, tableName := range tableNames {
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
			TableName: aws.String(tableName),
		}

		if typeIndex, ok := tableToIndex[tableName]; ok {
			createTableInput.GlobalSecondaryIndexes = []*dynamodb.GlobalSecondaryIndex{
				{
					IndexName: aws.String(typeIndex.IndexName),
					KeySchema: []*dynamodb.KeySchemaElement{
						{
							AttributeName: aws.String(typeIndex.Column),
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
					AttributeName: aws.String(typeIndex.Column),
					AttributeType: aws.String("S"),
				},
			)
		}
		_, err := DBClient.CreateTable(createTableInput)

		if err != nil {
			if _, ok := err.(*dynamodb.ResourceInUseException); ok {
				fmt.Println("Table ", tableName, " already created")
				continue
			}
			fmt.Println(err)
		}
	}
}

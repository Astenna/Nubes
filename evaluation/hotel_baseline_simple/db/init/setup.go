package main

import (
	"fmt"

	"github.com/Astenna/Nubes/evaluation/hotel_baseline_simple/db"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type IndexDefinition struct {
	PartitionKeyColumn string
	SortKeyColumn      string
	IndexName          string
}

type TableDefinition struct {
	TableName    string
	PartitionKey string
	SortKey      string
	Indexes      []IndexDefinition
}

func InitializeTables() {

	tableDefinitions := []TableDefinition{
		{
			TableName:    db.HotelTable,
			PartitionKey: "HotelName",
			Indexes: []IndexDefinition{{
				PartitionKeyColumn: "CityName",
				IndexName:          db.HotelTableIndex,
			}},
		},
		{
			TableName:    db.CityTable,
			PartitionKey: "CityName",
		},
	}

	for _, tableDefinition := range tableDefinitions {
		createTableInput := &dynamodb.CreateTableInput{
			BillingMode: aws.String("PAY_PER_REQUEST"),
			AttributeDefinitions: []*dynamodb.AttributeDefinition{
				{
					AttributeName: aws.String(tableDefinition.PartitionKey),
					AttributeType: aws.String("S"),
				},
			},
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String(tableDefinition.PartitionKey),
					KeyType:       aws.String("HASH"),
				},
			},
			TableName: aws.String(tableDefinition.TableName),
		}

		if tableDefinition.SortKey != "" {
			createTableInput.AttributeDefinitions = append(createTableInput.AttributeDefinitions,
				&dynamodb.AttributeDefinition{
					AttributeName: aws.String(tableDefinition.SortKey),
					AttributeType: aws.String("S"),
				},
			)
			createTableInput.KeySchema = append(createTableInput.KeySchema,
				&dynamodb.KeySchemaElement{
					AttributeName: aws.String(tableDefinition.SortKey),
					KeyType:       aws.String("RANGE"),
				},
			)
		}

		if len(tableDefinition.Indexes) != 0 {
			for _, indexDefinition := range tableDefinition.Indexes {
				createTableInput.AttributeDefinitions = append(createTableInput.AttributeDefinitions,
					&dynamodb.AttributeDefinition{
						AttributeName: aws.String(indexDefinition.PartitionKeyColumn),
						AttributeType: aws.String("S"),
					},
				)
				createTableInput.GlobalSecondaryIndexes = append(createTableInput.GlobalSecondaryIndexes,
					&dynamodb.GlobalSecondaryIndex{
						IndexName: aws.String(indexDefinition.IndexName),
						KeySchema: []*dynamodb.KeySchemaElement{
							{
								AttributeName: aws.String(indexDefinition.PartitionKeyColumn),
								KeyType:       aws.String("HASH"),
							},
						},
						Projection: &dynamodb.Projection{
							ProjectionType: aws.String("KEYS_ONLY"),
						},
					},
				)
			}
		}
		_, err := db.DbClient.CreateTable(createTableInput)

		if err != nil {
			if _, ok := err.(*dynamodb.ResourceInUseException); ok {
				fmt.Println("Table ", tableDefinition, " already created")
				continue
			}
			fmt.Println(err)
		}
	}
}

func main() {
	InitializeTables()
}

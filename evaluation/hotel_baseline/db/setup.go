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

type TableDefinition struct {
	TableName    string
	PartitionKey string
	SortKey      string
	Indexes      []IndexDefinition
}

func InitializeTables() {

	tableDefinitions := []TableDefinition{
		{
			TableName:    UserTable,
			PartitionKey: "Email",
		},
		{
			TableName:    HotelTable, // shared with CityTable
			PartitionKey: "CityName",
			SortKey:      "HotelName",
		},
		{
			TableName:    RoomTable,
			PartitionKey: "HotelName",
			SortKey:      "RoomId",
		},
		{
			TableName:    ReservationTable,
			PartitionKey: "RoomId",
			SortKey:      "DateIn",
			Indexes: []IndexDefinition{
				{
					IndexName: UsersReservationsIndex,
					Column:    "UserId",
				},
			},
		},
	}

	for _, tableDefinition := range tableDefinitions {
		createTableInput := &dynamodb.CreateTableInput{
			BillingMode: aws.String("PROVISIONED"),
			ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(1),
				WriteCapacityUnits: aws.Int64(1),
			},
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
					AttributeName: aws.String(tableDefinition.PartitionKey),
					KeyType:       aws.String("RANGE"),
				},
			)
		}

		if len(tableDefinition.Indexes) != 0 {
			for _, indexDefinition := range tableDefinition.Indexes {
				createTableInput.GlobalSecondaryIndexes = append(createTableInput.GlobalSecondaryIndexes,
					&dynamodb.GlobalSecondaryIndex{
						IndexName: aws.String(indexDefinition.IndexName),
						KeySchema: []*dynamodb.KeySchemaElement{
							{
								AttributeName: aws.String(indexDefinition.Column),
								KeyType:       aws.String("HASH"),
							},
						},
						Projection: &dynamodb.Projection{
							ProjectionType: aws.String("ALL"),
						},
						ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
							ReadCapacityUnits:  aws.Int64(1),
							WriteCapacityUnits: aws.Int64(1),
						},
					},
				)
			}
		}
		_, err := dbClient.CreateTable(createTableInput)

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

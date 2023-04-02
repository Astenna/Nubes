package lib

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type ReferenceNavigationList[T Nobject] struct {
	setup ReferenceNavigationListSetup[T]
}

func NewReferenceNavigationList[T Nobject](ownerId, ownerTypeName, referringFieldName string, isManyToMany bool) *ReferenceNavigationList[T] {
	r := new(ReferenceNavigationList[T])
	r.setup = NewReferenceNavigationListSetup[T](ownerId, ownerTypeName, referringFieldName, isManyToMany)
	r.setup.build()
	return r
}

func (r ReferenceNavigationList[T]) GetIds() ([]string, error) {

	if r.setup.UsesIndex {
		out, err := GetByIndex(r.setup.GetQueryByIndexParam())
		return out, err
	}

	if r.setup.IsManyToMany && !r.setup.UsesIndex {
		input, err := r.setup.GetQueryByPartitionKeyParam()
		if err != nil {
			return nil, err
		}
		out, err := GetSortKeysByPartitionKey(input)
		return out, err
	}

	return nil, fmt.Errorf("invalid initialization of ReferenceNavigationList")
}

func (r ReferenceNavigationList[T]) Get() ([]*T, error) {
	ids, err := r.GetIds()
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return []*T{}, nil
	}

	res, err := LoadBatch[T](ids)
	return res, err
}

func (r ReferenceNavigationList[T]) GetStubs() ([]T, error) {
	var ids []string
	var err error
	if r.setup.UsesIndex {
		ids, err = GetByIndex(r.setup.GetQueryByIndexParam())
		if err != nil {
			return nil, err
		}
	} else if r.setup.IsManyToMany && !r.setup.UsesIndex {
		input, err := r.setup.GetQueryByPartitionKeyParam()
		if err != nil {
			return nil, err
		}
		ids, err = GetSortKeysByPartitionKey(input)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("invalid initialization of ReferenceNavigationList")
	}

	batch, err := GetStubsInBatch[T](ids)

	if err != nil {
		return nil, fmt.Errorf("error occurred while retriving the objects from DB: %w", err)
	}
	return *batch, err
}

func (r ReferenceNavigationList[T]) AddToManyToMany(newId string) error {

	if newId == "" {
		return fmt.Errorf("missing id")
	}

	if r.setup.IsManyToMany {

		typeName := (*new(T)).GetTypeName()
		exists, err := IsInstanceAlreadyCreated(IsInstanceAlreadyCreatedParam{Id: newId, TypeName: typeName})
		if err != nil {
			return fmt.Errorf("error occurred while checking if typename %s with id %s exists. Error %w", typeName, newId, err)
		}
		if !exists {
			return fmt.Errorf("only existing instances can be added to many to many relationships. Typename %s with id %s not found", typeName, newId)
		}

		return InsertToManyToManyTable(r.setup.GetInsertToManyToManyTableParam(newId))
	}

	return fmt.Errorf(`can not add elements to ReferenceNavigationList used as OneToMany relationship. 
						To do it, export or just set the Reference field of the instance with the correct Id`)
}

func (r ReferenceNavigationList[T]) DeleteBatchFromManyToMany(ids []string) error {
	if len(ids) == 0 {
		return fmt.Errorf("missing ids of objects to delete")
	}

	param := r.setup.GetDeleteFromManyToManyParam(ids)
	return DeleteFromManyToManyTable(param)
}

type InsertToManyToManyTableParam struct {
	PartitionKeyName  string
	SortKeyName       string
	PartitionKeyValue string
	SortKeyValue      string
}

func (i InsertToManyToManyTableParam) Verify() error {
	if i.PartitionKeyName == "" {
		return fmt.Errorf("missing PartitionKeyName")
	}
	if i.SortKeyName == "" {
		return fmt.Errorf("missing SortKeyName")
	}
	if i.PartitionKeyValue == "" {
		return fmt.Errorf("missing PartitionKeyValue")
	}
	if i.SortKeyValue == "" {
		return fmt.Errorf("missing SortKeyValue")
	}

	return nil
}

func InsertToManyToManyTable(param InsertToManyToManyTableParam) error {
	if err := param.Verify(); err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(param.PartitionKeyName + param.SortKeyName),
		Item: map[string]*dynamodb.AttributeValue{
			param.PartitionKeyName: {
				S: aws.String(param.PartitionKeyValue),
			},
			param.SortKeyName: {
				S: aws.String(param.SortKeyValue),
			},
		},
	}

	_, err := dbClient.PutItem(input)
	return err
}

func DeleteFromManyToManyTable(param DeleteFromManyToManyParam) error {
	if err := param.Verify(); err != nil {
		return err
	}

	input := dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			param.TableName: {},
		},
	}
	if param.AreIdsToDeletePartitionKeys {
		for _, id := range param.IdsToDelete {
			input.RequestItems[param.TableName] = append(input.RequestItems[param.TableName],
				&dynamodb.WriteRequest{
					DeleteRequest: &dynamodb.DeleteRequest{
						Key: map[string]*dynamodb.AttributeValue{
							param.PartitionKeyName: {
								S: aws.String(id),
							},
							param.SortKeyName: {
								S: aws.String(param.SortKeyValue),
							},
						},
					},
				})
		}

	} else {
		for _, id := range param.IdsToDelete {
			input.RequestItems[param.TableName] = append(input.RequestItems[param.TableName],
				&dynamodb.WriteRequest{
					DeleteRequest: &dynamodb.DeleteRequest{
						Key: map[string]*dynamodb.AttributeValue{
							param.PartitionKeyName: {
								S: aws.String(param.PartitionKeyValue),
							},
							param.SortKeyName: {
								S: aws.String(id),
							},
						},
					},
				})
		}
	}

	_, err := dbClient.BatchWriteItem(&input)
	return err
}

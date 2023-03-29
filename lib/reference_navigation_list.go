package lib

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type ReferenceNavigationList[T Nobject] struct {
	ownerId       string
	ownerTypeName string

	isManyToMany             bool
	usesIndex                bool
	queryByIndexParam        QueryByIndexParam
	queryByPartitionKeyParam QueryByPartitionKeyParam
}

func NewReferenceNavigationList[T Nobject](ownerId, ownerTypeName, referringFieldName string, isManyToMany bool) *ReferenceNavigationList[T] {
	r := new(ReferenceNavigationList[T])
	r.ownerId = ownerId
	r.ownerTypeName = ownerTypeName
	r.isManyToMany = isManyToMany

	if isManyToMany {
		r.setupManyToManyRelationship()
	} else {
		r.setupOneToManyRelationship(referringFieldName)
	}

	if r.usesIndex {
		r.queryByIndexParam.KeyAttributeValue = r.ownerId
	}

	return r
}

func (r ReferenceNavigationList[T]) GetIds() ([]string, error) {

	if r.usesIndex {
		out, err := GetByIndex(r.queryByIndexParam)
		return out, err
	}

	if r.isManyToMany && !r.usesIndex {
		out, err := GetSortKeysByPartitionKey(r.queryByPartitionKeyParam)
		return out, err
	}

	return nil, fmt.Errorf("invalid initialization of ReferenceNavigationList")
}

func (r ReferenceNavigationList[T]) GetLoaded() ([]*T, error) {
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
	if r.usesIndex {
		ids, err = GetByIndex(r.queryByIndexParam)
		if err != nil {
			return nil, err
		}
	} else if r.isManyToMany && !r.usesIndex {
		ids, err = GetSortKeysByPartitionKey(r.queryByPartitionKeyParam)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("invalid initialization of ReferenceNavigationList")
	}

	batch, err := GetBatch[T](ids)

	if err != nil {
		return nil, fmt.Errorf("error occurred while retriving the objects from DB: %w", err)
	}
	return *batch, err
}

func (r ReferenceNavigationList[T]) AddToManyToMany(newId string) error {

	if newId == "" {
		return fmt.Errorf("missing id")
	}

	if r.isManyToMany {

		typeName := (*new(T)).GetTypeName()
		exists, err := IsInstanceAlreadyCreated(IsInstanceAlreadyCreatedParam{Id: newId, TypeName: typeName})
		if err != nil {
			return fmt.Errorf("error occurred while checking if typename %s with id %s exists. Error %w", typeName, newId, err)
		}
		if !exists {
			return fmt.Errorf("only existing instances can be added to many to many relationships. Typename %s with id %s not found", typeName, newId)
		}

		if r.usesIndex {
			return InsertToManyToManyTable(typeName, r.ownerTypeName, newId, r.ownerId)
		}
		return InsertToManyToManyTable(r.ownerTypeName, typeName, r.ownerId, newId)
	}

	return fmt.Errorf("can not add elements to ReferenceNavigationList of OneToMany relationship")
}

func (r *ReferenceNavigationList[T]) setupOneToManyRelationship(referringFieldName string) {
	otherTypeName := (*(new(T))).GetTypeName()
	r.queryByIndexParam.KeyAttributeName = referringFieldName
	r.queryByIndexParam.OutputAttributeName = "Id"
	r.usesIndex = true
	r.queryByIndexParam.TableName = otherTypeName
	r.queryByIndexParam.IndexName = otherTypeName + referringFieldName
}

func (r *ReferenceNavigationList[T]) setupManyToManyRelationship() {
	otherTypeName := (*(new(T))).GetTypeName()

	for index := 0; ; index++ {

		if index >= len(r.ownerTypeName) {
			r.queryByPartitionKeyParam.TableName = r.ownerTypeName + otherTypeName
			r.usesIndex = false
			break
		}
		if index >= len(otherTypeName) {
			r.queryByIndexParam.TableName = otherTypeName + r.ownerTypeName
			r.queryByIndexParam.IndexName = r.queryByIndexParam.TableName + "Reversed"
			r.usesIndex = true
			break
		}

		if r.ownerTypeName[index] < otherTypeName[index] {
			r.queryByPartitionKeyParam.TableName = r.ownerTypeName + otherTypeName
			r.usesIndex = false
			break
		} else if r.ownerTypeName[index] > otherTypeName[index] {
			r.queryByIndexParam.TableName = otherTypeName + r.ownerTypeName
			r.queryByIndexParam.IndexName = r.queryByIndexParam.TableName + "Reversed"
			r.usesIndex = true
			break
		}
	}

	if r.usesIndex {
		r.queryByIndexParam.KeyAttributeName = r.ownerTypeName
		r.queryByIndexParam.OutputAttributeName = otherTypeName
	} else {
		r.queryByPartitionKeyParam.PartitionAttributeName = r.ownerTypeName
		r.queryByPartitionKeyParam.PatritionAttributeValue = r.ownerId
		r.queryByPartitionKeyParam.OutputAttributeName = otherTypeName
	}
}

func InsertToManyToManyTable(partitionKeyName, sortKeyName, partitonKey, sortKey string) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String(partitionKeyName + sortKeyName),
		Item: map[string]*dynamodb.AttributeValue{
			partitionKeyName: {
				S: aws.String(partitonKey),
			},
			sortKeyName: {
				S: aws.String(sortKey),
			},
		},
	}

	_, err := DBClient.PutItem(input)
	return err
}

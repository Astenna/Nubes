package lib

import "fmt"

type ReferenceNavigationListSetup[T Nobject] struct {
	ownerId            string
	ownerTypeName      string
	referringFieldName string

	TableName    string
	IsManyToMany bool
	UsesIndex    bool
	IndexName    string
}

func NewReferenceNavigationListSetup[T Nobject](ownerId, ownerTypeName, referringFieldName string, isManyToMany bool) ReferenceNavigationListSetup[T] {
	setup := ReferenceNavigationListSetup[T]{
		ownerId:            ownerId,
		ownerTypeName:      ownerTypeName,
		referringFieldName: referringFieldName,
		IsManyToMany:       isManyToMany,
	}

	setup.build()
	return setup
}

func (r *ReferenceNavigationListSetup[T]) build() {
	otherTypeName := (*(new(T))).GetTypeName()

	// in order to build the name of the join table,
	// determine the lexicographical order of the two typenames
	// the first in order is the Primary Key, the second is the sort key
	if r.IsManyToMany {
		for index := 0; ; index++ {

			if index >= len(r.ownerTypeName) {
				r.TableName = r.ownerTypeName + otherTypeName
				r.UsesIndex = false
				break
			}
			if index >= len(otherTypeName) {
				r.TableName = otherTypeName + r.ownerTypeName
				r.IndexName = r.TableName + "Reversed"
				r.UsesIndex = true
				break
			}

			if r.ownerTypeName[index] < otherTypeName[index] {
				r.TableName = r.ownerTypeName + otherTypeName
				r.UsesIndex = false
				break
			} else if r.ownerTypeName[index] > otherTypeName[index] {
				r.TableName = otherTypeName + r.ownerTypeName
				r.IndexName = r.TableName + "Reversed"
				r.UsesIndex = true
				break
			}
		}

		return
	}

	// it's a bidirectional one-to-many relationship
	r.UsesIndex = true
	r.TableName = otherTypeName
	r.IndexName = otherTypeName + r.referringFieldName
}

func (r ReferenceNavigationListSetup[T]) GetQueryByIndexParam() QueryByIndexParam {
	result := QueryByIndexParam{}

	result.TableName = r.TableName
	result.IndexName = r.IndexName
	result.KeyAttributeValue = r.ownerId

	if r.IsManyToMany {
		result.KeyAttributeName = r.ownerTypeName
		result.OutputAttributeName = (*(new(T))).GetTypeName()
	} else {
		result.KeyAttributeName = r.referringFieldName
		result.OutputAttributeName = "Id"
	}

	return result
}

func (r ReferenceNavigationListSetup[T]) GetQueryByPartitionKeyParam() (QueryByPartitionKeyParam, error) {
	result := QueryByPartitionKeyParam{}

	if !r.IsManyToMany {
		return result, fmt.Errorf("queries by Partition Key can only be used in ManyToMany relationships")
	}

	result.TableName = r.TableName
	result.PartitionAttributeName = r.ownerTypeName
	result.PatritionAttributeValue = r.ownerId
	result.OutputAttributeName = (*(new(T))).GetTypeName()

	return result, nil
}

func (r ReferenceNavigationListSetup[T]) GetInsertToManyToManyTableParam(newId string) InsertToManyToManyTableParam {
	result := InsertToManyToManyTableParam{}

	if r.UsesIndex {
		result.PartitionKeyName = (*(new(T))).GetTypeName()
		result.SortKeyName = r.ownerTypeName
		result.PartitionKeyValue = newId
		result.SortKeyValue = r.ownerId
		return result
	}

	result.PartitionKeyName = r.ownerTypeName
	result.SortKeyName = (*(new(T))).GetTypeName()
	result.PartitionKeyValue = r.ownerId
	result.SortKeyValue = newId
	return result
}

type DeleteFromManyToManyParam struct {
	TableName                   string
	PartitionKeyName            string
	PartitionKeyValue           string
	SortKeyName                 string
	SortKeyValue                string
	IdsToDelete                 []string
	AreIdsToDeletePartitionKeys bool
}

func (d DeleteFromManyToManyParam) Verify() error {
	if d.TableName == "" {
		return fmt.Errorf("missing TableName")
	}
	if d.PartitionKeyName == "" {
		return fmt.Errorf("missing PartitionKeyName")
	}
	if d.SortKeyName == "" {
		return fmt.Errorf("missing SortKeyName")
	}
	if len(d.IdsToDelete) == 0 {
		return fmt.Errorf("missing IdsToDelete")
	}

	return nil
}

func (r ReferenceNavigationListSetup[T]) GetDeleteFromManyToManyParam(ids []string) DeleteFromManyToManyParam {
	result := DeleteFromManyToManyParam{IdsToDelete: ids}
	result.TableName = r.TableName

	if r.UsesIndex {
		result.PartitionKeyName = (*(new(T))).GetTypeName()
		result.SortKeyName = r.ownerTypeName
		result.SortKeyValue = r.ownerId
		result.AreIdsToDeletePartitionKeys = true
		return result
	}

	result.PartitionKeyName = r.ownerTypeName
	result.SortKeyName = (*(new(T))).GetTypeName()
	result.PartitionKeyValue = r.ownerId
	return result
}

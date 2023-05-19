package lib

import "fmt"

type referenceNavigationListSetup struct {
	ownerId            string
	ownerTypeName      string
	otherTypeName      string
	referringFieldName string

	TableName    string
	IsManyToMany bool
	UsesIndex    bool
	IndexName    string
}

func newReferenceNavigationListSetup(param ReferenceNavigationListParam) referenceNavigationListSetup {
	setup := referenceNavigationListSetup{
		ownerId:            param.OwnerId,
		ownerTypeName:      param.OwnerTypeName,
		otherTypeName:      param.OtherTypeName,
		referringFieldName: param.ReferringFieldName,
		IsManyToMany:       param.IsManyToMany,
	}

	setup.build()
	return setup
}

func (r *referenceNavigationListSetup) build() {
	// in order to build the name of the join table,
	// determine the lexicographical order of the two typenames
	// the first in order is the Primary Key, the second is the sort key
	if r.IsManyToMany {
		for index := 0; ; index++ {

			if index >= len(r.ownerTypeName) {
				r.TableName = r.ownerTypeName + r.otherTypeName
				r.UsesIndex = false
				break
			}
			if index >= len(r.otherTypeName) {
				r.TableName = r.otherTypeName + r.ownerTypeName
				r.IndexName = r.TableName + "Reversed"
				r.UsesIndex = true
				break
			}

			if r.ownerTypeName[index] < r.otherTypeName[index] {
				r.TableName = r.ownerTypeName + r.otherTypeName
				r.UsesIndex = false
				break
			} else if r.ownerTypeName[index] > r.otherTypeName[index] {
				r.TableName = r.otherTypeName + r.ownerTypeName
				r.IndexName = r.TableName + "Reversed"
				r.UsesIndex = true
				break
			}
		}

		return
	}

	// it's a bidirectional one-to-many relationship
	r.UsesIndex = true
	r.TableName = r.otherTypeName
	r.IndexName = r.otherTypeName + r.referringFieldName
}

func (r referenceNavigationListSetup) GetQueryByIndexParam() QueryByIndexParam {
	result := QueryByIndexParam{}

	result.TableName = r.TableName
	result.IndexName = r.IndexName
	result.KeyAttributeValue = r.ownerId

	if r.IsManyToMany {
		result.KeyAttributeName = r.ownerTypeName
		result.OutputAttributeName = r.otherTypeName
	} else {
		result.KeyAttributeName = r.referringFieldName
		result.OutputAttributeName = "Id"
	}

	return result
}

func (r referenceNavigationListSetup) GetQueryByPartitionKeyParam() (QueryByPartitionKeyParam, error) {
	result := QueryByPartitionKeyParam{}

	if !r.IsManyToMany {
		return result, fmt.Errorf("queries by Partition Key can only be used in ManyToMany relationships")
	}

	result.TableName = r.TableName
	result.PartitionAttributeName = r.ownerTypeName
	result.PatritionAttributeValue = r.ownerId
	result.OutputAttributeName = r.otherTypeName

	return result, nil
}

func (r referenceNavigationListSetup) GetInsertToManyToManyTableParam(newId string) InsertToManyToManyTableLibParam {
	result := InsertToManyToManyTableLibParam{}

	if r.UsesIndex {
		result.PartitionKeyName = r.otherTypeName
		result.SortKeyName = r.ownerTypeName
		result.PartitionKeyValue = newId
		result.SortKeyValue = r.ownerId
		return result
	}

	result.PartitionKeyName = r.ownerTypeName
	result.SortKeyName = r.otherTypeName
	result.PartitionKeyValue = r.ownerId
	result.SortKeyValue = newId
	return result
}

type DeleteFromManyToManyLibParam struct {
	TableName                   string
	PartitionKeyName            string
	PartitionKeyValue           string
	SortKeyName                 string
	SortKeyValue                string
	IdsToDelete                 []string
	AreIdsToDeletePartitionKeys bool
}

func (r referenceNavigationListSetup) GetDeleteFromManyToManyParam(ids []string) DeleteFromManyToManyLibParam {
	result := DeleteFromManyToManyLibParam{IdsToDelete: ids}
	result.TableName = r.TableName

	if r.UsesIndex {
		result.PartitionKeyName = r.otherTypeName
		result.SortKeyName = r.ownerTypeName
		result.SortKeyValue = r.ownerId
		result.AreIdsToDeletePartitionKeys = true
		return result
	}

	result.PartitionKeyName = r.ownerTypeName
	result.SortKeyName = r.otherTypeName
	result.PartitionKeyValue = r.ownerId
	return result
}

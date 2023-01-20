package lib

import "fmt"

type ReferenceNavigationList[T Nobject] struct {
	ownerId            string
	ownerTypeName      string
	referringFieldName string

	isManyToMany             bool
	usesIndex                bool
	queryByIndexParam        QueryByIndexParam
	QueryByPartitionKeyParam QueryByPartitionKeyParam
}

func NewReferenceNavigationList[T Nobject](ownerId, ownerTypeName, referringFieldName string, isManyToMany bool) *ReferenceNavigationList[T] {
	r := new(ReferenceNavigationList[T])
	r.ownerId = ownerId
	r.ownerTypeName = ownerTypeName
	r.referringFieldName = referringFieldName
	r.isManyToMany = isManyToMany

	if isManyToMany {
		r.setupManyToManyRelationship()
	} else {
		r.setupOneToManyRelationship()
	}

	if r.usesIndex {
		r.queryByIndexParam.KeyAttributeValue = r.ownerId
	}

	return r
}

func (r ReferenceNavigationList[T]) GetIds() ([]string, error) {

	if r.usesIndex {
		out, err := GetByIndex[T](r.queryByIndexParam)
		return out, err
	}

	if r.isManyToMany && !r.usesIndex {
		out, err := GetSortKeysByPartitionKey[T](r.QueryByPartitionKeyParam)
		return out, err
	}

	return nil, fmt.Errorf("invalid initialization of ReferenceNavigationList")
}

func (r *ReferenceNavigationList[T]) setupOneToManyRelationship() {
	otherTypeName := (*(new(T))).GetTypeName()
	r.queryByIndexParam.KeyAttributeName = r.referringFieldName
	r.queryByIndexParam.OutputAttributeName = "Id"
	r.usesIndex = true
	r.queryByIndexParam.TableName = otherTypeName
	r.queryByIndexParam.IndexName = otherTypeName + r.referringFieldName
}

func (r *ReferenceNavigationList[T]) setupManyToManyRelationship() {
	otherTypeName := (*(new(T))).GetTypeName()

	for index := 0; ; index++ {

		if index >= len(r.ownerTypeName) {
			r.QueryByPartitionKeyParam.TableName = r.ownerTypeName + otherTypeName
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
			r.QueryByPartitionKeyParam.TableName = r.ownerTypeName + otherTypeName
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
		r.QueryByPartitionKeyParam.PartitionAttributeName = r.ownerTypeName
		r.QueryByPartitionKeyParam.PatritionAttributeValue = r.ownerId
		r.QueryByPartitionKeyParam.OutputAttributeName = otherTypeName
	}
}

package lib

import (
	"fmt"
)

type ReferenceNavigationListHandlers struct {
	setup referenceNavigationListSetup
}

func NewReferenceNavigationListHandlers(param ReferenceNavigationListParam) *ReferenceNavigationListHandlers {
	r := new(ReferenceNavigationListHandlers)
	r.setup = newReferenceNavigationListSetup(param)
	return r
}

func (r ReferenceNavigationListHandlers) GetIds() ([]string, error) {

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

	return nil, fmt.Errorf("invalid initialization of ReferenceNavigationListHandlers")
}

func (r ReferenceNavigationListHandlers) Get() ([]string, error) {
	ids, err := r.GetIds()
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, nil
	}
	err = AreInstancesAlreadyCreated(LoadBatchParam{
		TypeName: r.setup.otherTypeName,
		Ids:      ids,
	})
	return ids, err
}

func (r ReferenceNavigationListHandlers) GetStubs() ([]interface{}, error) {
	ids, err := r.GetIds()
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, nil
	}

	return GetStubsInBatchWithTypeNameAsArg(GetBatchParam{
		Ids:      ids,
		TypeName: r.setup.otherTypeName,
	})
}

func (r ReferenceNavigationListHandlers) AddToManyToMany(newId string) error {

	if newId == "" {
		return fmt.Errorf("missing id")
	}

	if r.setup.IsManyToMany {

		typeName := r.setup.ownerTypeName
		exists, err := IsInstanceAlreadyCreated(IsInstanceAlreadyCreatedParam{Id: newId, TypeName: r.setup.otherTypeName})
		if err != nil {
			return fmt.Errorf("error occurred while checking if typename %s with id %s exists. Error %w", typeName, newId, err)
		}
		if !exists {
			return fmt.Errorf("only existing instances can be added to many to many relationships. Typename %s with id %s not found", typeName, newId)
		}

		return InsertToManyToManyTable(r.setup.GetInsertToManyToManyTableParam(newId))
	}

	return fmt.Errorf(`can not add elements to ReferenceNavigationListHandlers used as OneToMany relationship. 
						To do it, export or just set the Reference field of the instance with the correct Id`)
}

func (r ReferenceNavigationListHandlers) DeleteBatchFromManyToMany(ids []string) error {
	if len(ids) == 0 {
		return fmt.Errorf("missing ids of objects to delete")
	}

	param := r.setup.GetDeleteFromManyToManyParam(ids)
	return DeleteFromManyToManyTable(param)
}

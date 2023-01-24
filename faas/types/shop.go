package types

import (
	"github.com/Astenna/Nubes/lib"
)

type Shop struct {
	Id		string
	Name		string
	Owners		lib.ReferenceNavigationList[User]	`nubes:"hasMany-Shops" dynamodbav:"-"`
	Products	lib.ReferenceNavigationList[Product]	`nubes:"hasOne-SoldBy,readonly" dynamodbav:"-"`
	isInitialized	bool
}

func (Shop) GetTypeName() string {
	return "Shop"
}

func (s Shop) GetOwners() ([]string, error) {
	if s.isInitialized {
		fieldValue, _libError := lib.GetField(lib.GetFieldParam{Id: s.Id, TypeName: "Shop", FieldName: "Owners"})
		if _libError != nil {
			return *new([]string), _libError
		}
		s.Owners = fieldValue.(lib.ReferenceNavigationList[User])
	}
	return s.Owners.GetIds()
}

func (receiver *Shop) Init() {
	receiver.isInitialized = true
	receiver.Products = *lib.NewReferenceNavigationList[Product](receiver.Id, receiver.GetTypeName(), "SoldBy", false)
	receiver.Owners = *lib.NewReferenceNavigationList[User](receiver.Id, receiver.GetTypeName(), "", true)
}

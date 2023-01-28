package types

import (
	"fmt"

	"github.com/Astenna/Nubes/lib"
)

type Shop struct {
	Id            string
	Name          string
	Owners        lib.ReferenceNavigationList[User]    `nubes:"hasMany-Shops" dynamodbav:"-"`
	Products      lib.ReferenceNavigationList[Product] `nubes:"hasOne-SoldBy,readonly" dynamodbav:"-"`
	isInitialized bool
}

func (Shop) GetTypeName() string {
	return "Shop"
}

func (s Shop) GetOwners() ([]string, error) {
	if !s.isInitialized {
		return nil, fmt.Errorf(`fields of type ReferenceNavigationList can be used only after instance initialization. 
			Use lib.Load or lib.Export from the Nubes library to create initialized instances`)
	}
	return s.Owners.GetIds()
}

func (receiver *Shop) Init() {
	receiver.isInitialized = true
	receiver.Products = *lib.NewReferenceNavigationList[Product](receiver.Id, receiver.GetTypeName(), "SoldBy", false)
	receiver.Owners = *lib.NewReferenceNavigationList[User](receiver.Id, receiver.GetTypeName(), "", true)
}

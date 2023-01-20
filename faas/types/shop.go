package types

import (
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

func (receiver *Shop) Init() {
	receiver.isInitialized = true
	receiver.Products = *lib.NewReferenceNavigationList[Product](receiver.Id, receiver.GetTypeName(), "SoldBy", false)
	receiver.Owners = *lib.NewReferenceNavigationList[User](receiver.Id, receiver.GetTypeName(), "", true)
}

package types

import (
	"github.com/Astenna/Nubes/lib"
)

type Shop struct {
	Id       string
	Name     string
	Owners   lib.ReferenceNavigationList[User]    `nubes:"hasMany-Shops" dynamodbav:"-"`
	Products lib.ReferenceNavigationList[Product] `nubes:"hasOne-SoldBy,readonly" dynamodbav:"-"`
}

func (Shop) GetTypeName() string {
	return "Shop"
}

func (s *Shop) Init() error {
	s.Products = *lib.NewReferenceNavigationList[Product](s.Id, s.GetTypeName(), "SoldBy", false)
	s.Owners = *lib.NewReferenceNavigationList[User](s.Id, s.GetTypeName(), "", true)
	return nil
}

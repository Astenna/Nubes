package types

import (
	"github.com/Astenna/Nubes/lib"
)

type Shop struct {
	Id       string
	Name     string
	Owners   lib.ReferenceList[User]
	Products lib.ReferenceNavigationList[Product] `nubes:"hasOne-SoldBy,readonly"`
}

func (Shop) GetTypeName() string {
	return "Shop"
}

func (s *Shop) Init() error {
	s.Products = *lib.NewReferenceNavigationList[Product](s.Id, s.GetTypeName(), "SoldBy")
	return nil
}

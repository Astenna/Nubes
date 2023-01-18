package types

import (
	"github.com/Astenna/Nubes/lib"
)

type Shop struct {
	Id       string
	Name     string
	Owners   lib.ReferenceList[User]                    `nubes:"has-many-Shops"`
	Products lib.ReferenceNavigationList[Shop, Product] `nubes:"has-one-SoldBy"`
}

func (Shop) GetTypeName() string {
	return "Shop"
}

package types

import (
	"github.com/Astenna/Nubes/lib"
)

type Shop struct {
	Id	string
	Name	string
	Owners	*lib.ReferenceList[User]	`nubes:"index"`
}

func NewShop(shop Shop) (Shop, error) {
	out, _libError := lib.Insert(shop)
	if _libError != nil {
		return *new(Shop), _libError
	}
	shop.Id = out
	return shop, nil
}

func ReNewShop(id string) Shop {
	shop := new(Shop)
	shop.Id = id
	return *shop
}

func (Shop) GetTypeName() string {
	return "Shop"
}

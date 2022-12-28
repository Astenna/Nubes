package repositories

import (
	"errors"

	"github.com/Astenna/Nubes/faas/types"
	lib "github.com/Astenna/Nubes/lib"
)

func CreateShop(shop types.Shop) (string, error) {
	if shop.Name == "" {
		return "", errors.New("shop name can not be empty")
	}

	return lib.Insert(shop)
}

func GetShop(id string) (*types.Shop, error) {
	return lib.Get[types.Shop](id)
}

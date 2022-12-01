package faas

import (
	"errors"
	"github.com/Astenna/Thesis_PoC/faas/types"
	lib "github.com/Astenna/Thesis_PoC/faas_lib"
)

func CreateShop(shop types.Shop) error {
	if shop.Name == "" {
		return errors.New("shop name can not be empty")
	}

	return lib.Create(shop)
}

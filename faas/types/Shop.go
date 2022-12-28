package types

import (
	"github.com/Astenna/Nubes/lib"
)

type Shop struct {
	Id    string
	Name  string
	Owner *lib.Reference[User]
}

func (Shop) GetTypeName() string {
	return "Shop"
}

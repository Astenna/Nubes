package types

import "github.com/Astenna/Thesis_PoC/faas_lib"

type Shop struct {
	Id    int
	Name  string
	Owner *faas_lib.Reference[User]
}

func (Shop) GetTypeName() string {
	return "Shop"
}

func NewShop(ownerId string) *Shop {
	return &Shop{Owner: faas_lib.NewReference[User](ownerId)}
}

func (s *Shop) ChangeName(name string) {
	s.Name = name
}

func (s *Shop) ChangeOwner(ownerId string) {
	s.Owner = faas_lib.NewReference[User](ownerId)
}

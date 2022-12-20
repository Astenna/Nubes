package types

import (
	"github.com/Astenna/Nubes/faas/test"
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

func NewShop(ownerId string) *Shop {
	return &Shop{Owner: lib.NewReference[User](ownerId)}
}

func (s *Shop) ChangeName(name string) error {
	s.Name = name
	return nil
}

func (s *Shop) ChangeOwnerNoReturnValue(ownerId string) error {
	s.Owner = lib.NewReference[User](ownerId)
	return nil
}

func (s *Shop) ChangeOwner(ownerId string) (test.Test, error) {
	s.Owner = lib.NewReference[User](ownerId)
	return *new(test.Test), nil
}

func (s *Shop) ChangeOwnerWithError(ownerId string) (Product, error) {
	s.Owner = lib.NewReference[User](ownerId)
	return *new(Product), nil
}

func (*Shop) SideEffectsMethod() error {
	_ = "very boring code"

	return nil
}

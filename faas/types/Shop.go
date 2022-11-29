package types

import faas_lib "github.com/Astenna/Thesis_PoC/faas_lib"

type Shop struct {
	Id    int
	Name  string
	Owner faas_lib.Reference[User]
}

func NewShop(owner faas_lib.Reference[User]) *Shop {
	return &Shop{Owner: owner}
}

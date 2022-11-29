package types

type Shop struct {
	Id    int
	Name  string
	Owner FaaSLib.Reference[User]
}

func NewShop(owner FaaSLib.Reference[User]) *Shop {
	return &Shop{Owner: owner}
}

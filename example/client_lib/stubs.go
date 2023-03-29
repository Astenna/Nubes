package client_lib

type ProductStub struct {
	Id string

	Name string

	QuantityAvailable int

	SoldBy Reference[shop] `dynamodbav:",omitempty"`

	Discount ReferenceList[discount]

	Price float64
}

func (ProductStub) GetTypeName() string {
	return "Product"
}

type ShippingStub struct {
	Id string

	Address string

	State ShippingState
}

func (ShippingStub) GetTypeName() string {
	return "Shipping"
}

type ShopStub struct {
	Id string

	Name string
}

func (ShopStub) GetTypeName() string {
	return "Shop"
}

type UserStub struct {
	FirstName string

	LastName string

	Email string `nubes:"id,readonly" dynamodbav:"Id"`

	Password string `nubes:"readonly"`

	Address string

	Orders ReferenceList[order]
}

func (UserStub) GetTypeName() string {
	return "User"
}

type DiscountStub struct {
	Id string

	Percentage string
}

func (DiscountStub) GetTypeName() string {
	return "Discount"
}

type OrderStub struct {
	Id string

	Products []OrderedProduct

	Buyer Reference[user]

	Shipping Reference[shipping]
}

func (OrderStub) GetTypeName() string {
	return "Order"
}

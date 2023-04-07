package client_lib


 

 



func DiscountReferenceList(capacity ...int) ReferenceList[discount] {
	if capacity != nil {
		return make(ReferenceList[discount], 0, capacity[0])
	}
	return *new(ReferenceList[discount])
}
 



func OrderReferenceList(capacity ...int) ReferenceList[order] {
	if capacity != nil {
		return make(ReferenceList[order], 0, capacity[0])
	}
	return *new(ReferenceList[order])
}
 

 



func ShopReferenceList(capacity ...int) ReferenceList[shop] {
	if capacity != nil {
		return make(ReferenceList[shop], 0, capacity[0])
	}
	return *new(ReferenceList[shop])
}
 



func ProductReferenceList(capacity ...int) ReferenceList[product] {
	if capacity != nil {
		return make(ReferenceList[product], 0, capacity[0])
	}
	return *new(ReferenceList[product])
}
 



func ShippingReferenceList(capacity ...int) ReferenceList[shipping] {
	if capacity != nil {
		return make(ReferenceList[shipping], 0, capacity[0])
	}
	return *new(ReferenceList[shipping])
}
 



func UserReferenceList(capacity ...int) ReferenceList[user] {
	if capacity != nil {
		return make(ReferenceList[user], 0, capacity[0])
	}
	return *new(ReferenceList[user])
}
 

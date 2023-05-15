package client_lib

const (
	InPreparation   ShippingState = "InPreparation"
	InTransit       ShippingState = "InTransit"
	PickupAvailavle ShippingState = "PickupAvailable"
	Delivered       ShippingState = "Delivered"
)

type ShippingState string

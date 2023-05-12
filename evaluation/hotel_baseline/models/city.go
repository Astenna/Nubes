package models

type City struct {
	CityName string
	// only for dynamoDB queries (= sort key)
	HotelName   string
	Region      string
	Description string
}

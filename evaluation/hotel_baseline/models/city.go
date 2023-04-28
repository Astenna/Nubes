package models

type City struct {
	CityName string
	// only for dynamoDB queries (= sort key)
	HotelName   string `json:"-"`
	Region      string
	Description string
}

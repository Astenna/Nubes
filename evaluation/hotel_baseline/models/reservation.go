package models

import (
	"time"
)

type Reservation struct {
	CityHotelRoomId string // CityName_HotelName_RoomId
	DateIn          time.Time
	DateOut         time.Time
}

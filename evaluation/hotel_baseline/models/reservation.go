package models

import (
	"time"
)

type Reservation struct {
	CityHotelNameRoomId string // CityName_HotelName_RoomId
	DateIn              time.Time
	DateOut             time.Time
	UserEmail           string
}

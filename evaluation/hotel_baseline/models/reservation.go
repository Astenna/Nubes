package models

import (
	"time"
)

type Reservation struct {
	RoomId    string
	DateIn    time.Time
	DateOut   time.Time
	UserId    string
	HotelName string // to retrieve related room
}

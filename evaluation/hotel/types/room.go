package types

import (
	"time"

	"github.com/Astenna/Nubes/lib"
)

type Room struct {
	Id           string
	Name         string
	Description  string
	Hotel        lib.Reference[Hotel]
	Reservations []ReservationInOut
	Price        float32
}

type ReservationInOut struct {
	In  time.Time
	Out time.Time
}

func (o Room) GetTypeName() string {
	return "Room"
}

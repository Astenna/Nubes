package types

import (
	"time"

	"github.com/Astenna/Nubes/lib"
)

type Reservation struct {
	Id      string `nubes:"id,readonly" dynamodbav:"Id"` // Id: HotelId:dateIn
	Room    lib.Reference[Room]
	User    lib.Reference[User]
	DateOut time.Time
}

func (o Reservation) GetTypeName() string {
	return "Reservation"
}

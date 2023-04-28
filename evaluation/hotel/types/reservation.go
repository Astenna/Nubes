package types

import (
	"time"

	"github.com/Astenna/Nubes/lib"
)

type Reservation struct {
	Id              string `nubes:"id,readonly" dynamodbav:"Id"` // Id: HotelId_dateIn
	Room            lib.Reference[Room]
	User            lib.Reference[User] `dynamodbav:",omitempty"`
	DateOut         time.Time
	isInitialized   bool
	invocationDepth int
}

func (o Reservation) GetTypeName() string {
	return "Reservation"
}
func (receiver Reservation) GetId() string {
	return receiver.Id
}
func (receiver *Reservation) Init() {
	receiver.isInitialized = true
}
func (receiver *Reservation) saveChangesIfInitialized() error {
	if receiver.isInitialized && receiver.invocationDepth == 1 {
		_libError := lib.Upsert(receiver, receiver.Id)
		if _libError != nil {
			return _libError
		}
	}
	return nil
}

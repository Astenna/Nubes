package types

import (
	"time"

	"github.com/Astenna/Nubes/lib"
)

type Room struct {
	Id              string
	Name            string
	Description     string
	Hotel           lib.Reference[Hotel] `dynamodbav:",omitempty"`
	Reservations    []ReservationInOut
	Price           float32
	isInitialized   bool
	invocationDepth int
}

type ReservationInOut struct {
	In  time.Time
	Out time.Time
}

func (o Room) GetTypeName() string {
	return "Room"
}
func (receiver *Room) Init() {
	receiver.isInitialized = true
}
func (receiver *Room) saveChangesIfInitialized() error {
	if receiver.isInitialized && receiver.invocationDepth == 1 {
		_libError := lib.Upsert(receiver, receiver.Id)
		if _libError != nil {
			return _libError
		}
	}
	return nil
}

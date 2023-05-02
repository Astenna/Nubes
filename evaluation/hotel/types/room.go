package types

import (
	"github.com/Astenna/Nubes/lib"
)

type Room struct {
	Id              string
	Name            string
	Description     string
	Hotel           lib.Reference[Hotel] `dynamodbav:",omitempty"`
	Reservations    map[string][]ReservationInOut
	Price           float32
	isInitialized   bool
	invocationDepth int
}

type ReservationInOut struct {
	In  int
	Out int
}

func (o Room) GetTypeName() string {
	return "Room"
}

func (o Room) GetReservations() (map[string][]ReservationInOut, error) {
	if o.isInitialized {
		fieldValue := make(map[string][]ReservationInOut)
		_libError := lib.GetFieldOfType(lib.GetStateParam{Id: o.Id, TypeName: "Room", FieldName: "Reservations"}, &fieldValue)
		if _libError != nil {
			return *new(map[string][]ReservationInOut), _libError
		}
		o.Reservations = fieldValue
	}
	return o.Reservations, nil
}

func (o *Room) SetReservations(in map[string][]ReservationInOut) error {
	o.Reservations = in
	if o.isInitialized {
		_libError := lib.SetField(lib.SetFieldParam{Id: o.Id, TypeName: "Room", FieldName: "Reservations", Value: o.Reservations})
		if _libError != nil {
			return _libError
		}
	}
	return nil
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

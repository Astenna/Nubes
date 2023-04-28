package types

import (
	"github.com/Astenna/Nubes/lib"
	"github.com/jftuga/geodist"
)

type Hotel struct {
	HName           string `nubes:"Id" dynamodbav:"Id"`
	Street          string
	PostalCode      string
	Coordinates     geodist.Coord `nubes:"readonly"`
	Rate            float32
	Rooms           lib.ReferenceNavigationList[Room] `nubes:"hasOne-Hotel" dynamodbav:"-"`
	City            lib.Reference[City]               `dynamodbav:",omitempty"`
	isInitialized   bool
	invocationDepth int
}

func (o Hotel) GetTypeName() string {
	return "Hotel"
}
func (receiver Hotel) GetId() string {
	return receiver.HName
}
func (receiver *Hotel) Init() {
	receiver.isInitialized = true
	receiver.Rooms = *lib.NewReferenceNavigationList[Room](receiver.HName, receiver.GetTypeName(), "Hotel", false)
}
func (receiver *Hotel) saveChangesIfInitialized() error {
	if receiver.isInitialized && receiver.invocationDepth == 1 {
		_libError := lib.Upsert(receiver, receiver.HName)
		if _libError != nil {
			return _libError
		}
	}
	return nil
}

package types

import (
	"github.com/Astenna/Nubes/lib"
	"github.com/jftuga/geodist"
)

type Hotel struct {
	Id          string
	HName       string
	Street      string
	PostalCode  string
	Country     string
	Coordinates geodist.Coord `nubes:"readonly"`
	Rate        float32
	Rooms       lib.ReferenceNavigationList[Room] `nubes:"hasOne-Hotel"`
	City        lib.Reference[City]
}

func (o Hotel) GetTypeName() string {
	return "Hotel"
}

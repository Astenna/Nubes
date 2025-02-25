package types

import (
	"time"

	"github.com/Astenna/Nubes/lib"
)

type Discount struct {
	Id              string
	Percentage      string
	isInitialized   bool
	ValidFrom       time.Time
	ValidUntil      time.Time
	invocationDepth int
}

// NewDiscount is a very simple example of custom constructor definition
func NewDiscount() (Discount, error) {
	return Discount{
		ValidFrom: time.Now(),
	}, nil
}

func (d *Discount) SetValidFrom(date time.Time) error {
	d.ValidFrom = date
	if d.isInitialized {
		_libError := lib.SetField(lib.SetFieldParam{Id: d.Id, TypeName: "Discount", FieldName: "ValidFrom", Value: d.ValidFrom})
		if _libError != nil {
			return _libError
		}
	}
	return nil
}

func (d Discount) GetValidFrom() (time.Time, error) {
	if d.isInitialized {
		fieldValue := *new(time.Time)
		_libError := lib.GetFieldOfType(lib.GetStateParam{Id: d.Id, TypeName: "Discount", FieldName: "ValidFrom"}, &fieldValue)
		if _libError != nil {
			return *new(time.Time), _libError
		}
		d.ValidFrom = fieldValue
	}
	return d.ValidFrom, nil
}

func (Discount) GetTypeName() string {
	return "Discount"
}
func (receiver *Discount) Init() {
	receiver.isInitialized = true
}
func (receiver *Discount) saveChangesIfInitialized() error {
	if receiver.isInitialized && receiver.invocationDepth == 1 {
		_libError := lib.Upsert(receiver, receiver.Id)
		if _libError != nil {
			return _libError
		}
	}
	return nil
}

package types

import (
	"errors"

	"github.com/Astenna/Nubes/lib"
)

type Product struct {
	Id                string
	Name              string
	QuantityAvailable float64
	SoldBy            lib.Reference[Shop]
	Price             float64
	isInitialized     bool
}

func (Product) GetTypeName() string {
	return "Product"
}

func (p *Product) DecreaseAvailabilityBy(decreaseNum float64) error {
	if p.isInitialized {
		tempReceiverName, _libError := lib.GetObjectState[Product](p.Id)
		if _libError != nil {
			return _libError
		}
		p = tempReceiverName
		p.Init()
	}

	if p.QuantityAvailable-decreaseNum < 0 {
		return errors.New("not enough quantity available")
	}
	p.QuantityAvailable = p.QuantityAvailable - decreaseNum
	if p.isInitialized {
		_libError := lib.Upsert(p, p.Id)
		if _libError != nil {
			return _libError
		}
	}
	return nil
}

func (p Product) GetQuantityAvailable() (float64, error) {
	if p.isInitialized {
		fieldValue, _libError := lib.GetField(lib.GetFieldParam{Id: p.Id, TypeName: "Product", FieldName: "QuantityAvailable"})
		if _libError != nil {
			return *new(float64), _libError
		}
		p.QuantityAvailable = fieldValue.(float64)
	}
	return p.QuantityAvailable, nil
}

func (p Product) GetSoldBy() (lib.Reference[Shop], error) {
	if p.isInitialized {
		fieldValue, _libError := lib.GetField(lib.GetFieldParam{Id: p.Id, TypeName: "Product", FieldName: "SoldBy"})
		if _libError != nil {
			return *new(lib.Reference[Shop]), _libError
		}
		p.SoldBy = fieldValue.(lib.Reference[Shop])
	}
	return p.SoldBy, nil
}

func (p *Product) SetSoldBy(id string) error {
	p.SoldBy = lib.Reference[Shop](id)
	if p.isInitialized {
		_libError := lib.SetField(lib.SetFieldParam{Id: p.Id, TypeName: "Product", FieldName: "SoldBy", Value: p.SoldBy})
		if _libError != nil {
			return _libError
		}
	}
	return nil
}

func (receiver *Product) Init() {
	receiver.isInitialized = true
}

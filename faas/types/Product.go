package types

import (
	"errors"

	"github.com/Astenna/Nubes/lib"
)

type Product struct {
	Id                string
	Name              string
	QuantityAvailable int
	SoldBy            lib.Reference[Shop]
	Price             float32
}

func (Product) GetTypeName() string {
	return "Product"
}

func (p *Product) DecreaseAvailabilityBy(decreaseNum int) error {
	p, _libError := lib.GetObjectState[Product](p.Id)
	if _libError != nil {
		return _libError
	}
	if p.QuantityAvailable-decreaseNum < 0 {
		return errors.New("not enough quantity available")
	}

	p.QuantityAvailable = p.QuantityAvailable - decreaseNum
	_libError = lib.Upsert(p, p.Id)
	if _libError != nil {
		return _libError
	}
	return nil
}

func (p Product) GetSoldBy() (lib.Reference[Shop], error) {
	fieldValue, _libError := lib.GetField(p.Id, lib.GetFieldParam{TypeName: "Product", FieldName: "SoldBy"})
	if _libError != nil {
		return *new(lib.Reference[Shop]), _libError
	}

	fieldMap, _ := fieldValue.(map[string]string)
	p.SoldBy = *lib.NewReference[Shop](fieldMap["Id"])
	return p.SoldBy, nil
}

func (p *Product) SetSoldBy(id string) error {
	p.SoldBy = *lib.NewReference[Shop](id)
	_libError := lib.SetField(p.Id, lib.SetFieldParam{TypeName: "Product", FieldName: "SoldBy", Value: p.SoldBy})
	if _libError != nil {
		return _libError
	}
	return nil
}

package types

import (
	"errors"

	"github.com/Astenna/Nubes/lib"
)

type Product struct {
	Id			string
	Name			string
	QuantityAvailable	int
	SoldBy			lib.Reference[Shop]	`dynamodbav:",omitempty"`
	Discount		lib.ReferenceList[Discount]
	Price			float64
	isInitialized		bool
	invocationDepth		int
}

func (Product) GetTypeName() string {
	return "Product"
}

func (p *Product) DecreaseAvailabilityBy(decreaseNum int) error {
	p.invocationDepth++
	if p.isInitialized && p.invocationDepth == 1 {
		tempReceiverName, _libError := lib.GetObjectState[Product](p.Id)
		if _libError != nil {
			p.invocationDepth--
			return _libError
		}
		p = tempReceiverName
		p.Init()
	}
	for index, discount := range p.Discount {
		_, _ = index, discount
	}

	if p.QuantityAvailable-decreaseNum < 0 {
		p.invocationDepth--
		return errors.New("not enough quantity available")
	}
	p.QuantityAvailable = p.QuantityAvailable - decreaseNum
	_libUpsertError := p.saveChangesIfInitialized()
	p.invocationDepth--

	return _libUpsertError
}

func (p *Product) privateDecreaseAvailabilityBy(decreaseNum int) error {
	for index, discount := range p.Discount {
		_, _ = index, discount
	}

	if p.QuantityAvailable-decreaseNum < 0 {
		p.invocationDepth--
		return errors.New("not enough quantity available")
	}
	p.QuantityAvailable = p.QuantityAvailable - decreaseNum
	return nil
}

func (p Product) GetQuantityAvailable() (int, error) {
	if p.isInitialized {
		fieldValue, _libError := lib.GetFieldOfType[int](lib.GetStateParam{Id: p.Id, TypeName: "Product", FieldName: "QuantityAvailable"})
		if _libError != nil {
			return *new(int), _libError
		}
		p.QuantityAvailable = fieldValue
	}
	return p.QuantityAvailable, nil
}

func (p Product) GetSoldBy() (lib.Reference[Shop], error) {
	if p.isInitialized {
		fieldValue, _libError := lib.GetFieldOfType[lib.Reference[Shop]](lib.GetStateParam{Id: p.Id, TypeName: "Product", FieldName: "SoldBy"})
		if _libError != nil {
			return *new(lib.Reference[Shop]), _libError
		}
		p.SoldBy = fieldValue
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
func (receiver *Product) saveChangesIfInitialized() error {
	if receiver.isInitialized && receiver.invocationDepth == 1 {
		_libError := lib.Upsert(receiver, receiver.Id)
		if _libError != nil {
			return _libError
		}
	}
	return nil
}

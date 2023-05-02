package types

import (
	"errors"
	"time"

	"github.com/Astenna/Nubes/lib"
)

type Product struct {
	Id                string
	Name              string
	QuantityAvailable int
	SoldBy            lib.Reference[Shop] `dynamodbav:",omitempty"`
	Discount          lib.ReferenceList[Discount]
	Price             float64
	isInitialized     bool
	invocationDepth   int
}

func (Product) GetTypeName() string {
	return "Product"
}

func (p *Product) DecreaseAvailabilityBy(decreaseNum int) error {
	p.invocationDepth++
	if p.isInitialized && p.invocationDepth == 1 {
		_libError := lib.GetStub(p.Id, p)
		if _libError != nil {
			p.invocationDepth--
			return _libError
		}
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

func (p Product) GetQuantityAvailable() (int, error) {
	if p.isInitialized {
		fieldValue := *new(int)
		_libError := lib.GetFieldOfType(lib.GetStateParam{Id: p.Id, TypeName: "Product", FieldName: "QuantityAvailable"}, &fieldValue)
		if _libError != nil {
			return *new(int), _libError
		}
		p.QuantityAvailable = fieldValue
	}
	return p.QuantityAvailable, nil
}

func (p Product) GetName() (string, error) {
	if p.isInitialized {
		fieldValue := *new(string)
		_libError := lib.GetFieldOfType(lib.GetStateParam{Id: p.Id, TypeName: "Product", FieldName: "Name"}, &fieldValue)
		if _libError != nil {
			return *new(string), _libError
		}
		p.Name = fieldValue
	}
	return p.Name, nil
}

func (p Product) GetSoldBy() (lib.Reference[Shop], error) {
	if p.isInitialized {
		fieldValue := *new(lib.Reference[Shop])
		_libError := lib.GetFieldOfType(lib.GetStateParam{Id: p.Id, TypeName: "Product", FieldName: "SoldBy"}, &fieldValue)
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

// Example of a method accepting an Nobject as an input parameter.
// In such case, the invocations from client projects provide
// objects in an uninitialized state, thus the passed discount
// is exported with lib.Export[Discount](discount)
func (p *Product) AddNewDiscountByCopy(discount Discount) error {
	p.invocationDepth++
	if p.isInitialized && p.invocationDepth == 1 {
		_libError := lib.GetStub(p.Id, p)
		if _libError != nil {
			p.invocationDepth--
			return _libError
		}
	}

	timeFrom, err := discount.GetValidFrom()
	if err != nil {
		p.invocationDepth--
		return err
	}
	if timeFrom.IsZero() {
		if err := discount.SetValidFrom(time.Now()); err != nil {
			p.invocationDepth--
			return err
		}
	}

	exportedDiscount, err := lib.Export[Discount](discount)
	if err != nil {
		p.invocationDepth--
		return err
	}

	p.Discount = append(p.Discount, exportedDiscount.Id)
	_libUpsertError := p.saveChangesIfInitialized()
	p.invocationDepth--
	return _libUpsertError
}

// Example of a method accepting a reference to a Nobject as an input parameter.
// References always refer to initialized objects, hence there is no need
// to export the passed object as in the method 'AddNewDiscountByCopy'
func (p *Product) AddNewDiscountByReference(discount lib.Reference[Discount]) error {
	p.invocationDepth++
	if p.isInitialized && p.invocationDepth == 1 {
		_libError := lib.GetStub(p.Id, p)
		if _libError != nil {
			p.invocationDepth--
			return _libError
		}
	}
	discountInitialized, err := discount.Get()
	if err != nil {
		p.invocationDepth--
		return err
	}
	timeFrom, err := discountInitialized.GetValidFrom()
	if err != nil {
		p.invocationDepth--
		return err
	}
	if timeFrom.IsZero() {
		if err := discountInitialized.SetValidFrom(time.Now()); err != nil {
			p.invocationDepth--
			return err
		}
	}

	p.Discount = append(p.Discount, discountInitialized.Id)
	_libUpsertError := p.saveChangesIfInitialized()
	p.invocationDepth--
	return _libUpsertError
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

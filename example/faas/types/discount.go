package types

import "github.com/Astenna/Nubes/lib"

type Discount struct {
	Id              string
	Percentage      string
	isInitialized   bool
	invocationDepth int
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

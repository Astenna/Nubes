package lib

import "fmt"

type HandlerParameters struct {
	// Indicates the ID of associated object instance
	Id string
	// Indicates the ID of associated object instance
	TypeName string
	// Parameter of the orginal function
	// from which the handler is generated
	Parameter interface{}
}

type AddToManyToManyParam struct {
	TypeName      string
	NewId         string
	OwnerTypeName string
	OwnerId       string
	UsesIndex     bool
}

func (a AddToManyToManyParam) Verify() error {
	if a.TypeName == "" {
		return fmt.Errorf("missing TypeName")
	}
	if a.NewId == "" {
		return fmt.Errorf("missing TypeName")
	}
	if a.OwnerTypeName == "" {
		return fmt.Errorf("missing TypeName")
	}
	if a.OwnerId == "" {
		return fmt.Errorf("missing TypeName")
	}
	return nil
}

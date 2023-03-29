package lib

import "fmt"

type HandlerParameters struct {
	// Indicates the ID of associated object instance
	Id string
	// Indicates the TypeName of associated object instance
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
		return fmt.Errorf("missing NewId")
	}
	if a.OwnerTypeName == "" {
		return fmt.Errorf("missing OwnerTypeName")
	}
	if a.OwnerId == "" {
		return fmt.Errorf("missing OwnerId")
	}
	return nil
}

type LoadBatchParam struct {
	TypeName string
	Ids      []string
}

func (l LoadBatchParam) Verify() error {
	if l.TypeName == "" {
		return fmt.Errorf("missing TypeName")
	}
	if l.Ids == nil || len(l.Ids) == 0 {
		return fmt.Errorf("missing Ids")
	}

	return nil
}

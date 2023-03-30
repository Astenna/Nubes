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

type GetStateParam struct {
	Id        string
	FieldName string
	TypeName  string
	GetStub   bool
}

func (s GetStateParam) Validate() error {
	if s.Id == "" {
		return fmt.Errorf("missing Id of object's field to get")
	}
	if s.TypeName == "" {
		return fmt.Errorf("missing TypeName of object's field to get")
	}

	return nil
}

type SetFieldParam struct {
	Id        string
	FieldName string
	TypeName  string
	Value     interface{}
}

func (s SetFieldParam) Validate() error {
	if s.Id == "" {
		return fmt.Errorf("missing Id of object's field to get")
	}
	if s.FieldName == "" {
		return fmt.Errorf("missing FieldName of object's field to get")
	}
	if s.TypeName == "" {
		return fmt.Errorf("missing TypeName of object's field to get")
	}

	return nil
}

type QueryByPartitionKeyParam struct {
	TableName               string
	PartitionAttributeName  string
	PatritionAttributeValue string
	OutputAttributeName     string
}

func (q QueryByPartitionKeyParam) Validate() error {
	if q.TableName == "" {
		return fmt.Errorf("missing TableName")
	}
	if q.PartitionAttributeName == "" {
		return fmt.Errorf("missing PartitionAttributeName")
	}
	if q.PatritionAttributeValue == "" {
		return fmt.Errorf("missing PatritionAttributeValue")
	}
	if q.OutputAttributeName == "" {
		return fmt.Errorf("missing OutputAttributeName")
	}
	return nil
}

type QueryByIndexParam struct {
	TableName           string
	IndexName           string
	KeyAttributeName    string
	KeyAttributeValue   string
	OutputAttributeName string
}

func (q QueryByIndexParam) Validate() error {
	if q.TableName == "" {
		return fmt.Errorf("missing TableName")
	}
	if q.IndexName == "" {
		return fmt.Errorf("missing IndexName")
	}
	if q.KeyAttributeName == "" {
		return fmt.Errorf("missing KeyAttributeName")
	}
	if q.KeyAttributeValue == "" {
		return fmt.Errorf("missing KeyAttributeValue")
	}
	if q.OutputAttributeName == "" {
		return fmt.Errorf("missing OutputAttributeName")
	}
	return nil
}

type IsInstanceAlreadyCreatedParam struct {
	Id       string
	TypeName string
}

type GetBatchParam struct {
	Ids      []string
	TypeName string
}

func (q GetBatchParam) Validate() error {
	if q.Ids == nil || len(q.Ids) < 1 {
		return fmt.Errorf("empty list of Ids to retrieve")
	}
	if q.TypeName == "" {
		return fmt.Errorf("missing TypeName")
	}
	return nil
}

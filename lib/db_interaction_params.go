package lib

import "fmt"

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

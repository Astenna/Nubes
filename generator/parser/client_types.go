package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
)

type ClientTypesParser struct {
	DefinedTypes          map[string]*StructTypeDefinition
	CustomCtorDefinitions []CustomCtorDefinition
	OtherDecls            OtherDecls
	tokenSet              *token.FileSet
	packs                 map[string]*ast.Package
}

type StructTypeDefinition struct {
	PackageName             string
	Imports                 string
	NobjectImplementation   string
	CustomIdImplementation  string
	CustomIdReceiverName    string
	TypeNameLower           string
	TypeNameOrginalCase     string
	MemberFunctions         []MethodDefinition
	FieldDefinitions        []FieldDefinition
	OneToManyRelationships  []OneToManyRelationshipField
	ManyToManyRelationships []ManyToManyRelationshipField
}

type MethodDefinition struct {
	ReceiverName            string
	FuncName                string
	InputParamType          string
	OptionalReturnType      string
	OptionalReturnTypeUpper string
	IsReturnTypeNobject     bool
}

type FieldDefinition struct {
	FieldNameUpper  string
	FieldName       string
	FieldType       string
	FieldTypeUpper  string
	Tags            string
	IsReference     bool
	IsReferenceList bool
	IsReadonly      bool
}

type OtherDecls struct {
	Consts   []string
	GenDecls []string
}

func NewClientTypesParser(path string) (*ClientTypesParser, error) {
	typeSpec := new(ClientTypesParser)
	typeSpec.tokenSet = token.NewFileSet()
	packs, err := parser.ParseDir(typeSpec.tokenSet, path, nil, 0)
	if err != nil {
		return nil, err
	}

	typeSpec.packs = packs
	typeSpec.DefinedTypes = make(map[string]*StructTypeDefinition)
	typeSpec.OtherDecls = OtherDecls{}
	return typeSpec, nil
}

func (t *ClientTypesParser) Run() {
	t.detectGenDecls()
	t.detectFuncs()
	t.detectAndSetNobjectsReturnTypes()
}

func (t *ClientTypesParser) detectAndSetNobjectsReturnTypes() {
	for _, typeDefinition := range t.DefinedTypes {
		for i, function := range typeDefinition.MemberFunctions {
			if isReturnTypeDefined(function) && isReturnTypeNobject(function, t.DefinedTypes) {
				typeDefinition.MemberFunctions[i].IsReturnTypeNobject = true
				typeDefinition.MemberFunctions[i].OptionalReturnTypeUpper = function.OptionalReturnType
				typeDefinition.MemberFunctions[i].OptionalReturnType = lowerCasedFirstChar(function.OptionalReturnType)
			}
			if isOptionalParamNobject(function, t.DefinedTypes) {
				typeDefinition.MemberFunctions[i].InputParamType = lowerCasedFirstChar(function.InputParamType)
			}
		}
	}
}

func isOptionalParamNobject(f MethodDefinition, defTypes map[string]*StructTypeDefinition) bool {
	return f.InputParamType != "" && defTypes[f.InputParamType] != nil && defTypes[f.InputParamType].NobjectImplementation != ""
}

func isReturnTypeDefined(f MethodDefinition) bool {
	return f.OptionalReturnType != ""
}

func isReturnTypeNobject(f MethodDefinition, defTypes map[string]*StructTypeDefinition) bool {
	return defTypes[f.OptionalReturnType] != nil && defTypes[f.OptionalReturnType].NobjectImplementation != ""
}

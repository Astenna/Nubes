package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"strings"
)

type ClientTypesParser struct {
	DefinedTypes          map[string]*StructTypeDefinition
	CustomCtorDefinitions []CustomCtorDefinition
	OtherDecls            OtherDecls
	// functions field temporary stores functions (not methods)
	// that do not belong to any type (have no receiver type)
	functions []*ast.FuncDecl
	tokenSet  *token.FileSet
	packs     map[string]*ast.Package
}

type StructTypeDefinition struct {
	PackageName             string
	Imports                 string
	NobjectImplementation   string
	CustomIdFieldName       string
	TypeNameLower           string
	TypeNameOrginalCase     string
	MemberFunctions         []MethodDefinition
	FieldDefinitions        []FieldDefinition
	OneToManyRelationships  []OneToManyRelationshipField
	ManyToManyRelationships []ManyToManyRelationshipField
	CustomExportInputType   string
	CustomDeleteInputType   string
}

type MethodDefinition struct {
	ReceiverName            string
	FuncName                string
	InputParamType          string
	OptionalReturnType      string
	OptionalReturnTypeUpper string
	IsReturnTypeNobject     bool
	IsReturnTypeList        bool
	IsInputParamNobject     bool
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
	t.detectCustomImplementations()
}

func (t *ClientTypesParser) detectAndSetNobjectsReturnTypes() {
	for _, typeDefinition := range t.DefinedTypes {
		for i, function := range typeDefinition.MemberFunctions {
			adjustSubtypesIfInputOrOuputParamsAreReferences(&typeDefinition.MemberFunctions[i])

			if isReturnTypeDefined(function) && isNobject(function.OptionalReturnType, t.DefinedTypes) {
				typeDefinition.MemberFunctions[i].IsReturnTypeNobject = true
				typeDefinition.MemberFunctions[i].OptionalReturnType = function.OptionalReturnType + "Stub"
			}
			if isReturnTypeDefined(function) && isReturnTypeList(function.OptionalReturnType) {
				typeDefinition.MemberFunctions[i].IsReturnTypeList = true
			}
			if isOptionalParamNobject(function, t.DefinedTypes) {
				typeDefinition.MemberFunctions[i].IsInputParamNobject = true
				typeDefinition.MemberFunctions[i].InputParamType = function.InputParamType
			}
		}
	}
}

func (t *ClientTypesParser) detectCustomImplementations() {
	for _, fn := range t.functions {
		if strings.HasPrefix(fn.Name.Name, CustomExportPrefix) {
			typeName := strings.TrimPrefix(fn.Name.Name, CustomExportPrefix)

			if !isNobject(typeName, t.DefinedTypes) {
				continue
			}

			param, err := getFunctionParm(fn.Type.Params, t.DefinedTypes)
			if err != nil {
				fmt.Println("Maximum allowed number of parameters is 1. Custom export generation for " + fn.Name.Name + "skipped")
				continue
			}
			if isNobject(param, t.DefinedTypes) {
				param = param + "Stub"
			}

			t.DefinedTypes[typeName].CustomExportInputType = param

		} else if strings.HasPrefix(fn.Name.Name, CustomDeletePrefix) {
			typeName := strings.TrimPrefix(fn.Name.Name, CustomDeletePrefix)

			if !isNobject(typeName, t.DefinedTypes) {
				continue
			}

			param, err := getFunctionParm(fn.Type.Params, t.DefinedTypes)
			if err != nil {
				fmt.Println("Maximum allowed number of parameters is 1. Custom delete generation for " + fn.Name.Name + "skipped")
				continue
			}
			if isNobject(param, t.DefinedTypes) {
				param = param + "Stub"
			}

			t.DefinedTypes[typeName].CustomDeleteInputType = param

		} else if strings.HasPrefix(fn.Name.Name, ConstructorPrefix) {
			typeName := strings.TrimPrefix(fn.Name.Name, ConstructorPrefix)
			if !isNobject(typeName, t.DefinedTypes) {
				continue
			}
			param, err := getFunctionParm(fn.Type.Params, t.DefinedTypes)
			if err != nil {
				fmt.Println("Maximum allowed number of parameters is 1. Custom constructor generation for " + fn.Name.Name + "skipped")
				continue
			}
			t.CustomCtorDefinitions = append(t.CustomCtorDefinitions, CustomCtorDefinition{
				TypeName:               typeName,
				OptionalParamType:      param,
				IsOptionalParamNobject: isNobject(param, t.DefinedTypes),
			})
		}
	}
}

func isOptionalParamNobject(f MethodDefinition, defTypes map[string]*StructTypeDefinition) bool {
	return f.InputParamType != "" && defTypes[f.InputParamType] != nil && defTypes[f.InputParamType].NobjectImplementation != ""
}

func isReturnTypeDefined(f MethodDefinition) bool {
	return f.OptionalReturnType != ""
}

func isNobject(typeName string, defTypes map[string]*StructTypeDefinition) bool {
	trimmedListPrefix := strings.TrimPrefix(typeName, "[]")
	return defTypes[trimmedListPrefix] != nil && defTypes[trimmedListPrefix].NobjectImplementation != ""
}

func isReturnTypeList(typeName string) bool {
	return strings.HasPrefix(typeName, "[]")
}

func getFunctionParm(params *ast.FieldList, definedStructs map[string]*StructTypeDefinition) (string, error) {
	if params.List == nil || len(params.List) == 0 {
		return "", nil
	} else if len(params.List) > 1 {
		return "", fmt.Errorf("maximum allowed number of parameters is 1")
	}

	inputParamType := types.ExprString(params.List[0].Type)
	return inputParamType, nil
}

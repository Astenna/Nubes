package parser

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/exp/maps"
)

type TypeDefinition struct {
	PackageName            string
	Imports                string
	NobjectImplementation  string
	CustomIdImplementation string
	CustomIdReceiverName   string
	TypeNameLower          string
	TypeNameUpper          string
	MemberFunctions        []MemberFunction
	FieldDefinitions       []FieldDefinition
}

type MemberFunction struct {
	ReceiverName            string
	FuncName                string
	InputParamType          string
	OptionalReturnType      string
	OptionalReturnTypeUpper string
	IsReturnTypeNobject     bool
}

type FieldDefinition struct {
	FieldNameUpper string
	FieldName      string
	FieldType      string
	FieldTypeUpper string
	Tags           string
	IsReference    bool
	IsReadonly     bool
}

type OtherDecls struct {
	Consts      []string
	GenDecls    []string
	PackageName string
}

func PrepareTypes(path string) ([]*TypeDefinition, OtherDecls) {
	set := token.NewFileSet()
	packs, err := parser.ParseDir(set, path, nil, 0)
	assertDirParsed(err)

	definedTypes := make(map[string]*TypeDefinition)
	otherTypesDecls := OtherDecls{}

	for _, pack := range packs {
		for _, f := range pack.Files {

			ast.Inspect(f, func(n ast.Node) bool {
				if typeSpec, ok := n.(*ast.TypeSpec); ok {
					if strctType, ok := typeSpec.Type.(*ast.StructType); ok {
						typeName := strings.TrimPrefix(typeSpec.Name.Name, "*")
						MakeFieldsUnexported(strctType.Fields)
						if _, ok := definedTypes[typeName]; !ok {
							definedTypes[typeName] = &TypeDefinition{}
						}
						definedTypes[typeName].TypeNameUpper = typeName
						definedTypes[typeName].TypeNameLower = MakeFirstCharacterLowerCase(typeName)
						definedTypes[typeName].FieldDefinitions = GetFieldDefinitions(typeName, strctType)
					} else {
						def, err := GetTypeSpecAsString(set, typeSpec)
						if err != nil {
							fmt.Println(err)
						} else {
							otherTypesDecls.GenDecls = append(otherTypesDecls.GenDecls, def)
						}
					}
				}
				return true
			})

			for _, d := range f.Decls {

				if genDecl, ok := d.(*ast.GenDecl); ok {
					if genDecl.Tok == token.CONST {
						constStr, err := GetConstAsString(set, genDecl)
						if err != nil {
							fmt.Println(err)
						}
						otherTypesDecls.Consts = append(otherTypesDecls.Consts, constStr)
						continue
					}
				}

				if fn, isFn := d.(*ast.FuncDecl); isFn {
					if fn.Recv == nil {
						continue
					}

					typeName := strings.TrimPrefix(types.ExprString(fn.Recv.List[0].Type), "*")
					if fn.Name.Name == NobjectImplementationMethod {
						funcString, err := getFunctionBodyAsString(set, fn.Body)
						if err != nil {
							fmt.Println("error occurred when parsing GetTypeName of " + typeName)
							continue
						}

						if _, ok := definedTypes[typeName]; !ok {
							definedTypes[typeName] = &TypeDefinition{}
						}
						definedTypes[typeName].NobjectImplementation = funcString
						continue
					}

					if fn.Name.Name == CustomIdImplementationMethod {
						funcString, err := getFunctionBodyAsString(set, fn.Body)
						if err != nil {
							fmt.Println("error occurred when parsing GetTypeName of " + typeName)
							continue
						}

						if _, ok := definedTypes[typeName]; !ok {
							definedTypes[typeName] = &TypeDefinition{}
						}
						definedTypes[typeName].CustomIdImplementation = funcString
						if len(fn.Recv.List[0].Names) > 0 {
							definedTypes[typeName].CustomIdReceiverName = fn.Recv.List[0].Names[0].Name
						}
						continue
					}

					memberFunction, err := PrepareMemberFunction(fn)
					if err != nil {
						fmt.Println("Function "+fn.Name.Name+"not generated in client lib", err)
						continue
					}

					if _, ok := definedTypes[typeName]; !ok {
						definedTypes[typeName] = &TypeDefinition{}
					}
					definedTypes[typeName].MemberFunctions = append(definedTypes[typeName].MemberFunctions, *memberFunction)
				}
			}
		}
	}

	DetectAndSetNobjectsReturnTypes(definedTypes)
	return maps.Values(definedTypes), otherTypesDecls
}

func DetectAndSetNobjectsReturnTypes(definedTypes map[string]*TypeDefinition) {
	for _, typeDefinition := range definedTypes {
		for i, function := range typeDefinition.MemberFunctions {
			if isReturnTypeDefined(function) && isReturnTypeNobject(function, definedTypes) {
				typeDefinition.MemberFunctions[i].IsReturnTypeNobject = true
				typeDefinition.MemberFunctions[i].OptionalReturnTypeUpper = function.OptionalReturnType
				typeDefinition.MemberFunctions[i].OptionalReturnType = MakeFirstCharacterLowerCase(function.OptionalReturnType)
			}
		}
	}
}

func isReturnTypeDefined(f MemberFunction) bool {
	return f.OptionalReturnType != ""
}

func isReturnTypeNobject(f MemberFunction, defTypes map[string]*TypeDefinition) bool {
	return defTypes[f.OptionalReturnType] != nil && defTypes[f.OptionalReturnType].NobjectImplementation != ""
}

func GetFieldDefinitions(typeName string, strctType *ast.StructType) []FieldDefinition {
	fieldDefinitions := make([]FieldDefinition, 0, len(strctType.Fields.List)-1)

	for _, field := range strctType.Fields.List {
		newFieldDefinition := FieldDefinition{
			FieldNameUpper: MakeFirstCharacterUpperCase(field.Names[0].Name),
			FieldName:      field.Names[0].Name,
		}

		if field.Names[0].Name == "id" {
			newFieldDefinition.IsReadonly = true
		}

		newFieldDefinition.FieldType = strings.TrimPrefix(types.ExprString(field.Type), "*")
		if strings.Contains(newFieldDefinition.FieldType, ReferenceType) {
			newFieldDefinition.FieldType = strings.TrimPrefix(newFieldDefinition.FieldType, ReferenceType)
			newFieldDefinition.FieldTypeUpper = strings.Trim(newFieldDefinition.FieldType, "[]")
			newFieldDefinition.FieldType = MakeFirstCharacterLowerCase(newFieldDefinition.FieldTypeUpper)
			newFieldDefinition.IsReference = true
		}

		if field.Tag != nil && strings.Contains(field.Tag.Value, ReadonlyTag) {
			newFieldDefinition.IsReadonly = true
			newFieldDefinition.Tags = field.Tag.Value
		}

		fieldDefinitions = append(fieldDefinitions, newFieldDefinition)
	}

	return fieldDefinitions
}

func PrepareMemberFunction(fn *ast.FuncDecl) (*MemberFunction, error) {

	if fn.Type.Results == nil ||
		types.ExprString(fn.Type.Results.List[len(fn.Type.Results.List)-1].Type) != "error" {
		return nil, fmt.Errorf("methods belonging to nobjects must return error type")
	}
	if len(fn.Type.Results.List) > 2 {
		return nil, fmt.Errorf("methods belonging to nobjects can return at most 2 variables")
	}

	memberFunction := MemberFunction{
		FuncName: fn.Name.Name,
	}

	if len(fn.Recv.List[0].Names) > 0 {
		memberFunction.ReceiverName = fn.Recv.List[0].Names[0].Name
	}

	if len(fn.Type.Results.List) > 1 {
		memberFunction.OptionalReturnType = types.ExprString(fn.Type.Results.List[0].Type)
	}

	if len(fn.Type.Params.List) > 1 {
		return nil, fmt.Errorf("methods belonging to nobjects can have at most 1 parameter")
	} else if len(fn.Type.Params.List) == 1 {
		memberFunction.InputParamType = types.ExprString(fn.Type.Params.List[0].Type)
	}

	return &memberFunction, nil
}

func MakeFieldsUnexported(fieldList *ast.FieldList) {
	for _, field := range fieldList.List {
		field.Names[0].Name = MakeFirstCharacterLowerCase(field.Names[0].Name)
	}
}

func MakeFirstCharacterLowerCase(str string) string {
	if len(str) < 2 {
		return strings.ToLower(str)
	}

	bts := []byte(str)

	firstByte := bytes.ToLower([]byte{bts[0]})
	rest := bts[1:]

	str = string(bytes.Join([][]byte{firstByte, rest}, nil))
	return str
}

func MakeFirstCharacterUpperCase(str string) string {
	if len(str) < 2 {
		return strings.ToUpper(str)
	}

	bts := []byte(str)

	firstByte := bytes.ToUpper([]byte{bts[0]})
	rest := bts[1:]

	str = string(bytes.Join([][]byte{firstByte, rest}, nil))
	return str
}

func GetTypeSpecAsString(fset *token.FileSet, detectedStruct *ast.TypeSpec) (string, error) {
	var buf bytes.Buffer
	buf.WriteString("type ")
	err := printer.Fprint(&buf, fset, detectedStruct)
	if err != nil {
		return "", fmt.Errorf("error occurred when parsing the struct")
	}
	return buf.String(), nil
}

func GetConstAsString(fset *token.FileSet, detectedGenDecl *ast.GenDecl) (string, error) {
	var buf bytes.Buffer
	err := printer.Fprint(&buf, fset, detectedGenDecl)
	if err != nil {
		return "", fmt.Errorf("error occurred when parsing the struct")
	}
	return buf.String(), nil
}

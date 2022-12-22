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
	PackageName           string
	Imports               string
	StructDefinition      string
	NobjectImplementation string
	TypeNameLower         string
	TypeNameUpper         string
	MemberFunctions       []MemberFunction
	FieldDefinitions      []FieldDefinition
}

type MemberFunction struct {
	ReceiverName       string
	FuncName           string
	InputParamType     string
	OptionalReturnType string
}

type FieldDefinition struct {
	FieldNameUpper string
	FieldName      string
	FieldType      string
	IsReference    bool
}

func PrepareTypes(path string) []*TypeDefinition {
	set := token.NewFileSet()
	packs, err := parser.ParseDir(set, path, nil, 0)
	AssertDirParsed(err)

	definedTypes := make(map[string]*TypeDefinition)

	for _, pack := range packs {
		for _, f := range pack.Files {

			ast.Inspect(f, func(n ast.Node) bool {
				if typeSpec, ok := n.(*ast.TypeSpec); ok {
					if strctType, ok := typeSpec.Type.(*ast.StructType); ok {
						typeName := strings.TrimPrefix(typeSpec.Name.Name, "*")
						MakeFieldsUnexported(strctType.Fields)
						structString, err := GetStructAsString(set, typeSpec)
						if err == nil {
							if _, ok := definedTypes[typeName]; !ok {
								definedTypes[typeName] = &TypeDefinition{}
							}
							definedTypes[typeName].StructDefinition = structString
							definedTypes[typeName].TypeNameUpper = typeName
							definedTypes[typeName].TypeNameLower = MakeFirstCharacterLowerCase(typeName)
							definedTypes[typeName].FieldDefinitions = GetFieldDefinitions(typeName, strctType)
						}
					}
				}
				return true
			})

			for _, d := range f.Decls {
				if fn, isFn := d.(*ast.FuncDecl); isFn {
					if fn.Recv == nil {
						continue
					}

					typeName := strings.TrimPrefix(types.ExprString(fn.Recv.List[0].Type), "*")
					if fn.Name.Name == GetTypeName {
						funcString, err := GetFuncDeclAsString(set, fn)
						if err != nil {
							fmt.Println("error occurred when parsing GetTypeName of " + typeName)
							continue
						}

						if elem, ok := definedTypes[typeName]; !ok {
							definedTypes[typeName] = &TypeDefinition{
								NobjectImplementation: funcString,
							}
						} else {
							elem.NobjectImplementation = funcString
						}
						continue
					}

					memberFunction, err := PrepareMemberFunction(fn)
					if err != nil {
						fmt.Println("Function "+fn.Name.Name+"not generated in client lib", err)
						continue
					}

					if elem, ok := definedTypes[typeName]; !ok {
						definedTypes[typeName] = &TypeDefinition{
							MemberFunctions: []MemberFunction{
								*memberFunction,
							},
						}
					} else {
						elem.MemberFunctions = append(elem.MemberFunctions, *memberFunction)
					}
				}
			}
		}
	}

	return maps.Values(definedTypes)
}

func GetFieldDefinitions(typeName string, strctType *ast.StructType) []FieldDefinition {
	fieldDefinitions := make([]FieldDefinition, 0, len(strctType.Fields.List)-1)

	for _, field := range strctType.Fields.List {
		if field.Names[0].Name != "id" {

			newFieldDefinition := FieldDefinition{
				FieldNameUpper: MakeFirstCharacterUpperCase(field.Names[0].Name),
				FieldName:      field.Names[0].Name,
			}

			newFieldDefinition.FieldType = strings.TrimPrefix(types.ExprString(field.Type), "*")
			if strings.Contains(newFieldDefinition.FieldType, ReferenceType) {
				newFieldDefinition.FieldType = strings.TrimPrefix(newFieldDefinition.FieldType, ReferenceType)
				newFieldDefinition.FieldType = strings.Trim(newFieldDefinition.FieldType, "[]")
				newFieldDefinition.IsReference = true
			}

			fieldDefinitions = append(fieldDefinitions, newFieldDefinition)
		}
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

func GetStructAsString(fset *token.FileSet, detectedStruct *ast.TypeSpec) (string, error) {
	var buf bytes.Buffer
	err := printer.Fprint(&buf, fset, detectedStruct)
	if err != nil {
		return "", fmt.Errorf("error occurred when parsing the struct")
	}
	return buf.String(), nil
}

func GetFuncDeclAsString(fset *token.FileSet, f *ast.FuncDecl) (string, error) {
	var buf bytes.Buffer
	err := printer.Fprint(&buf, fset, f)
	if err != nil {
		return "", fmt.Errorf("error occurred when parsing the function body")
	}
	return buf.String(), nil
}

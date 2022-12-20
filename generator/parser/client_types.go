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
)

type TypeDefinition struct {
	PackageName      string
	Imports          string
	StructDefinition string
	MemberFunctions  []MemberFunction
}

type MemberFunction struct {
	ReceiverName       string
	ReceiverType       string
	FuncName           string
	InputParamType     string
	OptionalReturnType string
}

func PrepareTypes(path string) map[string]*TypeDefinition {
	set := token.NewFileSet()
	packs, err := parser.ParseDir(set, path, nil, 0)
	AssertDirParsed(err)

	var detectedTypes []*ast.StructType
	typeFiles := make(map[string]*TypeDefinition)

	for _, pack := range packs {
		for _, f := range pack.Files {

			ast.Inspect(f, func(n ast.Node) bool {
				if typeSpec, ok := n.(*ast.TypeSpec); ok {
					if strctType, ok := typeSpec.Type.(*ast.StructType); ok {
						MakeFieldsUnexported(strctType.Fields)
						detectedTypes = append(detectedTypes, strctType)

						structString, err := GetStructAsString(set, typeSpec)
						if err == nil {
							typeName := strings.TrimPrefix(typeSpec.Name.Name, "*")
							if elem, ok := typeFiles[typeName]; !ok {
								typeFiles[typeName] = &TypeDefinition{
									StructDefinition: structString,
								}
							} else {
								elem.StructDefinition = structString
							}
						}
					}
				}
				return true
			})

			for _, d := range f.Decls {
				if fn, isFn := d.(*ast.FuncDecl); isFn {
					if fn.Recv == nil || fn.Name.Name == GetTypeName {
						continue
					}

					memberFunction, err := PrepareMemberFunction(fn)
					if err != nil {
						fmt.Println("Function "+fn.Name.Name+"not generated in client lib", err)
						continue
					}

					if elem, ok := typeFiles[memberFunction.ReceiverType]; !ok {
						typeFiles[memberFunction.ReceiverType] = &TypeDefinition{
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

	return typeFiles
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
		ReceiverType: strings.TrimPrefix(types.ExprString(fn.Recv.List[0].Type), "*"),
		FuncName:     fn.Name.Name,
	}

	if len(fn.Recv.List[0].Names) > 0 {
		memberFunction.ReceiverName = fn.Recv.List[0].Names[0].Name
	}

	if len(fn.Type.Results.List) > 1 {
		memberFunction.InputParamType = types.ExprString(fn.Type.Results.List[0].Type)
	}

	return &memberFunction, nil
}

func MakeFieldsUnexported(fieldList *ast.FieldList) {
	for _, field := range fieldList.List {
		field.Names[0].Name = strings.ToLower(field.Names[0].Name)
	}
}

func GetStructAsString(fset *token.FileSet, detectedStruct *ast.TypeSpec) (string, error) {
	var buf bytes.Buffer
	err := printer.Fprint(&buf, fset, detectedStruct)
	if err != nil {
		return "", fmt.Errorf("error occurred when parsing the struct")
	}
	return buf.String(), nil
}

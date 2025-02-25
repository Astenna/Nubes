package parser

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"go/types"
	"strings"
)

func (t *ClientTypesParser) detectFuncs() {
	for _, pack := range t.packs {
		for _, f := range pack.Files {
			for _, decl := range f.Decls {
				if fn, isFn := decl.(*ast.FuncDecl); isFn {

					// ignore unexported functions (i.e. starting with lowercase letter)
					if fn.Name.IsExported() {
						if fn.Name.Name == InitFunctionName {
							continue
						}

						if fn.Recv == nil {
							t.functions = append(t.functions, fn)
							continue
						}

						typeName := strings.TrimPrefix(types.ExprString(fn.Recv.List[0].Type), "*")
						if fn.Name.Name == NobjectImplementationMethod {
							funcString, err := getFunctionBodyAsString(t.tokenSet, fn.Body)
							if err != nil {
								fmt.Println("error occurred when parsing GetTypeName of " + typeName)
								continue
							}

							if _, ok := t.DefinedTypes[typeName]; !ok {
								t.DefinedTypes[typeName] = &StructTypeDefinition{}
							}
							t.DefinedTypes[typeName].NobjectImplementation = funcString
							continue
						}

						if fn.Name.Name == CustomIdImplementationMethod {

							idFieldName, err := getIdFieldNameFromCustomIdImpl(fn)
							if err != nil {
								fmt.Println(err)
								continue
							}
							t.DefinedTypes[typeName].CustomIdFieldName = idFieldName
							continue
						}

						if isGetterOrSetterMethod(fn, typeName, t.DefinedTypes) {
							continue
						}

						// at this point, the method is recognized as a general,
						// state-changing method
						memberFunction, err := parseMethod(fn)
						if err != nil {
							fmt.Println("Function "+fn.Name.Name+"not generated in client lib", err)
							continue
						}

						if _, ok := t.DefinedTypes[typeName]; !ok {
							t.DefinedTypes[typeName] = &StructTypeDefinition{}
						}
						t.DefinedTypes[typeName].MemberFunctions = append(t.DefinedTypes[typeName].MemberFunctions, *memberFunction)
					}
				}
			}
		}
	}
}

func getFunctionBodyAsString(fset *token.FileSet, body *ast.BlockStmt) (string, error) {
	var buf bytes.Buffer
	err := printer.Fprint(&buf, fset, body)
	if err != nil {
		return "", fmt.Errorf("error occurred when parsing the function body")
	}
	return buf.String(), nil
}

func parseMethod(fn *ast.FuncDecl) (*MethodDefinition, error) {

	if fn.Type.Results == nil ||
		types.ExprString(fn.Type.Results.List[len(fn.Type.Results.List)-1].Type) != "error" {
		return nil, fmt.Errorf("methods belonging to nobjects must return error type")
	}
	if len(fn.Type.Results.List) > 2 {
		return nil, fmt.Errorf("methods belonging to nobjects can return at most 2 variables")
	}

	memberFunction := MethodDefinition{
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

// adjustSubtypesIfInputOrOuputParamsAreReferences changes type specification
// of generics reference types in input and output parameters so that
// initialized subtypes are used e.g. Reference<User> -> Reference<user>
// or ReferenceList<User> -> ReferenceList<user>

func adjustSubtypesIfInputOrOuputParamsAreReferences(methodDefinition *MethodDefinition) {
	if strings.Contains(methodDefinition.InputParamType, ReferenceListType) {
		methodDefinition.InputParamType = strings.TrimPrefix(methodDefinition.InputParamType, ReferenceListType)
		methodDefinition.InputParamType = strings.Trim(methodDefinition.InputParamType, "[]")
		methodDefinition.InputParamType = ReferenceListType + "[" + lowerCasedFirstChar(methodDefinition.InputParamType) + "]"
	} else if strings.Contains(methodDefinition.InputParamType, ReferenceType) {
		methodDefinition.InputParamType = strings.TrimPrefix(methodDefinition.InputParamType, ReferenceType)
		methodDefinition.InputParamType = strings.Trim(methodDefinition.InputParamType, "[]")
		methodDefinition.InputParamType = ReferenceType + "[" + lowerCasedFirstChar(methodDefinition.InputParamType) + "]"
	}

	if strings.Contains(methodDefinition.OptionalReturnType, ReferenceListType) {
		methodDefinition.OptionalReturnType = strings.TrimPrefix(methodDefinition.OptionalReturnType, ReferenceListType)
		methodDefinition.OptionalReturnType = strings.Trim(methodDefinition.OptionalReturnType, "[]")
		methodDefinition.OptionalReturnType = ReferenceListType + "[" + lowerCasedFirstChar(methodDefinition.OptionalReturnType) + "]"
	} else if strings.Contains(methodDefinition.OptionalReturnType, ReferenceType) {
		methodDefinition.OptionalReturnType = strings.TrimPrefix(methodDefinition.OptionalReturnType, ReferenceType)
		methodDefinition.OptionalReturnType = strings.Trim(methodDefinition.OptionalReturnType, "[]")
		methodDefinition.OptionalReturnType = ReferenceType + "[" + lowerCasedFirstChar(methodDefinition.OptionalReturnType) + "]"
	}
}

func getIdFieldNameFromCustomIdImpl(fn *ast.FuncDecl) (string, error) {
	var returnResult string
	ast.Inspect(fn, func(n ast.Node) bool {
		if ret, ok := n.(*ast.ReturnStmt); ok {
			returnResult = types.ExprString(ret.Results[0])
			return false
		}
		return true
	})

	if returnResult == "" {
		return "", errors.New("unable to detect Id field based on custom id interface implementation for" + fn.Recv.List[0].Names[0].Name)
	}
	splitted := strings.Split(returnResult, ".")
	return splitted[len(splitted)-1], nil
}

func isGetterOrSetterMethod(fn *ast.FuncDecl, typeName string, definedTypes map[string]*StructTypeDefinition) bool {
	if strings.HasPrefix(fn.Name.Name, GetPrefix) {
		fieldName := strings.TrimPrefix(fn.Name.Name, GetPrefix)

		for _, field := range definedTypes[typeName].FieldDefinitions {
			if field.FieldNameUpper == fieldName {
				return true
			}
		}
		for _, field := range definedTypes[typeName].OneToManyRelationships {
			if field.FromFieldNameUpper == fieldName || field.FromFieldName == fieldName {
				return true
			}
		}
		for _, field := range definedTypes[typeName].ManyToManyRelationships {
			if field.FromFieldNameUpper == fieldName || field.FromFieldName == fieldName {
				return true
			}
		}
	} else if strings.HasPrefix(fn.Name.Name, SetPrefix) {
		fieldName := strings.TrimPrefix(fn.Name.Name, SetPrefix)

		for _, field := range definedTypes[typeName].FieldDefinitions {
			if field.FieldNameUpper == fieldName {
				return true
			}
		}
		for _, field := range definedTypes[typeName].OneToManyRelationships {
			if field.FromFieldNameUpper == fieldName || field.FromFieldName == fieldName {
				return true
			}
		}
		for _, field := range definedTypes[typeName].ManyToManyRelationships {
			if field.FromFieldNameUpper == fieldName || field.FromFieldName == fieldName {
				return true
			}
		}
	}

	return false
}

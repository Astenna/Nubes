package parser

import (
	"errors"
	"fmt"
	"go/ast"
	"go/types"
	"strings"
)

func (t *ClientTypesParser) detectFuncs() {
	for _, pack := range t.packs {
		for _, f := range pack.Files {
			for _, decl := range f.Decls {
				if fn, isFn := decl.(*ast.FuncDecl); isFn {

					if fn.Name.Name == InitFunctionName {
						continue
					}

					if fn.Recv == nil && strings.HasPrefix(fn.Name.Name, ConstructorPrefix) {
						typeName := strings.TrimPrefix(fn.Name.Name, ConstructorPrefix)
						param, isOptitionalParamNobject, err := getFunctionParm(fn.Type.Params, t.DefinedTypes)
						if err != nil {
							fmt.Println("Maximum allowed number of parameters is 1. Handler generation for " + f.Name.Name + "skipped")
							continue
						}
						t.CustomCtorDefinitions = append(t.CustomCtorDefinitions, CustomCtorDefinition{
							TypeName:               typeName,
							OptionalParamType:      param,
							IsOptionalParamNobject: isOptitionalParamNobject,
						})
						continue
					}

					if fn.Recv == nil {
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
						funcString, err := getFunctionBodyAsString(t.tokenSet, fn.Body)
						if err != nil {
							fmt.Println("error occurred when parsing GetTypeName of " + typeName)
							continue
						}

						if _, ok := t.DefinedTypes[typeName]; !ok {
							t.DefinedTypes[typeName] = &StructTypeDefinition{}
						}
						t.DefinedTypes[typeName].CustomIdImplementation = funcString
						if len(fn.Recv.List[0].Names) > 0 {
							t.DefinedTypes[typeName].CustomIdReceiverName = fn.Recv.List[0].Names[0].Name
						}
						continue
					}

					if isGetterOrSetterMethod(fn, typeName, t.DefinedTypes) {
						continue
					}

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

func getFunctionParm(params *ast.FieldList, definedStructs map[string]*StructTypeDefinition) (string, bool, error) {
	if params.List == nil || len(params.List) == 0 {
		return "", false, nil
	} else if len(params.List) > 1 {
		return "", false, fmt.Errorf("maximum allowed number of parameters is 1")
	}

	inputParamType := types.ExprString(params.List[0].Type)
	if definedStructs[inputParamType] != nil && definedStructs[inputParamType].NobjectImplementation == "" {
		return inputParamType, true, nil
	}
	return inputParamType, false, nil
}

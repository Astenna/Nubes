package parser

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"go/types"
	"strings"
)

type StateChangingHandler struct {
	OrginalPackage      string
	OrginalPackageAlias string
	Imports             string
	MethodName          string
	ReceiverType        string
	ReceiverIdFieldName string
	OptionalReturnType  string
	OptionalInputType   string
}

type detectedFunction struct {
	Function *ast.FuncDecl
	Imports  []*ast.ImportSpec
}

func (t *TypeSpecParser) prepareDataForHandlers() {
	fileFunctionsMap := t.detectedFunctions

	for _, functions := range fileFunctionsMap {
		for _, detectedFunction := range functions {
			f := detectedFunction.Function

			if f.Recv == nil {

				param, err := getHandlerInputParam(f.Type.Params, t.Output.TypeFields)
				if err != nil {
					fmt.Println("Maximum allowed number of parameters is 1. Handler generation for " + f.Name.Name + "skipped")
					continue
				}
				if strings.HasPrefix(f.Name.Name, ConstructorPrefix) {
					typeName := strings.TrimPrefix(f.Name.Name, ConstructorPrefix)
					t.CustomCtors = append(t.CustomCtors, CustomCtorDefinition{
						OrginalPackage:      t.Output.ImportPath,
						OrginalPackageAlias: OrginalPackageAlias,
						TypeName:            typeName,
						OptionalParamType:   param,
					})
				} else if strings.Contains(f.Name.Name, CustomExportPrefix) {
					typeName := strings.TrimPrefix(f.Name.Name, CustomExportPrefix)
					if isNobject, isPresent := t.Output.IsNobjectInOrginalPackage[typeName]; !isPresent || !isNobject {
						fmt.Println(`Custom exports must be a concatenation of the 'Export' and valid Nobject type name! " + typeName + " is not a valid nobject, 
						custom export definition in funcition " + f.Name.Name + " skipped`)
						continue
					}

					t.Output.TypesWithCustomExport[typeName] = CustomExportDefinition{
						InputParameterType: param,
					}

				} else if strings.Contains(f.Name.Name, CustomDeletePrefix) {
					typeName := strings.TrimPrefix(f.Name.Name, CustomDeletePrefix)
					if isNobject, isPresent := t.Output.IsNobjectInOrginalPackage[typeName]; !isPresent || !isNobject {
						fmt.Println(`Custom exports must be a concatenation of the 'Export' and valid Nobject type name! " + typeName + " is not a valid nobject, 
						custom export definition in funcition " + f.Name.Name + " skipped`)
						continue
					}

					t.Output.TypesWithCustomDelete[typeName] = CustomDeleteDefinition{
						InputParameterType: param,
					}
				}

			} else {
				receiverTypeName := getFunctionReceiverTypeAsString(f.Recv)

				if strings.HasPrefix(f.Name.Name, SetPrefix) {
					withoutPrefix := strings.TrimPrefix(f.Name.Name, SetPrefix)
					if _, exists := t.Output.TypeFields[receiverTypeName][withoutPrefix]; exists {
						continue
					}
				} else if strings.HasPrefix(f.Name.Name, GetPrefix) {
					withoutPrefix := strings.TrimPrefix(f.Name.Name, GetPrefix)
					if _, exists := t.Output.TypeFields[receiverTypeName][withoutPrefix]; exists {
						continue
					}
				}

				if isNobject, isPresent := t.Output.IsNobjectInOrginalPackage[receiverTypeName]; !isPresent || !isNobject {
					fmt.Println("Member type does not implement Nobject interface. Handler generation for " + f.Name.Name + "skipped")
					continue
				}

				newHandler := StateChangingHandler{
					OrginalPackage:      t.Output.ImportPath,
					OrginalPackageAlias: OrginalPackageAlias,
					MethodName:          f.Name.Name,
					ReceiverType:        receiverTypeName,
					ReceiverIdFieldName: Id,
					Imports:             getImportsAsString(t.tokenSet, detectedFunction.Imports),
				}

				if customIdFieldName, hasCustomId := t.Output.TypesWithCustomId[receiverTypeName]; hasCustomId {
					newHandler.ReceiverIdFieldName = customIdFieldName
				}

				if retParamsVerifier.Check(f) {
					if len(f.Type.Results.List) > 1 {
						newHandler.OptionalReturnType = types.ExprString(f.Type.Results.List[0].Type)
						if isNobject, isPresent := t.Output.IsNobjectInOrginalPackage[newHandler.OptionalReturnType]; isNobject && isPresent {
							newHandler.OptionalReturnType = newHandler.OrginalPackageAlias + "." + newHandler.OptionalReturnType
						} else if strings.Contains(newHandler.OptionalReturnType, ReferenceListType) {
							newHandler.OptionalReturnType = strings.TrimPrefix(newHandler.OptionalReturnType, ReferenceListType)
							newHandler.OptionalReturnType = strings.Trim(newHandler.OptionalReturnType, "[]")
							newHandler.OptionalReturnType = ReferenceListType + "[" + OrginalPackageAlias + "." + newHandler.OptionalReturnType + "]"
						} else if strings.Contains(newHandler.OptionalReturnType, ReferenceType) {
							newHandler.OptionalReturnType = strings.TrimPrefix(newHandler.OptionalReturnType, ReferenceType)
							newHandler.OptionalReturnType = strings.Trim(newHandler.OptionalReturnType, "[]")
							newHandler.OptionalReturnType = ReferenceType + "[" + OrginalPackageAlias + "." + newHandler.OptionalReturnType + "]"
						}
					}
				}

				param, err := getHandlerInputParam(f.Type.Params, t.Output.TypeFields)
				if err != nil {
					fmt.Println("Maximum allowed number of parameters is 1. Handler generation for " + f.Name.Name + "skipped")
					continue
				}
				newHandler.OptionalInputType = param
				t.Handlers = append(t.Handlers, newHandler)
			}
		}
	}
}

func getImportsAsString(fset *token.FileSet, imports []*ast.ImportSpec) string {
	var buf bytes.Buffer
	for _, imp := range imports {
		err := printer.Fprint(&buf, fset, imp)
		buf.WriteString("\n")
		if err != nil {
			fmt.Println(err)
		}
	}

	return buf.String()
}

func getHandlerInputParam(params *ast.FieldList, typeFieldsInPkg map[string]map[string]string) (string, error) {
	if params.List == nil || len(params.List) == 0 {
		return "", nil
	} else if len(params.List) > 1 {
		return "", fmt.Errorf("maximum allowed number of parameters is 1")
	}

	inputParamType := types.ExprString(params.List[0].Type)
	if _, isPresent := typeFieldsInPkg[inputParamType]; isPresent {
		inputParamType = OrginalPackageAlias + "." + inputParamType
	} else if strings.Contains(inputParamType, ReferenceListType) {
		inputParamType = strings.TrimPrefix(inputParamType, ReferenceListType)
		inputParamType = strings.Trim(inputParamType, "[]")
		inputParamType = ReferenceListType + "[" + OrginalPackageAlias + "." + inputParamType + "]"
	} else if strings.Contains(inputParamType, ReferenceType) {
		inputParamType = strings.TrimPrefix(inputParamType, ReferenceType)
		inputParamType = strings.Trim(inputParamType, "[]")
		inputParamType = ReferenceType + "[" + OrginalPackageAlias + "." + inputParamType + "]"
	}

	return inputParamType, nil
}

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
	Invocation          string
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

			if strings.HasPrefix(f.Name.Name, SetPrefix) || strings.HasPrefix(f.Name.Name, GetPrefix) {
				// TODO FIX SET/GET prefix is not enough - the following string sequence
				// must refer to existing type's field name
				continue
			} else if f.Recv == nil {

				param, err := getHandlerFunctionParam(f.Type.Params, t.Output.IsNobjectInOrginalPackage)
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
				}

			} else {
				receiverTypeName := getFunctionReceiverTypeAsString(f.Recv)
				if isNobject := t.Output.IsNobjectInOrginalPackage[receiverTypeName]; !isNobject {
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
						if _, isPresent := t.Output.IsNobjectInOrginalPackage[newHandler.OptionalReturnType]; isPresent {
							newHandler.OptionalReturnType = newHandler.OrginalPackageAlias + "." + newHandler.OptionalReturnType
						}
					}
				}

				parameters, err := getHandlerFunctionParam(f.Type.Params, t.Output.IsNobjectInOrginalPackage)
				if err != nil {
					fmt.Println("Maximum allowed number of parameters is 1. Handler generation for " + f.Name.Name + "skipped")
					continue
				}
				if parameters == "" {
					newHandler.Invocation = f.Name.Name + "()"
				} else {

					newHandler.Invocation = f.Name.Name + "(" + HandlerInputParameterName + "." + HandlerInputParameterFieldName + ".(" + parameters + "))"
				}
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

func getHandlerFunctionParam(params *ast.FieldList, isNobjectInOrgPkg map[string]bool) (string, error) {
	if params.List == nil || len(params.List) == 0 {
		return "", nil
	} else if len(params.List) > 1 {
		return "", fmt.Errorf("maximum allowed number of parameters is 1")
	}

	inputParamType := types.ExprString(params.List[0].Type)
	if _, isPresent := isNobjectInOrgPkg[inputParamType]; isPresent {
		inputParamType = OrginalPackageAlias + "." + inputParamType
	}
	return inputParamType, nil
}

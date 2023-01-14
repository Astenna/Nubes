package parser

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"go/types"
)

type StateChangingHandler struct {
	OrginalPackage      string
	OrginalPackageAlias string
	Imports             string
	MethodName          string
	ReceiverType        string
	OptionalReturnType  string
	Invocation          string
}

type detectedFunction struct {
	Function *ast.FuncDecl
	Imports  []*ast.ImportSpec
}

func ParseStateChangingHandlers(path string, parsedPackage ParsedPackage) []StateChangingHandler {
	set := token.NewFileSet()
	packs, err := parser.ParseDir(set, path, nil, 0)
	assertDirParsed(err)

	fileFunctionsMap := getPackageFuncs(packs)
	isNobjectInOrgPkg := parsedPackage.IsNobjectInOrginalPackage

	handlerFuncs := []StateChangingHandler{}
	for _, functions := range fileFunctionsMap {
		for _, detectedFunction := range functions {
			f := detectedFunction.Function

			if f.Name.Name == NobjectImplementationMethod || f.Name.Name == CustomIdImplementationMethod {
				continue
			} else if f.Recv == nil {

				// TODO: detect & generate handler(s) with CTORs

			} else {
				ownerType := getFunctionReceiverTypeAsString(f.Recv)
				if isNobject := isNobjectInOrgPkg[ownerType]; !isNobject {
					fmt.Println("Member type does not implement Nobject interface. Handler generation for " + f.Name.Name + "skipped")
					continue
				}

				newHandler := StateChangingHandler{
					OrginalPackage:      parsedPackage.ImportPath,
					OrginalPackageAlias: OrginalPackageAlias,
					MethodName:          f.Name.Name,
					ReceiverType:        ownerType,
					Imports:             getImportsAsString(set, detectedFunction.Imports),
				}

				if retParamsVerifier.Check(f) {
					if len(f.Type.Results.List) > 1 {
						newHandler.OptionalReturnType = types.ExprString(f.Type.Results.List[0].Type)
						if _, isPresent := isNobjectInOrgPkg[newHandler.OptionalReturnType]; isPresent {
							newHandler.OptionalReturnType = newHandler.OrginalPackageAlias + "." + newHandler.OptionalReturnType
						}
					}
				}

				parameters, err := getStateChangingFuncParams(f.Type.Params, isNobjectInOrgPkg)
				if err != nil {
					fmt.Println("Maximum allowed number of parameters is 1. Handler generation for " + f.Name.Name + "skipped")
					continue
				}
				newHandler.Invocation = f.Name.Name + "(" + parameters + ")"
				handlerFuncs = append(handlerFuncs, newHandler)
			}
		}
	}

	return handlerFuncs
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

func getStateChangingFuncParams(params *ast.FieldList, isNobjectInOrgPkg map[string]bool) (string, error) {
	if params.List == nil || len(params.List) == 0 {
		return "", nil
	} else if len(params.List) > 1 {
		return "", fmt.Errorf("maximum allowed number of parameters is 1")
	}

	inputParamType := types.ExprString(params.List[0].Type)
	if _, isPresent := isNobjectInOrgPkg[inputParamType]; isPresent {
		inputParamType = OrginalPackageAlias + "." + inputParamType
	}
	return HandlerInputParameterName + "." + HandlerInputEmbededOrginalFunctionParameterName + ".(" + inputParamType + ")", nil
}

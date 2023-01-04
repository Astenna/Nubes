package parser

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"go/types"
	"strconv"
	"strings"
)

type StateChangingHandler struct {
	OrginalPackage      string
	OrginalPackageAlias string
	Imports             string
	MethodName          string
	ReceiverType        string
	OptionalReturnType  string
	Invocation          string
	Stateless           bool
	ErrorReturned       bool
}

type detectedFunction struct {
	Function *ast.FuncDecl
	Imports  []*ast.ImportSpec
}

func ParseStateChangingHandlers(path string, nobjectsImportPath string, definedInOrgPackage map[string]bool) []StateChangingHandler {
	set := token.NewFileSet()
	packs, err := parser.ParseDir(set, path, nil, 0)
	AssertDirParsed(err)

	fileFunctionsMap := GetPackageFuncs(packs)

	handlerFuncs := []StateChangingHandler{}
	for _, functions := range fileFunctionsMap {
		for _, detectedFunction := range functions {
			f := detectedFunction.Function
			if f.Recv == nil || f.Name.Name == NobjectImplementationMethod || f.Name.Name == CustomIdImplementationMethod {
				continue
			}

			ownerType := strings.TrimPrefix(types.ExprString(f.Recv.List[0].Type), "*")
			if _, ok := definedInOrgPackage[ownerType]; !ok {
				fmt.Println("Member type does not implement Nobject interface. Handler generation for " + f.Name.Name + "skipped")
				continue
			}

			newHandler := StateChangingHandler{
				OrginalPackage:      nobjectsImportPath,
				OrginalPackageAlias: OrginalPackageAlias,
				MethodName:          f.Name.Name,
				ReceiverType:        ownerType,
				Imports:             GetImports(set, detectedFunction.Imports),
			}

			errorTypeFound := false
			if f.Type.Results != nil {
				errorTypeFound = types.ExprString(f.Type.Results.List[len(f.Type.Results.List)-1].Type) == "error"
				if (!errorTypeFound && len(f.Type.Results.List) > 1) || (errorTypeFound && len(f.Type.Results.List) > 2) {
					fmt.Println("Maximum allowed number of non-error return parameters is 1. Handler generation for " + f.Name.Name + "skipped")
					continue
				}

				if errorTypeFound {
					newHandler.ErrorReturned = true
					if len(f.Type.Results.List) > 1 {
						newHandler.OptionalReturnType = types.ExprString(f.Type.Results.List[0].Type)
						if _, ok := definedInOrgPackage[newHandler.OptionalReturnType]; ok {
							newHandler.OptionalReturnType = newHandler.OrginalPackageAlias + "." + newHandler.OptionalReturnType
						}
					}
				}
			}

			// stateless method, instance will be created just to invoke the method
			// stateful method, create instance to invoke the method and then save state changes
			newHandler.Stateless = f.Recv.List[0].Names == nil

			parameters, err := GetStateChangingFuncParams(f.Type.Params, definedInOrgPackage)
			if err != nil {
				fmt.Println("Maximum allowed number of parameters is 1. Handler generation for " + f.Name.Name + "skipped")
				continue
			}
			newHandler.Invocation = f.Name.Name + "(" + parameters + ")"
			handlerFuncs = append(handlerFuncs, newHandler)
		}
	}

	return handlerFuncs
}

func GetStateChangingFuncParams(params *ast.FieldList, definedInOrgPackage map[string]bool) (string, error) {
	if params.List == nil || len(params.List) == 0 {
		return "", nil
	} else if len(params.List) > 1 {
		return "", fmt.Errorf("maximum allowed number of parameters is 1")
	}

	inputParamType := types.ExprString(params.List[0].Type)
	if _, ok := definedInOrgPackage[inputParamType]; ok {
		inputParamType = OrginalPackageAlias + "." + inputParamType
	}
	return HandlerInputParameterName + "." + HandlerInputEmbededOrginalFunctionParameterName + ".(" + inputParamType + ")", nil
}

func GetReturnTypesDefinition(results *ast.FieldList, nobjectTypes map[string]bool) (string, error) {

	if len(results.List) > 1 {
		return "", fmt.Errorf("maximum allowed number of non-error return parameters is 1, found " + strconv.Itoa(len(results.List)))
	}
	if types.ExprString(results.List[len(results.List)-1].Type) != "error" {
		return "", fmt.Errorf("the last return parameter type of repository function must be error")
	}

	if len(results.List) == 1 {
		return "error", nil
	}

	// given the prev. conditions now we're sure
	// the case is: (optionalReturnType, error)
	returnParam := types.ExprString(results.List[0].Type)

	if _, ok := nobjectTypes[returnParam]; ok {
		return "( " + OrginalPackageAlias + "." + returnParam + " , error)", nil
	}

	return "( " + returnParam + " , error)", nil
}

func GetFunctionBody(fset *token.FileSet, body *ast.BlockStmt) (string, error) {
	var buf bytes.Buffer
	err := printer.Fprint(&buf, fset, body)
	if err != nil {
		return "", fmt.Errorf("error occurred when parsing the function body")
	}
	return buf.String(), nil
}

func GetImports(fset *token.FileSet, imports []*ast.ImportSpec) string {
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

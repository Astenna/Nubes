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
	OrginalPackage           string
	OrginalPackageAlias      string
	Imports                  string
	Signature                string
	OwnerVariableName        string
	OwnerType                string
	ReturnFromInvocation     string
	Invocation               string
	HandlerNameWithoutSuffix string
	Stateless                bool
	OptionalReturnVar        string
	OptionalReturnType       string
}

type DefaultRepoHandler struct {
	TypesPackageAlias string
	TypesPackagePath  string
	TypeName          string
	OperationName     string
}

type CustomRepoHandler struct {
	Imports       string
	OperationName string
	TypeName      string
	Parameters    string
	ReturnValues  string
	Body          string
}

type detectedFunction struct {
	Function *ast.FuncDecl
	Imports  []*ast.ImportSpec
}

func ParseRepoHandlers(path string, nobjectsImportPath string, nobjectTypes map[string]struct{}) ([]CustomRepoHandler, []DefaultRepoHandler) {
	set := token.NewFileSet()
	packs, err := parser.ParseDir(set, path, nil, 0)
	AssertDirParsed(err)

	functions := GetPackageFuncs(packs)

	isNoObjectMethodDefined := make(map[string]map[string]bool, len(nobjectTypes))
	for i := range nobjectTypes {
		isNoObjectMethodDefined[i] = map[string]bool{GetPrefix: false, CreatePrefix: false, DeletePrefix: false, UpdatePrefix: false}
	}

	repoCustomFuncs := []CustomRepoHandler{}
	for _, detectedFunction := range functions {
		f := detectedFunction.Function
		var methodType string

		switch {
		case strings.HasPrefix(f.Name.Name, GetPrefix):
			methodType = GetPrefix
		case strings.HasPrefix(f.Name.Name, CreatePrefix):
			methodType = CreatePrefix
		case strings.HasPrefix(f.Name.Name, DeletePrefix):
			methodType = DeletePrefix
		case strings.HasPrefix(f.Name.Name, UpdatePrefix):
			methodType = UpdatePrefix
		default:
			continue
		}

		for typeName := range nobjectTypes {
			if strings.HasSuffix(f.Name.Name, typeName) {
				params, err := GetCustomRepoFuncParams(f.Type.Params)
				if err != nil {
					fmt.Println("faas handler " + f.Name.Name + " custom definition replaced with default definition")
					continue
				}

				returnParams, err := GetReturnTypesDefinition(f.Type.Results, nobjectTypes)
				if err != nil {
					fmt.Println("faas handler " + f.Name.Name + " custom definition replaced with default definition")
					continue
				}

				body, err := GetFunctionBody(set, f.Body)
				if err != nil {
					fmt.Printf("faas handler " + f.Name.Name + " custom definition replaced with default definition")
					continue
				}

				isNoObjectMethodDefined[typeName][methodType] = true
				repoCustomFuncs = append(repoCustomFuncs, CustomRepoHandler{
					OperationName: methodType,
					TypeName:      typeName,
					Parameters:    params,
					ReturnValues:  returnParams,
					Body:          body,
					Imports:       GetImports(set, detectedFunction.Imports),
				})
			}
		}
	}

	repoDefaultFuncs := GetDefaultRepoHandler(isNoObjectMethodDefined, nobjectsImportPath)
	return repoCustomFuncs, repoDefaultFuncs
}

func GetDefaultRepoHandler(isNoObjectMethodDefined map[string]map[string]bool, nobjectsImportPath string) []DefaultRepoHandler {
	var defaultFuncs []DefaultRepoHandler

	for typeName, typeMethodsMap := range isNoObjectMethodDefined {
		for method, isCustom := range typeMethodsMap {
			if !isCustom {
				defaultFuncs = append(defaultFuncs, DefaultRepoHandler{
					TypesPackageAlias: OrginalPackageAlias,
					TypesPackagePath:  nobjectsImportPath,
					TypeName:          typeName,
					OperationName:     method,
				})
			}
		}
	}

	return defaultFuncs
}

func ParseStateChangingHandlers(path string, nobjectsImportPath string, nobjectTypes map[string]struct{}) []StateChangingHandler {
	set := token.NewFileSet()
	packs, err := parser.ParseDir(set, path, nil, 0)
	AssertDirParsed(err)

	functions := GetPackageFuncs(packs)

	handlerFuncs := []StateChangingHandler{}
	for _, detectedFunction := range functions {
		f := detectedFunction.Function
		if f.Recv == nil || f.Name.Name == GetTypeName {
			continue
		}

		ownerType := strings.TrimPrefix(types.ExprString(f.Recv.List[0].Type), "*")
		if _, ok := nobjectTypes[ownerType]; !ok {
			fmt.Println("Member type does not implement Nobject interface. Handler generation for " + f.Name.Name + "skipped")
			continue
		}

		newHandler := StateChangingHandler{
			OrginalPackage:           nobjectsImportPath,
			OrginalPackageAlias:      OrginalPackageAlias,
			HandlerNameWithoutSuffix: f.Name.Name,
			Signature:                "func " + f.Name.Name + HandlerSuffix + HandlerParameters,
			OwnerType:                ownerType,
			Imports:                  GetImports(set, detectedFunction.Imports),
		}

		// 4 cases:
		// C1: no return parameters
		// C2: 1 return: error
		// C3: 1 return: non-error
		// C4: 2 return: non-error, error
		errorTypeFound := false
		if f.Type.Results == nil {
			// C1
			newHandler.Signature += " error"
		} else {
			errorTypeFound = types.ExprString(f.Type.Results.List[len(f.Type.Results.List)-1].Type) == "error"
			if !errorTypeFound && len(f.Type.Results.List) > 1 {
				fmt.Println("Maximum allowed number of non-error return parameters is 1. Handler generation for " + f.Name.Name + "skipped")
				continue
			} else {

				newHandler.OptionalReturnType = types.ExprString(f.Type.Results.List[0].Type)
				if _, ok := nobjectTypes[newHandler.OptionalReturnType]; ok {
					newHandler.OptionalReturnType = newHandler.OrginalPackageAlias + "." + newHandler.OptionalReturnType
				}

				if !errorTypeFound {
					// C3
					newHandler.Signature += "(" + newHandler.OptionalReturnType + ", error)"
					newHandler.ReturnFromInvocation = "result :="
					newHandler.OptionalReturnVar = "result"
				} else {
					if len(f.Type.Results.List) == 1 {
						// C2
						newHandler.Signature += " error"
						newHandler.ReturnFromInvocation = "_err :="
					} else {
						// C4
						newHandler.Signature += "(" + newHandler.OptionalReturnType + " ,error)"
						newHandler.ReturnFromInvocation = "result, _err :="
						newHandler.OptionalReturnVar = "result"
					}
				}
			}
		}

		var ownerTypeName string
		if f.Recv.List[0].Names == nil {
			// stateless method, instance will be created just to invoke the method
			newHandler.Stateless = true
			newHandler.OwnerVariableName = "typeInstance"
			ownerTypeName = "typeInstance"
		} else {
			newHandler.Stateless = false
			// stateful method, create instance to invoke the method and then save state changes
			ownerTypeName = f.Recv.List[0].Names[0].Name
			newHandler.OwnerVariableName = f.Recv.List[0].Names[0].Name
		}

		parameters, err := GetStateChangingFuncParams(f.Type.Params)
		if err != nil {
			fmt.Println("Maximum allowed number of parameters is 1. Handler generation for " + f.Name.Name + "skipped")
			continue
		}
		newHandler.Invocation = ownerTypeName + "." + f.Name.Name + "(" + parameters + ")"

		handlerFuncs = append(handlerFuncs, newHandler)
	}

	return handlerFuncs
}

func GetPackageFuncs(packs map[string]*ast.Package) []detectedFunction {
	var detectedFunctions []detectedFunction

	for _, pack := range packs {
		for _, f := range pack.Files {
			for _, d := range f.Decls {
				if fn, isFn := d.(*ast.FuncDecl); isFn {
					detectedFunctions = append(detectedFunctions, detectedFunction{
						Function: fn,
						Imports:  f.Imports,
					})
				}
			}
		}
	}

	return detectedFunctions
}

func GetStateChangingFuncParams(params *ast.FieldList) (string, error) {
	if params.List == nil || len(params.List) == 0 {
		return "", nil
	} else if len(params.List) > 1 {
		return "", fmt.Errorf("maximum allowed number of parameters is 1")
	}

	return HandlerInputParameterName + "." + HandlerInputEmbededOrginalFunctionParameterName + ".(" + types.ExprString(params.List[0].Type) + ")", nil
}

func GetCustomRepoFuncParams(params *ast.FieldList) (string, error) {
	if params.List == nil || len(params.List) == 0 {
		return "", nil
	} else if len(params.List) > 1 {
		return "", fmt.Errorf("maximum allowed number of parameters is 1")
	}

	return params.List[0].Names[0].Name + " " + types.ExprString(params.List[0].Type), nil
}

func GetReturnTypesDefinition(results *ast.FieldList, nobjectTypes map[string]struct{}) (string, error) {

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

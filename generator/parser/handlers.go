package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"strings"
)

type HandlerFunc struct {
	OrginalPackage       string
	OrginalPackageAlias  string
	Imports              string
	Signature            string
	OwnerVariableName    string
	OwnerType            string
	ReturnFromInvocation string
	Invocation           string
	HandlerName          string
	Stateless            bool
	OptionalReturnVar    string
}

func PrepareHandlersFromFunctions(path string, moduleName string) []HandlerFunc {
	return nil
}

func PrepareHandlersFromMethods(path string, moduleName string) []HandlerFunc {
	set := token.NewFileSet()
	packs, err := parser.ParseDir(set, path, nil, 0)
	AssertDirParsed(err)

	funcsMap := make(map[string][]*ast.FuncDecl)

	for packageName, pack := range packs {
		for _, f := range pack.Files {
			for _, d := range f.Decls {
				if fn, isFn := d.(*ast.FuncDecl); isFn {

					if funcsMap[packageName] == nil {
						funcsMap[packageName] = []*ast.FuncDecl{}
					}
					funcsMap[packageName] = append(funcsMap[packageName], fn)
				}
			}
		}
	}

	handlerFuncs := []HandlerFunc{}
	for packageName, funcs := range funcsMap {
		for _, f := range funcs {
			if f.Recv == nil || f.Name.Name == GetTypeName {
				continue
			}

			parametersList := "(" + HandlerInputParameterName + " " + HandlerInputParameterType + ")"
			//strings.SplitAfter("(id string, "+strings.TrimPrefix(types.ExprString(f.Type), "func("), ")")[0]
			newHandler := HandlerFunc{
				OrginalPackage:      moduleName + "/" + packageName,
				OrginalPackageAlias: OrginalPackageAlias,
				HandlerName:         f.Name.Name + HandlerSuffix,
				Signature:           "func " + f.Name.Name + HandlerSuffix + parametersList,
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
				} else if !errorTypeFound {
					// C3
					newHandler.Signature += "(" + newHandler.OrginalPackageAlias + "." + types.ExprString(f.Type.Results.List[0].Type) + ", error)"
					newHandler.ReturnFromInvocation = "result :="
					newHandler.OptionalReturnVar = "result"
				} else {
					if len(f.Type.Results.List) == 1 {
						// C2
						newHandler.Signature += " error"
						newHandler.ReturnFromInvocation = "err :="
					} else {
						// C4
						newHandler.Signature += "(" + newHandler.OrginalPackageAlias + "." + types.ExprString(f.Type.Results.List[0].Type) + " ,error)"
						newHandler.ReturnFromInvocation = "result, err :="
						newHandler.OptionalReturnVar = "result"
					}
				}
			}

			newHandler.OwnerType = strings.TrimPrefix(types.ExprString(f.Recv.List[0].Type), "*")
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

			parameters, err := GetOrginalFunctionParameters(f.Type.Params)
			if err != nil {
				fmt.Println("Maximum allowed number of parameters is 1. Handler generation for " + f.Name.Name + "skipped")
				continue
			}
			newHandler.Invocation = ownerTypeName + "." + f.Name.Name + "(" + parameters + ")"

			handlerFuncs = append(handlerFuncs, newHandler)
		}
	}

	return handlerFuncs
}

func GetOrginalFunctionParameters(params *ast.FieldList) (string, error) {
	if params.List == nil || len(params.List) == 0 {
		return "", nil
	} else if len(params.List) > 1 {
		return "", fmt.Errorf("maximum allowed number of parameters is 1")
	}

	return HandlerInputParameterName + "." + HandlerInputEmbededOrginalFunctionParameterName + ".(" + types.ExprString(params.List[0].Type) + ")", nil
}

func AssertDirParsed(err error) {
	if err != nil {
		fmt.Println("Failed to parse files in the directory: %w", err)
		os.Exit(1)
	}
}

package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"strings"
)

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

func ParseRepoHandlers(path string, nobjectsImportPath string, nobjectTypes map[string]bool) ([]CustomRepoHandler, []DefaultRepoHandler) {
	set := token.NewFileSet()
	packs, err := parser.ParseDir(set, path, nil, 0)
	AssertDirParsed(err)

	fileFunctionsMap := GetPackageFuncs(packs)

	isNoObjectMethodDefined := make(map[string]map[string]bool, len(nobjectTypes))
	for i, isNobject := range nobjectTypes {
		isNoObjectMethodDefined[i] = map[string]bool{GetPrefix: !isNobject, CreatePrefix: !isNobject, DeletePrefix: !isNobject, UpdatePrefix: !isNobject}
	}

	repoCustomFuncs := []CustomRepoHandler{}
	for _, functions := range fileFunctionsMap {
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

func GetCustomRepoFuncParams(params *ast.FieldList) (string, error) {
	if params.List == nil || len(params.List) == 0 {
		return "", nil
	} else if len(params.List) > 1 {
		return "", fmt.Errorf("maximum allowed number of parameters is 1")
	}

	return params.List[0].Names[0].Name + " " + types.ExprString(params.List[0].Type), nil
}

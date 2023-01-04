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

func GetPackageTypes(path string, moduleName string) (map[string]bool, string) {
	set := token.NewFileSet()
	packs, err := parser.ParseDir(set, path, nil, 0)
	AssertDirParsed(err)

	isNobjectType := make(map[string]bool)

	var packageName string
	for pckgName, pack := range packs {
		for _, f := range pack.Files {
			for _, d := range f.Decls {
				if fn, isFn := d.(*ast.FuncDecl); isFn {
					if fn.Recv != nil && fn.Name.Name == NobjectImplementationMethod {
						ownerType := types.ExprString(fn.Recv.List[0].Type)
						isNobjectType[ownerType] = true
					}
				}

				if genDecl, ok := d.(*ast.GenDecl); ok {
					for _, elem := range genDecl.Specs {
						if typeSpec, ok := elem.(*ast.TypeSpec); ok {
							typeName := strings.TrimPrefix(typeSpec.Name.Name, "*")
							if _, isPresent := isNobjectType[typeName]; !isPresent {
								isNobjectType[typeName] = false
							}
						}
					}
				}
			}
		}
		packageName = pckgName
	}

	return isNobjectType, moduleName + "/" + packageName
}

func GetPackageFuncs(packs map[string]*ast.Package) map[string][]detectedFunction {
	detectedFunctions := make(map[string][]detectedFunction)

	for _, pack := range packs {
		for path, f := range pack.Files {
			for _, d := range f.Decls {
				if fn, isFn := d.(*ast.FuncDecl); isFn {
					detectedFunctions[path] = append(detectedFunctions[f.Name.Name], detectedFunction{
						Function: fn,
						Imports:  f.Imports,
					})
				}
			}
		}
	}

	return detectedFunctions
}

func AssertDirParsed(err error) {
	if err != nil {
		fmt.Println("Failed to parse files in the directory: %w", err)
		os.Exit(1)
	}
}

package parser

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"strings"
)

type ParsedPackage struct {
	ImportPath                string
	IsNobjectInOrginalPackage map[string]bool
	TypeFields                map[string]map[string]string
	TypeAttributesIndexes     map[string][]string
	TypesWithCustomId         map[string]string
}

func GetPackageTypes(path string, moduleName string) ParsedPackage {
	set := token.NewFileSet()
	packs, err := parser.ParseDir(set, path, nil, 0)
	assertDirParsed(err)

	result := ParsedPackage{
		IsNobjectInOrginalPackage: make(map[string]bool),
		TypesWithCustomId:         map[string]string{},
		TypeAttributesIndexes:     map[string][]string{},
		TypeFields:                map[string]map[string]string{},
	}

	for packageName, pack := range packs {
		for _, f := range pack.Files {
			for _, d := range f.Decls {
				if fn, isFn := d.(*ast.FuncDecl); isFn {
					if fn.Recv != nil && fn.Name.Name == NobjectImplementationMethod {
						ownerType := getFunctionReceiverTypeAsString(fn.Recv)
						result.IsNobjectInOrginalPackage[ownerType] = true
					}
					if fn.Recv != nil && fn.Name.Name == CustomIdImplementationMethod {
						ownerType := getFunctionReceiverTypeAsString(fn.Recv)
						idFieldName, err := getIdFieldNameFromCustomIdImpl(fn)
						if err != nil {
							fmt.Println(err)
							continue
						}
						result.TypesWithCustomId[ownerType] = idFieldName
					}
				}

				if genDecl, ok := d.(*ast.GenDecl); ok {
					for _, elem := range genDecl.Specs {
						if typeSpec, ok := elem.(*ast.TypeSpec); ok {
							typeName := strings.TrimPrefix(typeSpec.Name.Name, "*")
							if _, isPresent := result.IsNobjectInOrginalPackage[typeName]; !isPresent {
								result.IsNobjectInOrginalPackage[typeName] = false
							}

							if strctType, ok := typeSpec.Type.(*ast.StructType); ok {
								parseStructFields(strctType, typeName, &result)
							}
						}
					}
				}
			}
		}
		result.ImportPath = moduleName + "/" + packageName
	}

	// TODO: loop over all detected indexes, to verify if the attributes are named according to the convention
	return result
}

func parseStructFields(strctType *ast.StructType, typeName string, parsedPackage *ParsedPackage) {
	if strctType == nil || strctType.Fields == nil || len(strctType.Fields.List) == 0 {
		return
	}

	parsedPackage.TypeFields[typeName] = make(map[string]string, len(strctType.Fields.List))
	for _, field := range strctType.Fields.List {
		fieldType := types.ExprString(field.Type)
		parsedPackage.TypeFields[typeName][field.Names[0].Name] = fieldType

		if strings.Contains(fieldType, LibraryReferenceNavigationList) {
			relationshipOrginalOwner := strings.Split(fieldType, ",")[1]
			relationshipOrginalOwner = strings.TrimRight(relationshipOrginalOwner, "]")

			parsedPackage.TypeAttributesIndexes[relationshipOrginalOwner] = append(parsedPackage.TypeAttributesIndexes[relationshipOrginalOwner], typeName)
		}
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

func getPackageFuncs(packs map[string]*ast.Package) map[string][]detectedFunction {
	detectedFunctions := make(map[string][]detectedFunction)

	for _, pack := range packs {
		for path, f := range pack.Files {
			for _, d := range f.Decls {
				if fn, isFn := d.(*ast.FuncDecl); isFn {
					detectedFunctions[path] = append(detectedFunctions[path], detectedFunction{
						Function: fn,
						Imports:  f.Imports,
					})
				}
			}
		}
	}

	return detectedFunctions
}

func assertDirParsed(err error) {
	if err != nil {
		fmt.Println("Failed to parse files in the directory: %w", err)
		os.Exit(1)
	}
}

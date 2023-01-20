package parser

import (
	"errors"
	"fmt"
	"go/ast"
	"go/types"
	"os"
	"strings"
)

func (t *TypeSpecParser) detectAndAdjustDecls() {
	for _, pack := range t.packs {
		for path, f := range pack.Files {
			t.fileChanged[path] = false
			for _, d := range f.Decls {
				if genDecl, ok := d.(*ast.GenDecl); ok {
					for _, elem := range genDecl.Specs {
						if typeSpec, ok := elem.(*ast.TypeSpec); ok {
							typeName := strings.TrimPrefix(typeSpec.Name.Name, "*")
							if _, isPresent := t.Output.IsNobjectInOrginalPackage[typeName]; !isPresent {
								t.Output.IsNobjectInOrginalPackage[typeName] = false
							}

							if strctType, ok := typeSpec.Type.(*ast.StructType); ok {
								modified := parseStructFields(strctType, typeName, &t.Output)

								if !t.fileChanged[path] {
									t.fileChanged[path] = modified
								}
							}
						}
					}
				}
			}
		}
	}
}

// The parseStructFields returns true if the ast representing
// the struct was modified, otherwise false
func parseStructFields(strctType *ast.StructType, typeName string, parsedPackage *ParsedPackage) bool {
	structDefinitionModified := false
	if strctType == nil || strctType.Fields == nil || len(strctType.Fields.List) == 0 {
		return structDefinitionModified
	}

	parsedPackage.TypeFields[typeName] = make(map[string]string, len(strctType.Fields.List))
	fieldModified := false
	isNobject := parsedPackage.IsNobjectInOrginalPackage[typeName]
	for _, field := range strctType.Fields.List {
		fieldType := types.ExprString(field.Type)
		parsedPackage.TypeFields[typeName][field.Names[0].Name] = fieldType

		if isNobject {
			fieldModified = parseRelationshipsTags(field, typeName, fieldType, parsedPackage)
			if !structDefinitionModified {
				structDefinitionModified = fieldModified
			}
		}
	}

	if _, exists := parsedPackage.TypeFields[typeName][IsInitializedFieldName]; !exists && isNobject {
		strctType.Fields.List = append(strctType.Fields.List, &ast.Field{
			Names: []*ast.Ident{{Name: IsInitializedFieldName}}, Type: &ast.Ident{Name: "bool"},
		})
		return true
	}

	return structDefinitionModified
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

func assertDirParsed(err error) {
	if err != nil {
		fmt.Println("Failed to parse files in the directory: %w", err)
		os.Exit(1)
	}
}

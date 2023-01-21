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

func (t *TypeSpecParser) detectAndAdjustDecls() {
	for _, pack := range t.packs {
		for path, f := range pack.Files {
			t.fileChanged[path] = false
			for _, d := range f.Decls {
				if genDecl, ok := d.(*ast.GenDecl); ok {
					for _, elem := range genDecl.Specs {
						if typeSpec, ok := elem.(*ast.TypeSpec); ok {
							typeName := strings.TrimPrefix(typeSpec.Name.Name, "*")
							isNobjectType, isPresent := t.Output.IsNobjectInOrginalPackage[typeName]

							if !isPresent {
								t.Output.IsNobjectInOrginalPackage[typeName] = false
							}

							if strctType, ok := typeSpec.Type.(*ast.StructType); ok {
								modified := parseStructFieldsForTypeSpec(strctType, typeName, &t.Output)

								if !t.fileChanged[path] {
									t.fileChanged[path] = modified
								}

								if isNobjectType && !t.isInitAlreadyAdded[typeName] {
									t.addInitFunctionDefinition(f, typeName)
									t.fileChanged[path] = true
								}
							}
						}
					}
				}
			}
		}
	}
}

func (t *TypeSpecParser) addInitFunctionDefinition(f *ast.File, typeName string) {
	idFieldName := getIdFieldNameOfType(typeName, t.Output.TypesWithCustomId)
	function := getInitFunctionForType(typeName, idFieldName, t.Output.TypeNavListsReferringFieldName[typeName], t.Output.ManyToManyRelationships[typeName])
	f.Decls = append(f.Decls, function)
}

// The parseStructFieldsForTypeSpec returns true if the ast representing
// the struct was modified, otherwise false
func parseStructFieldsForTypeSpec(strctType *ast.StructType, typeName string, parsedPackage *ParsedPackage) bool {
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

func getTypeSpecAsString(fset *token.FileSet, detectedStruct *ast.TypeSpec) (string, error) {
	var buf bytes.Buffer
	buf.WriteString("type ")
	err := printer.Fprint(&buf, fset, detectedStruct)
	if err != nil {
		return "", fmt.Errorf("error occurred when parsing the struct")
	}
	return buf.String(), nil
}

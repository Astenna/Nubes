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

func (t *ClientTypesParser) detectGenDecls() {
	for _, pack := range t.packs {
		for _, f := range pack.Files {
			for _, decl := range f.Decls {
				if genDecl, ok := decl.(*ast.GenDecl); ok {
					for _, elem := range genDecl.Specs {
						if typeSpec, ok := elem.(*ast.TypeSpec); ok {
							if strctType, ok := typeSpec.Type.(*ast.StructType); ok {
								typeName := strings.TrimPrefix(typeSpec.Name.Name, "*")

								makeFieldsUnexported(strctType.Fields)
								if _, ok := t.DefinedTypes[typeName]; !ok {
									t.DefinedTypes[typeName] = &StructTypeDefinition{}
								}

								t.DefinedTypes[typeName].TypeNameOrginalCase = typeName
								t.DefinedTypes[typeName].TypeNameLower = lowerCasedFirstChar(typeName)

								t.parseStructFieldsForClients(strctType, typeName)
							} else {
								// DETECT AND SAVE CUSTOM TYPES (e.g. type MyInt int)
								def, err := getTypeSpecAsString(t.tokenSet, typeSpec)
								if err != nil {
									fmt.Println(err)
								} else {
									t.OtherDecls.GenDecls = append(t.OtherDecls.GenDecls, def)
								}
							}
						}
					}
					// DETECT AND SAVE CONST DECLARATIONS
					if genDecl.Tok == token.CONST {
						constStr, err := getConstAsString(t.tokenSet, genDecl)
						if err != nil {
							fmt.Println(err)
						}
						t.OtherDecls.Consts = append(t.OtherDecls.Consts, constStr)
						continue
					}
				}
			}
		}
	}
}

func (t *ClientTypesParser) parseStructFieldsForClients(astStrct *ast.StructType, typeName string) {
	if astStrct == nil || astStrct.Fields == nil || astStrct.Fields.List == nil {
		return
	}

	for _, field := range astStrct.Fields.List {
		fieldType := strings.TrimPrefix(types.ExprString(field.Type), "*")

		if strings.Contains(fieldType, LibraryReferenceNavigationList) {
			err := t.parseRelationshipsTagsClient(field, typeName, fieldType)
			if err != nil {
				fmt.Println("error occurred when parsing relationship tags", err)
			}

		} else {
			parseFields(field, fieldType, t.DefinedTypes[typeName])
		}
	}
}

func (t *ClientTypesParser) parseRelationshipsTagsClient(field *ast.Field, typeName string, fieldType string) error {
	tags, err := getParsedTags(field)
	if err != nil {
		return err
	}
	if tags == nil {
		return fmt.Errorf("invalid usage of %s. Missing tag definition", LibraryReferenceNavigationList)
	}

	if tag, _ := tags.Get(NubesTagKey); tag != nil {
		if strings.Contains(tag.Name, HasOneTag) {
			splitted := strings.Split(tag.Name, "-")
			navigationToFieldName := ""
			if len(splitted) > 0 {
				navigationToFieldName = splitted[1]
			} else {
				return fmt.Errorf("%s detected, but missing reffering field name. Referring field name should be specified after - charcter, e.g.: %s <referring_field_name>", HasOneTag, HasOneTag)
			}

			navigationToTypeName := strings.TrimPrefix(fieldType, LibraryReferenceNavigationList)
			navigationToTypeName = strings.Trim(navigationToTypeName, "[]")

			oneToMany := OneToManyRelationshipField{TypeName: navigationToTypeName, FieldName: navigationToFieldName, FromFieldName: field.Names[0].Name}
			t.DefinedTypes[typeName].OneToManyRelationships = append(t.DefinedTypes[typeName].OneToManyRelationships, oneToMany)
		} else if strings.Contains(tag.Name, HasManyTag) {
			navigationToTypeName := strings.TrimPrefix(fieldType, LibraryReferenceNavigationList)
			navigationToTypeName = strings.Trim(navigationToTypeName, "[]")

			newManyToManyRelationship := NewManyToManyRelationshipField(typeName, navigationToTypeName, field.Names[0].Name)
			newManyToManyRelationship.FromFieldName = field.Names[0].Name
			t.DefinedTypes[typeName].ManyToManyRelationships = append(t.DefinedTypes[typeName].ManyToManyRelationships, *newManyToManyRelationship)
		}
	} else {
		return fmt.Errorf("invalid usage of %s. Missing tag definition", LibraryReferenceNavigationList)
	}

	return nil
}

func parseFields(field *ast.Field, fieldType string, structDef *StructTypeDefinition) {
	newFieldDefinition := FieldDefinition{
		FieldNameUpper: upperCaseFirstChar(field.Names[0].Name),
		FieldName:      field.Names[0].Name,
		FieldType:      fieldType,
	}

	if strings.Contains(newFieldDefinition.FieldType, ReferenceListType) {
		newFieldDefinition.FieldType = strings.TrimPrefix(newFieldDefinition.FieldType, ReferenceListType)
		newFieldDefinition.FieldTypeUpper = strings.Trim(newFieldDefinition.FieldType, "[]")
		newFieldDefinition.FieldType = lowerCasedFirstChar(newFieldDefinition.FieldTypeUpper)
		newFieldDefinition.IsReferenceList = true
	} else if strings.Contains(newFieldDefinition.FieldType, ReferenceType) {
		newFieldDefinition.FieldType = strings.TrimPrefix(newFieldDefinition.FieldType, ReferenceType)
		newFieldDefinition.FieldTypeUpper = strings.Trim(newFieldDefinition.FieldType, "[]")
		newFieldDefinition.FieldType = lowerCasedFirstChar(newFieldDefinition.FieldTypeUpper)
		newFieldDefinition.IsReference = true
	}

	if field.Names[0].Name == "id" || (field.Tag != nil && isReadonly(field)) {
		newFieldDefinition.IsReadonly = true
	}
	if field.Tag != nil {
		newFieldDefinition.Tags = field.Tag.Value
	}

	structDef.FieldDefinitions = append(structDef.FieldDefinitions, newFieldDefinition)
}

func getConstAsString(fset *token.FileSet, detectedGenDecl *ast.GenDecl) (string, error) {
	var buf bytes.Buffer
	err := printer.Fprint(&buf, fset, detectedGenDecl)
	if err != nil {
		return "", fmt.Errorf("error occurred when parsing the struct")
	}
	return buf.String(), nil
}

func makeFieldsUnexported(fieldList *ast.FieldList) {
	for _, field := range fieldList.List {
		field.Names[0].Name = lowerCasedFirstChar(field.Names[0].Name)
	}
}

func lowerCasedFirstChar(str string) string {
	if len(str) < 2 {
		return strings.ToLower(str)
	}

	bts := []byte(str)

	firstByte := bytes.ToLower([]byte{bts[0]})
	rest := bts[1:]

	str = string(bytes.Join([][]byte{firstByte, rest}, nil))
	return str
}

func upperCaseFirstChar(str string) string {
	if len(str) < 2 {
		return strings.ToUpper(str)
	}

	bts := []byte(str)

	firstByte := bytes.ToUpper([]byte{bts[0]})
	rest := bts[1:]

	str = string(bytes.Join([][]byte{firstByte, rest}, nil))
	return str
}

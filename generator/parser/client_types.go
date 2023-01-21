package parser

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/exp/maps"
)

type StructTypeDefinition struct {
	PackageName             string
	Imports                 string
	NobjectImplementation   string
	CustomIdImplementation  string
	CustomIdReceiverName    string
	TypeNameLower           string
	TypeNameOrginalCase     string
	MemberFunctions         []MethodDefinition
	FieldDefinitions        []FieldDefinition
	OneToManyRelationships  []OneToManyRelationshipField
	ManyToManyRelationships []ManyToManyRelationshipField
	CustomCtorDefinition    CustomCtorDefinition
}

type MethodDefinition struct {
	ReceiverName            string
	FuncName                string
	InputParamType          string
	OptionalReturnType      string
	OptionalReturnTypeUpper string
	IsReturnTypeNobject     bool
}

type FieldDefinition struct {
	FieldNameUpper  string
	FieldName       string
	FieldType       string
	FieldTypeUpper  string
	Tags            string
	IsReference     bool
	IsReferenceList bool
	IsReadonly      bool
}

type OtherDecls struct {
	Consts      []string
	GenDecls    []string
	PackageName string
}

func ParsePackage(path string) ([]*StructTypeDefinition, OtherDecls) {
	set := token.NewFileSet()
	packs, err := parser.ParseDir(set, path, nil, 0)
	assertDirParsed(err)

	definedTypes := make(map[string]*StructTypeDefinition)
	otherTypesDecls := OtherDecls{}

	for _, pack := range packs {
		for _, f := range pack.Files {
			for _, decl := range f.Decls {

				if genDecl, ok := decl.(*ast.GenDecl); ok {
					for _, elem := range genDecl.Specs {
						if typeSpec, ok := elem.(*ast.TypeSpec); ok {
							if strctType, ok := typeSpec.Type.(*ast.StructType); ok {
								typeName := strings.TrimPrefix(typeSpec.Name.Name, "*")

								MakeFieldsUnexported(strctType.Fields)
								if _, ok := definedTypes[typeName]; !ok {
									definedTypes[typeName] = &StructTypeDefinition{}
								}

								definedTypes[typeName].TypeNameOrginalCase = typeName
								definedTypes[typeName].TypeNameLower = lowerCasedFirstChar(typeName)

								parseStructFieldsForClients(definedTypes[typeName], strctType, typeName)
							} else {
								// DETECT AND SAVE CUSTOM TYPES (e.g. type MyInt int)
								def, err := getTypeSpecAsString(set, typeSpec)
								if err != nil {
									fmt.Println(err)
								} else {
									otherTypesDecls.GenDecls = append(otherTypesDecls.GenDecls, def)
								}
							}
						}
					}
					// DETECT AND SAVE CONST DECLARATIONS
					if genDecl.Tok == token.CONST {
						constStr, err := getConstAsString(set, genDecl)
						if err != nil {
							fmt.Println(err)
						}
						otherTypesDecls.Consts = append(otherTypesDecls.Consts, constStr)
						continue
					}
				}

				if fn, isFn := decl.(*ast.FuncDecl); isFn {
					if fn.Recv == nil || fn.Name.Name == InitFunctionName {
						continue
					}

					typeName := strings.TrimPrefix(types.ExprString(fn.Recv.List[0].Type), "*")
					if fn.Name.Name == NobjectImplementationMethod {
						funcString, err := getFunctionBodyAsString(set, fn.Body)
						if err != nil {
							fmt.Println("error occurred when parsing GetTypeName of " + typeName)
							continue
						}

						if _, ok := definedTypes[typeName]; !ok {
							definedTypes[typeName] = &StructTypeDefinition{}
						}
						definedTypes[typeName].NobjectImplementation = funcString
						continue
					}

					if fn.Name.Name == CustomIdImplementationMethod {
						funcString, err := getFunctionBodyAsString(set, fn.Body)
						if err != nil {
							fmt.Println("error occurred when parsing GetTypeName of " + typeName)
							continue
						}

						if _, ok := definedTypes[typeName]; !ok {
							definedTypes[typeName] = &StructTypeDefinition{}
						}
						definedTypes[typeName].CustomIdImplementation = funcString
						if len(fn.Recv.List[0].Names) > 0 {
							definedTypes[typeName].CustomIdReceiverName = fn.Recv.List[0].Names[0].Name
						}
						continue
					}

					if isGetterOrSetterMethod(fn, typeName, definedTypes) {
						continue
					}

					memberFunction, err := PrepareMemberFunction(fn)
					if err != nil {
						fmt.Println("Function "+fn.Name.Name+"not generated in client lib", err)
						continue
					}

					if _, ok := definedTypes[typeName]; !ok {
						definedTypes[typeName] = &StructTypeDefinition{}
					}
					definedTypes[typeName].MemberFunctions = append(definedTypes[typeName].MemberFunctions, *memberFunction)
				}
			}
		}
	}

	detectAndSetNobjectsReturnTypes(definedTypes)
	return maps.Values(definedTypes), otherTypesDecls
}

func isGetterOrSetterMethod(fn *ast.FuncDecl, typeName string, definedTypes map[string]*StructTypeDefinition) bool {
	if strings.HasPrefix(fn.Name.Name, GetPrefix) {
		fieldName := strings.TrimPrefix(fn.Name.Name, GetPrefix)

		for _, field := range definedTypes[typeName].FieldDefinitions {
			if field.FieldNameUpper == fieldName {
				return true
			}
		}
	} else if strings.HasPrefix(fn.Name.Name, SetPrefix) {
		fieldName := strings.TrimPrefix(fn.Name.Name, SetPrefix)

		for _, field := range definedTypes[typeName].FieldDefinitions {
			if field.FieldNameUpper == fieldName {
				return true
			}
		}
	}

	return false
}

func detectAndSetNobjectsReturnTypes(definedTypes map[string]*StructTypeDefinition) {
	for _, typeDefinition := range definedTypes {
		for i, function := range typeDefinition.MemberFunctions {
			if isReturnTypeDefined(function) && isReturnTypeNobject(function, definedTypes) {
				typeDefinition.MemberFunctions[i].IsReturnTypeNobject = true
				typeDefinition.MemberFunctions[i].OptionalReturnTypeUpper = function.OptionalReturnType
				typeDefinition.MemberFunctions[i].OptionalReturnType = lowerCasedFirstChar(function.OptionalReturnType)
			}
		}
	}
}

func isReturnTypeDefined(f MethodDefinition) bool {
	return f.OptionalReturnType != ""
}

func isReturnTypeNobject(f MethodDefinition, defTypes map[string]*StructTypeDefinition) bool {
	return defTypes[f.OptionalReturnType] != nil && defTypes[f.OptionalReturnType].NobjectImplementation != ""
}

func parseStructFieldsForClients(structDef *StructTypeDefinition, astStrct *ast.StructType, typeName string) {
	if astStrct == nil || astStrct.Fields == nil || astStrct.Fields.List == nil {
		return
	}

	for _, field := range astStrct.Fields.List {
		fieldType := strings.TrimPrefix(types.ExprString(field.Type), "*")

		if strings.Contains(fieldType, LibraryReferenceNavigationList) {
			err := parseRelationshipsTagsClient(structDef, field, typeName, fieldType)
			if err != nil {
				fmt.Println("error occurred when parsing relationship tags", err)
			}

		} else {
			parseFields(field, fieldType, structDef)
		}
	}
}

func parseRelationshipsTagsClient(structDef *StructTypeDefinition, field *ast.Field, typeName string, fieldType string) error {
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
			structDef.OneToManyRelationships = append(structDef.OneToManyRelationships, oneToMany)
		} else if strings.Contains(tag.Name, HasManyTag) {
			navigationToTypeName := strings.TrimPrefix(fieldType, LibraryReferenceNavigationList)
			navigationToTypeName = strings.Trim(navigationToTypeName, "[]")

			newManyToManyRelationship := NewManyToManyRelationshipField(typeName, navigationToTypeName, field.Names[0].Name)
			newManyToManyRelationship.FromFieldName = field.Names[0].Name
			structDef.ManyToManyRelationships = append(structDef.ManyToManyRelationships, *newManyToManyRelationship)
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

func PrepareMemberFunction(fn *ast.FuncDecl) (*MethodDefinition, error) {

	if fn.Type.Results == nil ||
		types.ExprString(fn.Type.Results.List[len(fn.Type.Results.List)-1].Type) != "error" {
		return nil, fmt.Errorf("methods belonging to nobjects must return error type")
	}
	if len(fn.Type.Results.List) > 2 {
		return nil, fmt.Errorf("methods belonging to nobjects can return at most 2 variables")
	}

	memberFunction := MethodDefinition{
		FuncName: fn.Name.Name,
	}

	if len(fn.Recv.List[0].Names) > 0 {
		memberFunction.ReceiverName = fn.Recv.List[0].Names[0].Name
	}

	if len(fn.Type.Results.List) > 1 {
		memberFunction.OptionalReturnType = types.ExprString(fn.Type.Results.List[0].Type)
	}

	if len(fn.Type.Params.List) > 1 {
		return nil, fmt.Errorf("methods belonging to nobjects can have at most 1 parameter")
	} else if len(fn.Type.Params.List) == 1 {
		memberFunction.InputParamType = types.ExprString(fn.Type.Params.List[0].Type)
	}

	return &memberFunction, nil
}

func MakeFieldsUnexported(fieldList *ast.FieldList) {
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

func getTypeSpecAsString(fset *token.FileSet, detectedStruct *ast.TypeSpec) (string, error) {
	var buf bytes.Buffer
	buf.WriteString("type ")
	err := printer.Fprint(&buf, fset, detectedStruct)
	if err != nil {
		return "", fmt.Errorf("error occurred when parsing the struct")
	}
	return buf.String(), nil
}

func getConstAsString(fset *token.FileSet, detectedGenDecl *ast.GenDecl) (string, error) {
	var buf bytes.Buffer
	err := printer.Fprint(&buf, fset, detectedGenDecl)
	if err != nil {
		return "", fmt.Errorf("error occurred when parsing the struct")
	}
	return buf.String(), nil
}

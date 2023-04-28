package parser

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"github.com/fatih/structtag"
)

type structPath struct {
	strctType *ast.StructType
	path      string
}

func (t *TypeSpecParser) detectAndModifyAstStructs() {

	detectedStructTypeWithFile := map[string]structPath{}

	for _, pack := range t.packs {
		for path, f := range pack.Files {
			t.fileChanged[path] = false
			for _, d := range f.Decls {
				if genDecl, ok := d.(*ast.GenDecl); ok {
					for _, elem := range genDecl.Specs {
						if typeSpec, ok := elem.(*ast.TypeSpec); ok {
							typeName := strings.TrimPrefix(typeSpec.Name.Name, "*")

							if strctType, ok := typeSpec.Type.(*ast.StructType); ok {
								detectedStructTypeWithFile[typeName] = structPath{strctType: strctType, path: path}
								modified := t.parseStructFields(f, strctType, typeName)

								if !t.fileChanged[path] {
									t.fileChanged[path] = modified
								}

								if t.Output.IsNobjectInOrginalPackage[typeName] && !t.isInitAlreadyAdded[typeName] {
									t.addInitMethod(f, typeName)
									t.fileChanged[path] = true
								}

								if t.Output.IsNobjectInOrginalPackage[typeName] && !t.isSaveChangesAlreadyAdded[typeName] {

									t.addSaveChangesIfInitializedMethod(f, typeName)
									t.fileChanged[path] = true
								}
							}
						}
					}
				}
			}
		}
	}

	t.addIgnoreEmptyTagToBidirectionalOneToManyRel(detectedStructTypeWithFile)
}

// The addIgnoreEmptyTagToBidirectionalOneToManyRel method adds `dynamodbav:",omitempty"` tag
// to the side with Reference type of bidirectional one-to-many relationship
// without the tag, every Reference field one the side of bidirectional
// one-to-many relationship is mandatory
// for in-depth exaplanation, see: https://github.com/aws/aws-sdk-go/issues/1803 and
// https://notes.serverlessfirst.com/public/How+does+DynamoDB+handle+NULL%2C+empty+and+undefined+fields
func (t *TypeSpecParser) addIgnoreEmptyTagToBidirectionalOneToManyRel(detectedStructTypeWithFile map[string]structPath) {
	for _, oneToManyRels := range t.Output.BidrectionalOneToManyRel {
		for _, oneToMany := range oneToManyRels {
			strctWithReferenceField := detectedStructTypeWithFile[oneToMany.TypeName]
			for _, field := range strctWithReferenceField.strctType.Fields.List {
				if field.Names[0].Name == oneToMany.FieldName {
					// field to which the oneToMany relationships refers found

					tags, err := getParsedTags(field)
					if err != nil {
						fmt.Println("error occurerd while checking struct tags of:", oneToMany.TypeName, " field: ", field.Names[0].Name, ". Error: ", err)
					} else if tags != nil {
						dynamodbTag, err := tags.Get(DynamoDBTagKey)
						if err != nil {
							fmt.Println("error occurerd while checking struct tags of:", oneToMany.TypeName, " field: ", field.Names[0].Name, ". Error: ", err)
						} else if dynamodbTag == nil {
							// no dynamoDB tags added before

							tags.AddOptions(DynamoDBTagKey, DynamoDBIgnoreEmptyTagValue)
							field.Tag.Value = "`" + tags.String() + "`"
							t.fileChanged[strctWithReferenceField.path] = true
						} else if !strings.Contains(dynamodbTag.GoString(), DynamoDBIgnoreEmptyTagValue) {
							// dynamoDB tags added before, but not omniempty option

							dynamodbTag.Options = append(dynamodbTag.Options, DynamoDBIgnoreEmptyTagValue)
							field.Tag.Value = "`" + tags.String() + "`"
							t.fileChanged[strctWithReferenceField.path] = true
						}
					} else {
						// no tags added to field

						field.Tag = &ast.BasicLit{Value: "`" + DynamoDBIgnoreEmptyTag + "`"}
						t.fileChanged[strctWithReferenceField.path] = true
					}
					break
				}
			}
		}
	}
}

func (t *TypeSpecParser) addInitMethod(f *ast.File, typeName string) {
	idFieldName := getIdFieldNameOfType(typeName, t.Output.TypesWithCustomId)
	function := getInitFunctionForType(typeName, idFieldName, t.Output.BidrectionalOneToManyRel[typeName], t.Output.ManyToManyRelationships[typeName])
	f.Decls = append(f.Decls, function)
}

func (t *TypeSpecParser) addSaveChangesIfInitializedMethod(f *ast.File, typeName string) {
	function := getSaveChangesMethodForType(typeName, t.Output)
	f.Decls = append(f.Decls, function)
}

// The parseStructFields returns true if the ast representing
// the struct was modified, otherwise false
func (t *TypeSpecParser) parseStructFields(f *ast.File, strctType *ast.StructType, typeName string) bool {
	structDefinitionModified := false
	if strctType == nil || strctType.Fields == nil || len(strctType.Fields.List) == 0 {
		return structDefinitionModified
	}

	t.Output.TypeFields[typeName] = make(map[string]string, len(strctType.Fields.List))
	fieldModified, structModified := false, false

	isNobject := t.Output.IsNobjectInOrginalPackage[typeName]
	for _, field := range strctType.Fields.List {
		t.Output.TypeFields[typeName][field.Names[0].Name] = types.ExprString(field.Type)

		if isNobject {
			fieldModified = t.parseRelationshipsTags(field, typeName)
			structModified = t.addCustomIdImplementationIfNeeded(f, field, typeName)

			if !structDefinitionModified {
				structDefinitionModified = fieldModified || structModified
			}
		}
	}

	if _, exists := t.Output.TypeFields[typeName][IsInitializedFieldName]; !exists && isNobject {
		strctType.Fields.List = append(strctType.Fields.List, &ast.Field{
			Names: []*ast.Ident{{Name: IsInitializedFieldName}}, Type: &ast.Ident{Name: "bool"},
		})
		structDefinitionModified = true
	}
	if _, exists := t.Output.TypeFields[typeName][InvocationDepthFieldName]; !exists && isNobject {
		strctType.Fields.List = append(strctType.Fields.List, &ast.Field{
			Names: []*ast.Ident{{Name: InvocationDepthFieldName}}, Type: &ast.Ident{Name: "int"},
		})
		structDefinitionModified = true
	}

	return structDefinitionModified
}

func (t *TypeSpecParser) addCustomIdImplementationIfNeeded(f *ast.File, field *ast.Field, typeName string) bool {
	tags, err := getParsedTags(field)

	if err != nil {
		fmt.Println("error occurerd while checking struct tags of:", typeName, " field: ", field.Names[0].Name, ". Error: ", err)
	} else if tags != nil {
		if tag, _ := tags.Get(NubesTagKey); tag != nil {

			if strings.EqualFold(tag.Name, CustomIdTag) {
				if types.ExprString(field.Type) != "string" {
					fmt.Println("ERROR: The field selected as CustomId field must be a string.", field.Names[0].Name,
						"selected as CustomId field selected for type", typeName, "is not a string")
					return false
				}

				tagAdded := addDynamoDBIdTag(tags, typeName, field)

				if fieldName, exists := t.Output.TypesWithCustomId[typeName]; exists {
					if fieldName != field.Names[0].Name {
						fmt.Println(`ERROR: already existing implementation of CustomId interface (GetId method) 
						must be removed after different field is set to be the CustomId. Old  CustomId field: `,
							field.Names[0].Name, "the new one:", field.Names[0].Name)
					}
					return tagAdded
				}

				t.Output.TypesWithCustomId[typeName] = field.Names[0].Name
				f.Decls = append(f.Decls, getCustomIdImplementation(typeName, field.Names[0].Name))
				return true
			}
		}
	}

	return false
}

func addDynamoDBIdTag(tags *structtag.Tags, typeName string, field *ast.Field) bool {
	dynamodbTag, _ := tags.Get(DynamoDBTagKey)

	if dynamodbTag != nil && dynamodbTag.Name == DynamoDBIdTagValue {
		return false
	}

	if dynamodbTag != nil && dynamodbTag.Name != DynamoDBIdTagValue {
		fmt.Println("invalid definition of dynamodb struct tag fixed in", typeName, "field:", field.Names[0].Name, " replaced with mandatory ignore tag for", LibraryReferenceNavigationList)
	}
	field.Tag.Value = field.Tag.Value[0:len(field.Tag.Value)-1] + " " + DynamoDBIdTag + "`"
	return true
}

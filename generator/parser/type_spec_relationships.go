package parser

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/fatih/structtag"
)

type ManyToManyRelationshipField struct {
	FieldName      string
	FromFieldName  string
	PartionKeyName string
	SortKeyName    string
	TableName      string
}

type OneToManyRelationshipField struct {
	TypeName      string
	FieldName     string
	FromFieldName string
}

func NewManyToManyRelationshipField(typeName1, typeName2, fieldName string) *ManyToManyRelationshipField {
	// aproach: partion key id is always the "smaller" string
	// where "smaller" means: the ASCII number of the first distinct character
	// corresponds to lower value, or the string is shorter
	// if one typeName is the Prefix of another
	result := new(ManyToManyRelationshipField)

	for index := 0; ; index++ {

		if index >= len(typeName1) {
			result.PartionKeyName = typeName1
			result.SortKeyName = typeName2
			break
		}
		if index >= len(typeName2) {
			result.PartionKeyName = typeName2
			result.SortKeyName = typeName1
			break
		}

		if typeName1[index] < typeName2[index] {
			result.PartionKeyName = typeName1
			result.SortKeyName = typeName2
			break
		} else if typeName1[index] > typeName2[index] {
			result.PartionKeyName = typeName2
			result.SortKeyName = typeName1
			break
		}
	}

	result.FieldName = fieldName
	result.TableName = result.PartionKeyName + result.SortKeyName
	return result
}

// The parseRelationshipsTags detects Nubes' HasManyTag and HasOneTag tags.
// If tag is found, it adds dynamodbav:"-" tag
// so that the field is ignored in dynamodb interaction.
// If the dynamodb tag was already added, it does nothing.
// The parseRelationshipsTags return value indicates whether
// the ast was modified (= whether the dynamodb tag was added).
func parseRelationshipsTags(field *ast.Field, typeName string, fieldType string, parsedPackage *ParsedPackage) bool {
	tags, err := getParsedTags(field)

	if err != nil {
		fmt.Println("error occured while checking struct tags of:", typeName, " field: ", field.Names[0].Name, ". Error: ", err)
	} else if tags != nil {
		if tag, _ := tags.Get(NubesTagKey); tag != nil {

			if strings.Contains(tag.Name, HasOneTag) {

				splitted := strings.Split(tag.Name, "-")
				navigationToFieldName := ""
				if len(splitted) > 0 {
					navigationToFieldName = splitted[1]
				} else {
					fmt.Println(HasOneTag, " detected, but missing reffering field name. Referring field name should be specified after - charcter, e.g.: ", HasOneTag, "<referring_field_name>")
					return false
				}

				if strings.Contains(fieldType, LibraryReferenceNavigationList) {
					navigationToTypeName := strings.TrimPrefix(fieldType, LibraryReferenceNavigationList)
					navigationToTypeName = strings.Trim(navigationToTypeName, "[]")

					parsedPackage.TypeAttributesIndexes[navigationToTypeName] = append(parsedPackage.TypeAttributesIndexes[navigationToTypeName], navigationToFieldName)
					navToField := OneToManyRelationshipField{TypeName: navigationToTypeName, FieldName: navigationToFieldName, FromFieldName: field.Names[0].Name}
					parsedPackage.TypeNavListsReferringFieldName[typeName] = append(parsedPackage.TypeNavListsReferringFieldName[typeName], navToField)

					return addDynamoDBIgnoreTag(tags, field, typeName)
				} else {
					fmt.Println(HasManyTag, " or ", HasOneTag, " can be used only with ", LibraryReferenceNavigationList, " fields!")
					return false
				}
			} else if strings.Contains(tag.Name, HasManyTag) {

				if strings.Contains(fieldType, LibraryReferenceNavigationList) {
					navigationToTypeName := strings.TrimPrefix(fieldType, LibraryReferenceNavigationList)
					navigationToTypeName = strings.Trim(navigationToTypeName, "[]")

					newManyToManyRelationship := NewManyToManyRelationshipField(typeName, navigationToTypeName, field.Names[0].Name)
					newManyToManyRelationship.FromFieldName = field.Names[0].Name
					parsedPackage.ManyToManyRelationships[typeName] = append(parsedPackage.ManyToManyRelationships[typeName], *newManyToManyRelationship)

					return addDynamoDBIgnoreTag(tags, field, typeName)
				} else {
					fmt.Println(HasManyTag, " or ", HasOneTag, " can be used only with ", LibraryReferenceNavigationList, " fields!")
					return false
				}
			}
		}
	}

	return false
}

func addDynamoDBIgnoreTag(tags *structtag.Tags, field *ast.Field, typeName string) bool {
	dynamoTag, _ := tags.Get(DynamoDBKeyTag)
	if dynamoTag == nil {
		field.Tag.Value = field.Tag.Value[0:len(field.Tag.Value)-1] + " " + DynamoDBIgnoreTag + "`"
		return true
	}
	if dynamoTag.Name != "-" {
		fmt.Println("invalid definition of dynamodb struct tag fixed in", typeName, "field:", field.Names[0].Name, " replaced with mandatory ignore tag for", LibraryReferenceNavigationList)
		tags.Set(&structtag.Tag{Key: DynamoDBKeyTag, Name: DynamoDBIgnoreValueTag})
		field.Tag.Value = "`" + tags.String() + "`"
		return true
	}

	return false
}

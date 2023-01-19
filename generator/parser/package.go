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

	"github.com/fatih/structtag"
)

type ParsedPackage struct {
	ImportPath                     string
	IsNobjectInOrginalPackage      map[string]bool
	TypeFields                     map[string]map[string]string
	TypeAttributesIndexes          map[string][]string
	TypeNavListsReferringFieldName map[string][]NavigationToField
	ManyToManyRelationships        map[string][]ManyToManyRelationshipField
	TypesWithCustomId              map[string]string
}

type ManyToManyRelationshipField struct {
	FieldName      string
	PartionKeyName string
	SortKeyName    string
	TableName      string
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

type NavigationToField struct {
	TypeName  string
	FieldName string
}

func GetPackageTypes(path string, moduleName string) ParsedPackage {
	set := token.NewFileSet()
	packs, err := parser.ParseDir(set, path, nil, 0)
	assertDirParsed(err)

	result := ParsedPackage{
		IsNobjectInOrginalPackage:      make(map[string]bool),
		TypesWithCustomId:              map[string]string{},
		TypeAttributesIndexes:          map[string][]string{},
		TypeNavListsReferringFieldName: map[string][]NavigationToField{},
		ManyToManyRelationships:        map[string][]ManyToManyRelationshipField{},
		TypeFields:                     map[string]map[string]string{},
	}

	for packageName, pack := range packs {
		for path, f := range pack.Files {
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
								fileModified := parseStructFields(strctType, typeName, &result)
								if fileModified {
									saveAstChangesInFile(f, set, path)
								}
							}
						}
					}
				}
			}
		}
		result.ImportPath = moduleName + "/" + packageName
	}

	return result
}

func parseStructFields(strctType *ast.StructType, typeName string, parsedPackage *ParsedPackage) bool {
	structDefinitionModified := false
	if strctType == nil || strctType.Fields == nil || len(strctType.Fields.List) == 0 {
		return structDefinitionModified
	}

	parsedPackage.TypeFields[typeName] = make(map[string]string, len(strctType.Fields.List))
	for _, field := range strctType.Fields.List {
		fieldType := types.ExprString(field.Type)
		parsedPackage.TypeFields[typeName][field.Names[0].Name] = fieldType

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
						continue
					}

					if strings.Contains(fieldType, LibraryReferenceNavigationList) {
						navigationToTypeName := strings.TrimPrefix(fieldType, LibraryReferenceNavigationList)
						navigationToTypeName = strings.Trim(navigationToTypeName, "[]")

						parsedPackage.TypeAttributesIndexes[navigationToTypeName] = append(parsedPackage.TypeAttributesIndexes[navigationToTypeName], navigationToFieldName)
						parsedPackage.TypeNavListsReferringFieldName[typeName] = append(parsedPackage.TypeNavListsReferringFieldName[typeName], NavigationToField{TypeName: navigationToTypeName, FieldName: navigationToFieldName})

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
					} else {
						fmt.Println(HasManyTag, " or ", HasOneTag, " can be used only with ", LibraryReferenceNavigationList, " fields!")
					}
				} else if strings.Contains(tag.Name, HasManyTag) {

					if strings.Contains(fieldType, LibraryReferenceNavigationList) {
						navigationToTypeName := strings.TrimPrefix(fieldType, LibraryReferenceNavigationList)
						navigationToTypeName = strings.Trim(navigationToTypeName, "[]")

						newManyToManyRelationship := NewManyToManyRelationshipField(typeName, navigationToTypeName, field.Names[0].Name)
						parsedPackage.ManyToManyRelationships[typeName] = append(parsedPackage.ManyToManyRelationships[typeName], *newManyToManyRelationship)
					} else {
						fmt.Println(HasManyTag, " or ", HasOneTag, " can be used only with ", LibraryReferenceNavigationList, " fields!")
					}
				}
			}
		}
	}
	return false
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

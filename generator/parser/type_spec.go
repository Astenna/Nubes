package parser

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"go/types"
	"os"
	"strings"

	tp "github.com/Astenna/Nubes/generator/template"
)

type TypeSpecParser struct {
	Output      ParsedPackage
	Handlers    []StateChangingHandler
	CustomCtors []CustomCtorDefinition

	tokenSet                  *token.FileSet
	packs                     map[string]*ast.Package
	detectedFunctions         map[string][]detectedFunction
	isSaveChangesAlreadyAdded map[string]bool
	isInitAlreadyAdded        map[string]bool
	fileChanged               map[string]bool
}

type ParsedPackage struct {
	ImportPath                string
	IsNobjectInOrginalPackage map[string]bool
	TypeFields                map[string]map[string]string
	TypeAttributesIndexes     map[string][]string
	BidrectionalOneToManyRel  map[string][]OneToManyRelationshipField
	ManyToManyRelationships   map[string][]ManyToManyRelationshipField
	TypesWithCustomId         map[string]string
	TypesWithCustomExport     map[string]CustomExportDefinition
	TypesWithCustomDelete     map[string]CustomDeleteDefinition
}

type CustomCtorDefinition struct {
	OrginalPackageAlias    string
	OrginalPackage         string
	TypeName               string
	OptionalParamType      string
	IsOptionalParamNobject bool
}

type CustomExportDefinition struct {
	InputParameterType string
}

type CustomDeleteDefinition struct {
	InputParameterType string
}

func NewTypeSpecParser(path string) (*TypeSpecParser, error) {
	typeSpecParser := new(TypeSpecParser)
	typeSpecParser.tokenSet = token.NewFileSet()
	packg, err := parser.ParseDir(typeSpecParser.tokenSet, path, nil, parser.Mode(parser.ParseComments))
	if err != nil {
		return nil, fmt.Errorf("failed to parse package in path %s. Error: %w", path, err)
	}

	typeSpecParser.packs = packg
	typeSpecParser.Output = ParsedPackage{
		IsNobjectInOrginalPackage: make(map[string]bool),
		TypesWithCustomId:         map[string]string{},
		TypesWithCustomExport:     map[string]CustomExportDefinition{},
		TypesWithCustomDelete:     map[string]CustomDeleteDefinition{},
		TypeAttributesIndexes:     map[string][]string{},
		BidrectionalOneToManyRel:  map[string][]OneToManyRelationshipField{},
		ManyToManyRelationships:   map[string][]ManyToManyRelationshipField{},
		TypeFields:                map[string]map[string]string{},
	}
	typeSpecParser.Handlers = []StateChangingHandler{}
	typeSpecParser.CustomCtors = []CustomCtorDefinition{}
	typeSpecParser.fileChanged = map[string]bool{}
	typeSpecParser.detectedFunctions = make(map[string][]detectedFunction)
	typeSpecParser.isInitAlreadyAdded = map[string]bool{}
	typeSpecParser.isSaveChangesAlreadyAdded = map[string]bool{}

	return typeSpecParser, nil
}

func (t *TypeSpecParser) Run(moduleName string) {

	t.detectNobjectTypesAndFunctions(moduleName)
	t.detectAndModifyAstStructs()
	t.modifyAstMethods()
	t.prepareDataForHandlers()
	t.addNubesLibImportIfMissing()
	t.saveChangesInAst()
}

// The detectNobjectTypesAndFunctions detects object types
// and methods defined in the package.
// Nobject types are recognised as the types that implement
// Nobject interface (i.e. GetTypeName method)
func (t *TypeSpecParser) detectNobjectTypesAndFunctions(moduleName string) {
	for packageName, pack := range t.packs {
		for path, f := range pack.Files {
			for _, d := range f.Decls {
				if fn, isFn := d.(*ast.FuncDecl); isFn {

					if fn.Recv != nil {
						ownerType := getFunctionReceiverTypeAsString(fn.Recv)
						switch fn.Name.Name {
						case SaveChangesIfInitialized:
							t.isSaveChangesAlreadyAdded[ownerType] = true
							continue
						case NobjectImplementationMethod:
							t.Output.IsNobjectInOrginalPackage[ownerType] = true
							continue
						case InitFunctionName:
							t.isInitAlreadyAdded[ownerType] = true
							continue
						case CustomIdImplementationMethod:
							idFieldName, err := getIdFieldNameFromCustomIdImpl(fn)
							if err != nil {
								fmt.Println(err)
								continue
							}
							t.Output.TypesWithCustomId[ownerType] = idFieldName
							continue
						}
					}

					// ignore unexported functions (i.e. starting with lowercase letter)
					if fn.Name.IsExported() {
						if areReturnParamsValid(fn) {
							t.detectedFunctions[path] = append(t.detectedFunctions[path], detectedFunction{
								Function: fn,
								Imports:  f.Imports,
							})
						}
					}
				}
			}
		}

		t.Output.ImportPath = moduleName + "/" + packageName
	}
}

func (t TypeSpecParser) addNubesLibImportIfMissing() {
	for _, pack := range t.packs {
		for path, f := range pack.Files {
			// the import is added if missing only to the modified files
			// it is assumed that not modified ones do not require
			// the library as it was already added or is not needed at all
			if value, exists := t.fileChanged[path]; exists && value {
				libImported := false
				for _, imp := range f.Imports {
					if strings.Contains(imp.Path.Value, LibImportPath) {
						libImported = true
						break
					}
				}
				if !libImported {
					importNubes := &ast.GenDecl{
						TokPos: f.Package,
						Tok:    token.IMPORT,
						Specs:  []ast.Spec{&ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: LibImportPath}}},
					}
					f.Decls = prependElem[ast.Decl](f.Decls, importNubes)
				}
			}
		}
	}
}

func (t TypeSpecParser) saveChangesInAst() {
	for _, pack := range t.packs {
		for path, f := range pack.Files {

			if value, exists := t.fileChanged[path]; exists && value {

				var buf bytes.Buffer
				err := printer.Fprint(&buf, t.tokenSet, f)
				if err != nil {
					fmt.Println(err)
				}
				nobjectTypeFile, err := os.Create(path)
				if err != nil {
					fmt.Println(err)
				}
				buf.WriteTo(nobjectTypeFile)
				nobjectTypeFile.Close()
				tp.RunGoimportsOnFile(path)
			}
		}
	}
}

// areReturnParamsValid returns true if the number of parameters is equal to two or one,
// If exactly two return parameters are defined, then the second paramater
// must be an error type.
// If exactly one return parameter is defined, then the parameter must be
// an error type.
// If the above conditions do not hold, it prints relevant error message
// and returns false.
func areReturnParamsValid(f *ast.FuncDecl) bool {

	if f.Type.Results == nil || f.Type.Results.List == nil || !isErrorTypeReturned(f) {
		fmt.Println("error type must be defined as the last return type from type's method. Handler generation for " + f.Name.Name + " skipped")
		return false
	}
	if len(f.Type.Results.List) > 2 {
		fmt.Println("maximum allowed number of non-error return parameters is 1. Handler generation for " + f.Name.Name + " skipped")
		return false
	}

	return true
}

func isErrorTypeReturned(f *ast.FuncDecl) bool {
	return len(f.Type.Results.List) > 0 && types.ExprString(f.Type.Results.List[len(f.Type.Results.List)-1].Type) == "error"
}

func getFunctionReceiverTypeAsString(fieldList *ast.FieldList) string {
	return strings.TrimPrefix(types.ExprString(fieldList.List[0].Type), "*")
}

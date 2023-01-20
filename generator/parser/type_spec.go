package parser

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"strings"
)

type TypeSpecParser struct {
	Output   ParsedPackage
	Handlers []StateChangingHandler

	tokenSet          *token.FileSet
	packs             map[string]*ast.Package
	detectedFunctions map[string][]detectedFunction
	fileChanged       map[string]bool
}

type ParsedPackage struct {
	ImportPath                     string
	IsNobjectInOrginalPackage      map[string]bool
	TypeFields                     map[string]map[string]string
	TypeAttributesIndexes          map[string][]string
	TypeNavListsReferringFieldName map[string][]NavigationToField
	ManyToManyRelationships        map[string][]ManyToManyRelationshipField
	TypesWithCustomId              map[string]string
}

func NewTypeSpecParser(path string) (*TypeSpecParser, error) {
	typeSpecParser := new(TypeSpecParser)
	typeSpecParser.tokenSet = token.NewFileSet()
	packg, err := parser.ParseDir(typeSpecParser.tokenSet, path, nil, 0)
	typeSpecParser.packs = packg

	typeSpecParser.Output = ParsedPackage{
		IsNobjectInOrginalPackage:      make(map[string]bool),
		TypesWithCustomId:              map[string]string{},
		TypeAttributesIndexes:          map[string][]string{},
		TypeNavListsReferringFieldName: map[string][]NavigationToField{},
		ManyToManyRelationships:        map[string][]ManyToManyRelationshipField{},
		TypeFields:                     map[string]map[string]string{},
	}
	typeSpecParser.Handlers = []StateChangingHandler{}
	typeSpecParser.fileChanged = map[string]bool{}
	typeSpecParser.detectedFunctions = make(map[string][]detectedFunction)

	if err != nil {
		return nil, fmt.Errorf("failed to parse package in path %s. Error: %w", path, err)
	}
	return typeSpecParser, nil
}

func (t TypeSpecParser) Run(moduleName string) ParsedPackage {

	t.detectNobjectTypes(moduleName)
	t.detectAndAdjustDecls()

	isTypeNewCtorImplemented := make(map[string]bool)
	isTypeReNewCtorImplemented := make(map[string]bool)
	isTypeDestructorImplemented := make(map[string]bool)

	t.detectAndAdjustMethods(isTypeNewCtorImplemented, isTypeReNewCtorImplemented, isTypeDestructorImplemented)
	t.prepareHandleres()
	t.saveChangesInAst()

	return t.Output
}

func (t *TypeSpecParser) detectNobjectTypes(moduleName string) {
	for packageName, pack := range t.packs {
		for _, f := range pack.Files {
			for _, d := range f.Decls {
				if fn, isFn := d.(*ast.FuncDecl); isFn {
					if fn.Recv != nil && fn.Name.Name == NobjectImplementationMethod {
						ownerType := getFunctionReceiverTypeAsString(fn.Recv)
						t.Output.IsNobjectInOrginalPackage[ownerType] = true
					}
					if fn.Recv != nil && fn.Name.Name == CustomIdImplementationMethod {
						ownerType := getFunctionReceiverTypeAsString(fn.Recv)
						idFieldName, err := getIdFieldNameFromCustomIdImpl(fn)
						if err != nil {
							fmt.Println(err)
							continue
						}
						t.Output.TypesWithCustomId[ownerType] = idFieldName
					}
				}
			}
		}

		// TODO: recognize CustomIds via fieldType
		// TODO: what if there is more than one package?
		t.Output.ImportPath = moduleName + "/" + packageName
	}
}

func (t TypeSpecParser) saveChangesInAst() {
	for _, pack := range t.packs {
		for path, f := range pack.Files {

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
					f.Decls = prepend[ast.Decl](f.Decls, importNubes)
				}

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
			}
		}
	}
}

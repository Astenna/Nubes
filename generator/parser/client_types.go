package parser

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"go/types"
	"os"
	"strings"
)

type TypeDefinition struct {
	PackageName      string
	Imports          string
	StructDefinition string
	MemberFunctions  []MemberFunction
}

type MemberFunction struct {
	ReceiverName       string
	ReceiverType       string
	FuncName           string
	InputParamName     string
	InputParamType     string
	OptionalReturnType string
}

func PrepareTypes(path string) []TypeDefinition {
	set := token.NewFileSet()
	packs, err := parser.ParseDir(set, path, nil, 0)
	AssertDirParsed(err)

	var structs []*ast.StructType
	funcsMap := make(map[string][]*ast.FuncDecl)

	for _, pack := range packs {
		for _, f := range pack.Files {

			ast.Inspect(f, func(n ast.Node) bool {
				if typeSpec, ok := n.(*ast.TypeSpec); ok {
					if strctType, ok := typeSpec.Type.(*ast.StructType); ok {
						structs = append(structs, strctType)
						MakeFieldsUnexported(strctType.Fields)
						printer.Fprint(os.Stdout, set, strctType)
					}
				}
				return true
			})

			for _, d := range f.Decls {
				if fn, isFn := d.(*ast.FuncDecl); isFn {
					if fn.Recv == nil {
						continue
					}
					ownerType := types.ExprString(fn.Recv.List[0].Type)
					if funcsMap[ownerType] == nil {
						funcsMap[ownerType] = []*ast.FuncDecl{}
					}
					funcsMap[ownerType] = append(funcsMap[ownerType], fn)
				}
			}
		}
	}

	memberFuncsMap := make(map[string][]MemberFunction)
	for ownerType, funcs := range funcsMap {
		for _, f := range funcs {
			if f.Name.Name == GetTypeName {
				continue
			}
			memberFuncsMap[ownerType] = append(memberFuncsMap[ownerType], MemberFunction{
				Name:       f.Name.Name,
				Parameters: "TODO",
			})
		}
	}

	typeDefinitions := []TypeDefinition{}
	for _, strcs := range structs {
		typeDefinitions = append(typeDefinitions, TypeDefinition{
			Imports:          "TODO",
			StructDefinition: types.ExprString(strcs),
			MemberFunctions:  memberFuncsMap[types.ExprString(strcs)],
		})
	}

	return typeDefinitions
}

func MakeFieldsUnexported(fieldList *ast.FieldList) {
	for _, field := range fieldList.List {
		field.Names[0].Name = strings.ToLower(field.Names[0].Name)
	}
}

func PrepareTypesFiles(path string) []TypeDefinition {
	set := token.NewFileSet()
	packs, err := parser.ParseDir(set, path, nil, 0)
	AssertDirParsed(err)

	//structsMap := make(map[string][]*ast.StructType)
	//funcsMap := make(map[string][]*ast.FuncDecl)

	for _, pack := range packs {
		for _, file := range pack.Files {
			ast.Inspect(file, func(n ast.Node) bool {
				switch x := n.(type) {
				case *ast.FuncDecl:

					newCallStmt := &ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "fmt",
								},
								Sel: &ast.Ident{
									Name: "Println",
								},
							},
							Args: []ast.Expr{
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: `"instrumentation"`,
								},
							},
						},
					}

					x.Body.List = []ast.Stmt{newCallStmt}
				}
				return true
			})

			printer.Fprint(os.Stdout, set, file)
		}
	}

	return nil
}

// astutil.Apply(pack, func(c *astutil.Cursor) bool {
// 	n := c.Node()
// 	switch n.(type) {
// 	case *ast.FuncDecl:

// 		astutil.Apply(n, func(crs *astutil.Cursor) bool {
// 			if _, ok := crs.Node().(*ast.BlockStmt); ok {
// 				blckStmnt := ast.BlockStmt{
// 					List: {

// 					},
// 				}
// 				c.Replace(blckStmnt)
// 			}
// 			return false
// 		}, nil)

// 	}
// 	return false
// }, nil)

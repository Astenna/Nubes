package parser

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"go/types"
	"os"
)

type TypeDefinition struct {
	PackageName      string
	Import           string
	StructDefinition string
	MemberFunctions  []MemberFunction
}

type MemberFunction struct {
	Name       string
	Parameters string
	Body       string
}

func PrepareTypes(path string) []TypeDefinition {
	set := token.NewFileSet()
	packs, err := parser.ParseDir(set, path, nil, 0)
	AssertDirParsed(err)

	structsMap := make(map[string][]*ast.StructType)
	funcsMap := make(map[string][]*ast.FuncDecl)

	for packageName, pack := range packs {
		for _, f := range pack.Files {

			ast.Inspect(f, func(n ast.Node) bool {
				if str, ok := n.(*ast.StructType); ok {
					if structsMap[packageName] == nil {
						structsMap[packageName] = []*ast.StructType{}
					}
					structsMap[packageName] = append(structsMap[packageName], str)
					printer.Fprint(os.Stdout, set, str)
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
	for packageName, strcs := range structsMap {
		for _, str := range strcs {
			typeDefinitions = append(typeDefinitions, TypeDefinition{
				PackageName:      packageName,
				Import:           "TODO",
				StructDefinition: types.ExprString(str),
				MemberFunctions:  memberFuncsMap[types.ExprString(str)],
			})
		}
	}

	return typeDefinitions
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

package parser

import (
	"go/ast"
)

func getCustomIdImplementation(typeName, fieldName string) *ast.FuncDecl {
	receiverName := "receiver"

	return &ast.FuncDecl{
		Name: &ast.Ident{Name: CustomIdImplementationMethod},
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{{Name: receiverName}},
					Type:  &ast.Ident{Name: typeName},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{&ast.ReturnStmt{
				Results: []ast.Expr{&ast.SelectorExpr{
					X:   &ast.Ident{Name: receiverName},
					Sel: &ast.Ident{Name: fieldName},
				}},
			}},
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
			Results: &ast.FieldList{
				List: []*ast.Field{{Type: &ast.Ident{Name: "string"}}},
			},
		},
	}
}

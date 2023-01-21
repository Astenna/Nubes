package parser

import (
	"go/ast"
	"go/token"
	"go/types"
)

func getErrorCheckExpr(fn *ast.FuncDecl, errorVariableName string) ast.IfStmt {
	ifStmt := ast.IfStmt{
		Cond: &ast.BinaryExpr{
			X:  &ast.Ident{Name: errorVariableName},
			Op: token.NEQ,
			Y:  &ast.Ident{Name: "nil"},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{},
		},
	}

	// given faas handler reqs, there can be one optional return type apart from "error"
	if fn.Type.Results != nil && fn.Type.Results.List != nil {
		returnTypeName := types.ExprString(fn.Type.Results.List[0].Type)
		if returnTypeName != "error" {
			ifStmt.Body.List = []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.StarExpr{
							X: &ast.CallExpr{
								Args: []ast.Expr{&ast.Ident{Name: returnTypeName}},
								Fun:  &ast.Ident{Name: "new"},
							},
						},
						&ast.Ident{
							Name: LibErrorVariableName,
						},
					},
				},
			}

			return ifStmt
		}
	}

	ifStmt.Body.List = []ast.Stmt{
		&ast.ReturnStmt{
			Results: []ast.Expr{
				&ast.Ident{
					Name: LibErrorVariableName,
				},
			},
		},
	}
	return ifStmt
}

type getDBStmtsParam struct {
	idFieldName          string
	typeName             string
	receiverVariableName string
	fieldName            string
	fieldType            string
}

func getGetterDBStmts(fn *ast.FuncDecl, input getDBStmtsParam) ast.IfStmt {
	isInitializedCheck := getIsInitializedCheck(fn.Recv.List[0].Names[0].Name)
	getFieldFromLib := ast.AssignStmt{
		Tok: token.DEFINE,
		Lhs: []ast.Expr{
			&ast.Ident{Name: "fieldValue"},
			&ast.Ident{Name: LibErrorVariableName},
		},
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "lib"},
					Sel: &ast.Ident{Name: "GetField"},
				},
				Args: []ast.Expr{
					&ast.SelectorExpr{
						X:   &ast.Ident{Name: input.receiverVariableName},
						Sel: &ast.Ident{Name: input.idFieldName},
					},
					&ast.CompositeLit{
						Type: &ast.SelectorExpr{
							X:   &ast.Ident{Name: "lib"},
							Sel: &ast.Ident{Name: "GetFieldParam"},
						},
						Elts: []ast.Expr{
							&ast.KeyValueExpr{
								Key: &ast.Ident{Name: "TypeName"},
								Value: &ast.BasicLit{
									Kind:  token.STRING,
									Value: "\"" + input.typeName + "\"",
								},
							},
							&ast.KeyValueExpr{
								Key: &ast.Ident{Name: "FieldName"},
								Value: &ast.BasicLit{
									Kind:  token.STRING,
									Value: "\"" + input.fieldName + "\"",
								},
							},
						},
					},
				},
			}}}
	errorCheck := getErrorCheckExpr(fn, LibErrorVariableName)
	fieldAssign := ast.AssignStmt{
		Tok: token.ASSIGN,
		Lhs: []ast.Expr{
			&ast.SelectorExpr{
				X:   &ast.Ident{Name: input.receiverVariableName},
				Sel: &ast.Ident{Name: input.fieldName},
			},
		},
		Rhs: []ast.Expr{
			&ast.TypeAssertExpr{
				Type: &ast.Ident{Name: input.fieldType},
				X:    &ast.Ident{Name: "fieldValue"},
			},
		},
	}

	isInitializedCheck.Body.List = []ast.Stmt{&getFieldFromLib, &errorCheck, &fieldAssign}
	return isInitializedCheck
}

func getSetterDBStmts(fn *ast.FuncDecl, input getDBStmtsParam) ast.IfStmt {
	isInitializedCheck := getIsInitializedCheck(fn.Recv.List[0].Names[0].Name)
	getFieldFromLib := ast.AssignStmt{
		Tok: token.DEFINE,
		Lhs: []ast.Expr{
			&ast.Ident{Name: LibErrorVariableName},
		},
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "lib"},
					Sel: &ast.Ident{Name: "SetField"},
				},
				Args: []ast.Expr{
					&ast.SelectorExpr{
						X:   &ast.Ident{Name: input.receiverVariableName},
						Sel: &ast.Ident{Name: input.idFieldName},
					},
					&ast.CompositeLit{
						Type: &ast.SelectorExpr{
							X:   &ast.Ident{Name: "lib"},
							Sel: &ast.Ident{Name: "SetFieldParam"},
						},
						Elts: []ast.Expr{
							&ast.KeyValueExpr{
								Key: &ast.Ident{Name: "TypeName"},
								Value: &ast.BasicLit{
									Kind:  token.STRING,
									Value: "\"" + input.typeName + "\"",
								},
							},
							&ast.KeyValueExpr{
								Key: &ast.Ident{Name: "FieldName"},
								Value: &ast.BasicLit{
									Kind:  token.STRING,
									Value: "\"" + input.fieldName + "\"",
								},
							},
							&ast.KeyValueExpr{
								Key: &ast.Ident{Name: "Value"},
								Value: &ast.SelectorExpr{
									X:   &ast.Ident{Name: input.receiverVariableName},
									Sel: &ast.Ident{Name: input.fieldName},
								},
							},
						},
					},
				},
			}}}
	errorCheck := getErrorCheckExpr(fn, LibErrorVariableName)

	isInitializedCheck.Body.List = []ast.Stmt{&getFieldFromLib, &errorCheck}
	return isInitializedCheck
}

package parser

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strings"
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

func getNewCtorStmts(fn *ast.FuncDecl, typeName string, typesWithCustomId map[string]string) ([]ast.Stmt, error) {
	toInsertVariableName, err := getObjectVariableName(fn.Type.Params)
	if err != nil {
		return nil, err
	}

	idFieldName := ""
	if idField, isPresent := typesWithCustomId[typeName]; isPresent {
		idFieldName = idField
	} else {
		idFieldName = "Id"
	}

	insertInLib := ast.AssignStmt{
		Tok: token.DEFINE,
		Lhs: []ast.Expr{
			&ast.Ident{Name: "out"},
			&ast.Ident{Name: LibErrorVariableName},
		},
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "lib"},
					Sel: &ast.Ident{Name: "Insert"},
				},
				Args: []ast.Expr{
					&ast.Ident{Name: toInsertVariableName},
				},
			}}}
	errorCheck := getErrorCheckExpr(fn, LibErrorVariableName)
	idAssign := ast.AssignStmt{
		Tok: token.ASSIGN,
		Lhs: []ast.Expr{
			&ast.SelectorExpr{
				X:   &ast.Ident{Name: toInsertVariableName},
				Sel: &ast.Ident{Name: idFieldName},
			},
		},
		Rhs: []ast.Expr{
			&ast.Ident{Name: "out"},
		}}

	return []ast.Stmt{&insertInLib, &errorCheck, &idAssign}, nil
}

func getObjectVariableName(params *ast.FieldList) (string, error) {
	if params.List == nil || len(params.List) == 0 {
		return "", fmt.Errorf("object to be inserted not found in the parameters list")
	} else if len(params.List) > 1 {
		return "", fmt.Errorf("maximum allowed number of parameters is 1")
	}

	return params.List[0].Names[0].Name, nil
}

func getUpsertInLibExpr(fn *ast.FuncDecl, typesWithCustomId map[string]string) ast.AssignStmt {
	typeName := getFunctionReceiverTypeAsString(fn.Recv)
	idFieldName := ""
	if idField, isPresent := typesWithCustomId[typeName]; isPresent {
		idFieldName = idField
	} else {
		idFieldName = "Id"
	}

	return ast.AssignStmt{
		Tok: token.ASSIGN,
		Lhs: []ast.Expr{
			&ast.Ident{Name: LibErrorVariableName},
		},
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "lib"},
					Sel: &ast.Ident{Name: "Upsert"},
				},
				Args: []ast.Expr{
					&ast.Ident{Name: fn.Recv.List[0].Names[0].Name},
					&ast.SelectorExpr{
						X:   &ast.Ident{Name: fn.Recv.List[0].Names[0].Name},
						Sel: &ast.Ident{Name: idFieldName},
					}},
			},
		},
	}
}

func getPointerAssignStmt(receiverName string) ast.AssignStmt {
	return ast.AssignStmt{
		Tok: token.ASSIGN,
		Lhs: []ast.Expr{
			&ast.Ident{Name: receiverName},
		},
		Rhs: []ast.Expr{
			&ast.StarExpr{
				X: &ast.Ident{Name: TemporaryReceiverName},
			},
		},
	}
}

func getReadFromLibExpr(fn *ast.FuncDecl, typesWithCustomId map[string]string) (ast.AssignStmt, bool) {
	typeName := types.ExprString(fn.Recv.List[0].Type)
	isPointerReceiver := strings.Contains(typeName, "*")
	typeName = strings.TrimPrefix(typeName, "*")

	idFieldName := ""
	if idField, isPresent := typesWithCustomId[typeName]; isPresent {
		idFieldName = idField
	} else {
		idFieldName = "Id"
	}

	assignStmt := ast.AssignStmt{
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.IndexExpr{
					Index: &ast.Ident{Name: typeName},
					X: &ast.SelectorExpr{
						X:   &ast.Ident{Name: "lib"},
						Sel: &ast.Ident{Name: "Get"},
					},
				},
				Args: []ast.Expr{
					&ast.SelectorExpr{
						X:   &ast.Ident{Name: fn.Recv.List[0].Names[0].Name},
						Sel: &ast.Ident{Name: idFieldName},
					},
				},
			},
		},
	}

	if isPointerReceiver {
		assignStmt.Lhs = []ast.Expr{
			&ast.Ident{Name: fn.Recv.List[0].Names[0].Name},
			&ast.Ident{Name: LibErrorVariableName},
		}
	} else {
		assignStmt.Lhs = []ast.Expr{
			&ast.Ident{Name: TemporaryReceiverName},
			&ast.Ident{Name: LibErrorVariableName},
		}
	}

	return assignStmt, isPointerReceiver
}

package parser

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"
)

func getNobjectStateConditionalRetrieval(fn *ast.FuncDecl, parsedPackage ParsedPackage) ast.IfStmt {
	isInitializedCheck := getIsInitializedCheck(fn.Recv.List[0].Names[0].Name)
	readFromLibExpr, isPointerReceiver := getReadFromLibExpr(fn, parsedPackage.TypesWithCustomId)
	errorCheck := getErrorCheckExpr(fn, LibErrorVariableName)

	tempRecvAssignment := getTempRecvAssignStmt(fn.Recv.List[0].Names[0].Name, isPointerReceiver)
	initCall := getInitCall(fn.Recv.List[0].Names[0].Name)
	isInitializedCheck.Body.List = []ast.Stmt{&readFromLibExpr, &errorCheck, &tempRecvAssignment, &initCall}

	return isInitializedCheck
}

func getNobjectStateConditionalUpsert(fn *ast.FuncDecl, parsedPackage ParsedPackage) ast.IfStmt {
	isInitializedCheck := getIsInitializedCheck(fn.Recv.List[0].Names[0].Name)
	saveExpr := getUpsertInLibExpr(fn, parsedPackage.TypesWithCustomId)
	erorCheck := getErrorCheckExpr(fn, LibErrorVariableName)

	isInitializedCheck.Body.List = []ast.Stmt{&saveExpr, &erorCheck}
	return isInitializedCheck
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
						Sel: &ast.Ident{Name: LibraryGetObjectStateMethod},
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

	assignStmt.Lhs = []ast.Expr{
		&ast.Ident{Name: TemporaryReceiverName},
		&ast.Ident{Name: LibErrorVariableName},
	}

	return assignStmt, isPointerReceiver
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
		Tok: token.DEFINE,
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

func getTempRecvAssignStmt(receiverName string, isPointerReceiver bool) ast.AssignStmt {
	assign := ast.AssignStmt{
		Tok: token.ASSIGN,
		Lhs: []ast.Expr{
			&ast.Ident{Name: receiverName},
		},
	}

	if isPointerReceiver {
		assign.Rhs = []ast.Expr{&ast.Ident{Name: TemporaryReceiverName}}
	} else {
		assign.Rhs = []ast.Expr{
			&ast.StarExpr{
				X: &ast.Ident{Name: TemporaryReceiverName},
			},
		}
	}

	return assign
}

func getIsInitializedCheck(receiverVariableName string) ast.IfStmt {
	ifStmt := ast.IfStmt{
		Cond: &ast.SelectorExpr{
			X:   &ast.Ident{Name: receiverVariableName},
			Sel: &ast.Ident{Name: IsInitializedFieldName},
		},
		Body: &ast.BlockStmt{},
	}

	return ifStmt
}

func getInitCall(receiverVariableName string) ast.ExprStmt {
	return ast.ExprStmt{X: &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   &ast.Ident{Name: receiverVariableName},
			Sel: &ast.Ident{Name: InitFunctionName},
		},
		Args: []ast.Expr{},
	}}
}

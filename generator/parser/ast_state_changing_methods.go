package parser

import (
	"errors"
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

type getDBStmtsParam struct {
	idFieldName          string
	typeName             string
	receiverVariableName string
	fieldName            string
	fieldType            string
}

func getGetterDBStmts(fn *ast.FuncDecl, input getDBStmtsParam) ast.IfStmt {
	isInitializedCheck := getIsInitializeCheck(fn.Recv.List[0].Names[0].Name)

	var fieldValPrep *ast.AssignStmt
	// the second condition makes sures it's not a map with values of type list
	if strings.Contains(input.fieldType, "[]") && !strings.HasPrefix(input.fieldType, "map") {
		fieldValPrep = &ast.AssignStmt{
			Tok: token.DEFINE,
			Lhs: []ast.Expr{&ast.Ident{Name: "fieldValue"}},
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun:  &ast.Ident{Name: "make"},
					Args: []ast.Expr{&ast.ArrayType{Elt: &ast.Ident{Name: strings.TrimPrefix(input.fieldType, "[]")}}},
				},
			},
		}
	} else if strings.Contains(input.fieldType, "map") {
		key, err := getKeyTypeOfMap(input.fieldType)
		if err != nil {
			fmt.Print(err)
		}
		value, err := getValueTypeOfMap(input.fieldType, key)
		if err != nil {
			fmt.Print(err)
		}

		fieldValPrep = &ast.AssignStmt{
			Tok: token.DEFINE,
			Lhs: []ast.Expr{&ast.Ident{Name: "fieldValue"}},
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.Ident{Name: "make"},
					Args: []ast.Expr{
						&ast.MapType{Key: &ast.Ident{Name: key}, Value: value},
					},
				},
			},
		}
	} else {
		fieldValPrep = &ast.AssignStmt{
			Tok: token.DEFINE,
			Lhs: []ast.Expr{&ast.Ident{Name: "fieldValue"}},
			Rhs: []ast.Expr{
				&ast.StarExpr{X: &ast.CallExpr{
					Args: []ast.Expr{&ast.Ident{Name: input.fieldType}},
					Fun:  &ast.Ident{Name: "new"},
				}},
			},
		}
	}

	getFieldFromLib := ast.AssignStmt{
		Tok: token.DEFINE,
		Lhs: []ast.Expr{
			&ast.Ident{Name: LibErrorVariableName},
		},
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "lib"},
					Sel: &ast.Ident{Name: LibraryGetFieldOfType},
				},
				Args: []ast.Expr{
					&ast.CompositeLit{
						Type: &ast.SelectorExpr{
							X:   &ast.Ident{Name: "lib"},
							Sel: &ast.Ident{Name: GetStateParamType},
						},
						Elts: []ast.Expr{
							&ast.KeyValueExpr{
								Key: &ast.Ident{Name: Id},
								Value: &ast.SelectorExpr{
									X:   &ast.Ident{Name: input.receiverVariableName},
									Sel: &ast.Ident{Name: input.idFieldName},
								},
							},
							&ast.KeyValueExpr{
								Key: &ast.Ident{Name: TypeName},
								Value: &ast.BasicLit{
									Kind:  token.STRING,
									Value: "\"" + input.typeName + "\"",
								},
							},
							&ast.KeyValueExpr{
								Key: &ast.Ident{Name: FieldName},
								Value: &ast.BasicLit{
									Kind:  token.STRING,
									Value: "\"" + input.fieldName + "\"",
								},
							},
						},
					},
					&ast.UnaryExpr{
						Op: token.AND,
						X:  &ast.Ident{Name: "fieldValue"},
					},
				},
			},
		}}
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
			&ast.Ident{Name: "fieldValue"},
		},
	}

	isInitializedCheck.Body.List = []ast.Stmt{fieldValPrep, &getFieldFromLib, &errorCheck, &fieldAssign}
	return isInitializedCheck
}

func getValueTypeOfMap(mapType string, mapKey string) (ast.Expr, error) {
	splitted := strings.Split(mapType, "["+mapKey+"]")

	if len(splitted) > 0 {
		mapValue := splitted[1]

		if strings.Contains(mapValue, "[]") {

			return &ast.ArrayType{
				Elt: &ast.Ident{Name: mapValue[2:]},
			}, nil
		} else if strings.Contains(mapValue, "map") {

			key, err1 := getKeyTypeOfMap(mapValue)
			if err1 != nil {
				return nil, err1
			}
			value, err2 := getValueTypeOfMap(mapValue, key)
			if err2 != nil {
				return nil, err2
			}

			return &ast.MapType{
				Key:   &ast.Ident{Name: key},
				Value: value,
			}, nil
		}
	}

	return nil, errors.New("error occurred while parsing the field " + mapType)
}

func getKeyTypeOfMap(mapType string) (string, error) {
	res := strings.Split(mapType, "[")
	if len(res) > 0 {
		keyType := strings.Split(res[1], "]")
		if keyType != nil {
			return keyType[0], nil
		}
		return "", errors.New("error occurred while parsing the field " + mapType)
	}
	return "", errors.New("error occurred while parsing the field " + mapType)
}

func getReferenceNavigationListDBStmts(fn *ast.FuncDecl) ast.IfStmt {

	isUnInitializedCheck := ast.IfStmt{
		Cond: &ast.UnaryExpr{
			Op: token.NOT,
			X: &ast.SelectorExpr{
				X:   &ast.Ident{Name: fn.Recv.List[0].Names[0].Name},
				Sel: &ast.Ident{Name: IsInitializedFieldName},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{&ast.ReturnStmt{
				Results: []ast.Expr{
					&ast.Ident{Name: "nil"},
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   &ast.Ident{Name: "fmt"},
							Sel: &ast.Ident{Name: "Errorf"},
						},
						Args: []ast.Expr{&ast.BasicLit{
							Kind:  token.STRING,
							Value: "`fields of type ReferenceNavigationList can be used only after instance initialization. \n\t\t\tUse lib.Load or lib.Export from the Nubes library to create initialized instances`",
						}},
					},
				},
			}},
		},
	}
	return isUnInitializedCheck
}

func getSetterDBStmts(fn *ast.FuncDecl, input getDBStmtsParam) ast.IfStmt {
	isInitializedCheck := getIsInitializeCheck(fn.Recv.List[0].Names[0].Name)
	getFieldFromLib := ast.AssignStmt{
		Tok: token.DEFINE,
		Lhs: []ast.Expr{
			&ast.Ident{Name: LibErrorVariableName},
		},
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "lib"},
					Sel: &ast.Ident{Name: SetField},
				},
				Args: []ast.Expr{
					&ast.CompositeLit{
						Type: &ast.SelectorExpr{
							X:   &ast.Ident{Name: "lib"},
							Sel: &ast.Ident{Name: SetFieldParam},
						},
						Elts: []ast.Expr{
							&ast.KeyValueExpr{
								Key: &ast.Ident{Name: Id},
								Value: &ast.SelectorExpr{
									X:   &ast.Ident{Name: input.receiverVariableName},
									Sel: &ast.Ident{Name: input.idFieldName},
								},
							},
							&ast.KeyValueExpr{
								Key: &ast.Ident{Name: TypeName},
								Value: &ast.BasicLit{
									Kind:  token.STRING,
									Value: "\"" + input.typeName + "\"",
								},
							},
							&ast.KeyValueExpr{
								Key: &ast.Ident{Name: FieldName},
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

func getNobjectFunctionProlog(fn *ast.FuncDecl, parsedPackage ParsedPackage) []ast.Stmt {
	invocationDepthInc := getInvocationDepthInceremntStmt(fn.Recv.List[0].Names[0].Name)
	isInitializedCheck := getIsInitializedAndInvocationDepthEqOneCheck(fn.Recv.List[0].Names[0].Name)
	readFromLibExpr := getReadFromLibExpr(fn, parsedPackage.TypesWithCustomId)
	errorCheck := getErrorCheckExpr(fn, LibErrorVariableName)

	isInitializedCheck.Body.List = []ast.Stmt{&readFromLibExpr, &errorCheck}
	return []ast.Stmt{&invocationDepthInc, &isInitializedCheck}
}

func getNobjectStateConditionalUpsert(typeName, receiverVarName string, parsedPackage ParsedPackage) ast.IfStmt {
	isInitializedCheck := getIsInitializedAndInvocationDepthEqOneCheck(receiverVarName)
	saveExpr := getUpsertInLibExpr(typeName, receiverVarName, parsedPackage.TypesWithCustomId)
	erorCheck := ast.IfStmt{
		Cond: &ast.BinaryExpr{
			X:  &ast.Ident{Name: LibErrorVariableName},
			Op: token.NEQ,
			Y:  &ast.Ident{Name: "nil"},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.Ident{
							Name: LibErrorVariableName,
						},
					},
				},
			},
		},
	}

	isInitializedCheck.Body.List = []ast.Stmt{&saveExpr, &erorCheck}
	return isInitializedCheck
}

func getSaveChangesMethodForType(typeName string, parsedPackage ParsedPackage) *ast.FuncDecl {
	receiverVarName := "receiver"
	ifStmt := getNobjectStateConditionalUpsert(typeName, receiverVarName, parsedPackage)

	function := &ast.FuncDecl{
		Name: &ast.Ident{Name: SaveChangesIfInitialized},
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{{Name: receiverVarName}},
					Type:  &ast.StarExpr{X: &ast.Ident{Name: typeName}},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ifStmt,
				&ast.ReturnStmt{Results: []ast.Expr{&ast.Ident{Name: "nil"}}},
			},
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{Type: &ast.Ident{Name: "error"}},
				},
			},
		},
	}

	return function
}

func invokeSaveChangesMethodForType(fn *ast.FuncDecl, parsedPackage ParsedPackage) *ast.AssignStmt {
	return &ast.AssignStmt{
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: fn.Recv.List[0].Names[0].Name},
					Sel: &ast.Ident{Name: SaveChangesIfInitialized},
				},
			},
		},
		Lhs: []ast.Expr{
			&ast.Ident{Name: UpsertLibErrorVariableName},
		},
	}
}

func getReadFromLibExpr(fn *ast.FuncDecl, typesWithCustomId map[string]string) ast.AssignStmt {
	typeName := types.ExprString(fn.Recv.List[0].Type)
	isPointerReceiver := strings.Contains(typeName, "*")
	typeName = strings.TrimPrefix(typeName, "*")

	idFieldName := ""
	if idField, isPresent := typesWithCustomId[typeName]; isPresent {
		idFieldName = idField
	} else {
		idFieldName = Id
	}

	assignStmt := ast.AssignStmt{
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "lib"},
					Sel: &ast.Ident{Name: LibraryGetObjectStateMethod},
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
		(assignStmt.Rhs[0].(*ast.CallExpr)).Args = append((assignStmt.Rhs[0].(*ast.CallExpr)).Args, ast.NewIdent(fn.Recv.List[0].Names[0].Name))
	} else {
		(assignStmt.Rhs[0].(*ast.CallExpr)).Args = append((assignStmt.Rhs[0].(*ast.CallExpr)).Args, &ast.UnaryExpr{Op: token.AND, X: ast.NewIdent(fn.Recv.List[0].Names[0].Name)})
	}

	assignStmt.Lhs = []ast.Expr{
		&ast.Ident{Name: LibErrorVariableName},
	}

	return assignStmt
}

func getUpsertInLibExpr(typeName, receiverVariableName string, typesWithCustomId map[string]string) ast.AssignStmt {
	idFieldName := ""
	if idField, isPresent := typesWithCustomId[typeName]; isPresent {
		idFieldName = idField
	} else {
		idFieldName = Id
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
					Sel: &ast.Ident{Name: Upsert},
				},
				Args: []ast.Expr{
					&ast.Ident{Name: receiverVariableName},
					&ast.SelectorExpr{
						X:   &ast.Ident{Name: receiverVariableName},
						Sel: &ast.Ident{Name: idFieldName},
					}},
			},
		},
	}
}

func getInvocationDepthInceremntStmt(receiverVariableName string) ast.IncDecStmt {
	return ast.IncDecStmt{
		Tok: token.INC,
		X: &ast.SelectorExpr{
			Sel: &ast.Ident{Name: InvocationDepthFieldName},
			X:   &ast.Ident{Name: receiverVariableName},
		},
	}
}

func getInvocationDepthDecremntStmt(receiverVariableName string) *ast.IncDecStmt {
	return &ast.IncDecStmt{
		Tok: token.DEC,
		X: &ast.SelectorExpr{
			Sel: &ast.Ident{Name: InvocationDepthFieldName},
			X:   &ast.Ident{Name: receiverVariableName},
		},
	}
}

func getIsInitializedAndInvocationDepthEqOneCheck(receiverVariableName string) ast.IfStmt {
	ifStmt := ast.IfStmt{
		Cond: &ast.BinaryExpr{
			Op: token.LAND,
			X: &ast.SelectorExpr{
				X:   &ast.Ident{Name: receiverVariableName},
				Sel: &ast.Ident{Name: IsInitializedFieldName},
			},
			Y: &ast.BinaryExpr{
				Op: token.EQL,
				X: &ast.SelectorExpr{
					X:   &ast.Ident{Name: receiverVariableName},
					Sel: &ast.Ident{Name: InvocationDepthFieldName},
				},
				Y: &ast.BasicLit{
					Kind:  token.INT,
					Value: "1",
				},
			},
		},

		Body: &ast.BlockStmt{},
	}

	return ifStmt
}

func getIsInitializeCheck(receiverVariableName string) ast.IfStmt {
	ifStmt := ast.IfStmt{
		Cond: &ast.SelectorExpr{
			X:   &ast.Ident{Name: receiverVariableName},
			Sel: &ast.Ident{Name: IsInitializedFieldName},
		},
		Body: &ast.BlockStmt{},
	}

	return ifStmt
}

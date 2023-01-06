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
)

func AddReadWriteOpToMethods(path string, parsedPackage ParsedPackage) {
	set := token.NewFileSet()
	packs, err := parser.ParseDir(set, path, nil, 0)
	assertDirParsed(err)

	for _, pack := range packs {
		for filePath, f := range pack.Files {
			fileModified := false
			for _, d := range f.Decls {
				if fn, isFn := d.(*ast.FuncDecl); isFn {

					if fn.Recv != nil && fn.Name.Name != NobjectImplementationMethod && f.Name.Name != CustomIdImplementationMethod {
						typeName := getFunctionReceiverTypeAsString(fn.Recv)

						if isNobject := parsedPackage.IsNobjectInOrginalPackage[typeName]; isNobject && !isFunctionStateless(fn.Recv) {
							if retParamsVerifier.Check(fn) && !isReadOperationAlreadyAdded(fn, set) {

								fileModified = true

								SaveExpr := getSaveChangesInLibExpr(fn, parsedPackage.TypesWithCustomId)
								ReadFromLibExpr, isPointerReceiver := getReadFromLibExpr(fn, parsedPackage.TypesWithCustomId)
								ErrorCheck := getErrorCheckExpr(fn, LibErrorVariableName)

								fn.Body.List = prepend[ast.Stmt](fn.Body.List, &ErrorCheck)
								if !isPointerReceiver {
									pointerStms := getPointerAssignStmt(fn.Recv.List[0].Names[0].Name)
									fn.Body.List = prepend[ast.Stmt](fn.Body.List, &pointerStms)
								}
								fn.Body.List = prepend[ast.Stmt](fn.Body.List, &ReadFromLibExpr)
								fn.Body.List = prependBeforeLastElem[ast.Stmt](fn.Body.List, &SaveExpr)
								fn.Body.List = prependBeforeLastElem[ast.Stmt](fn.Body.List, &ErrorCheck)
							}
						}
					}
				}
			}

			if fileModified {
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
				err := printer.Fprint(&buf, set, f)
				if err != nil {
					fmt.Println(err)
				}
				nobjectTypeFile, err := os.Create(filePath)
				if err != nil {
					fmt.Println(err)
				}
				buf.WriteTo(nobjectTypeFile)
				nobjectTypeFile.Close()
			}
		}
	}
}

func isReadOperationAlreadyAdded(fn *ast.FuncDecl, set *token.FileSet) bool {
	if len(fn.Body.List) > 2 {
		assign, _ := fn.Body.List[0].(*ast.AssignStmt)
		secLastElem, _ := getFunctionBodyStmtAsString(set, assign)
		return strings.Contains(secLastElem, "lib.Get")
	}
	return false
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

func getSaveChangesInLibExpr(fn *ast.FuncDecl, typesWithCustomId map[string]string) ast.AssignStmt {
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

func prependBeforeLastElem[T any](stmtList []T, toInsert T) []T {
	x := append(stmtList, *new(T))
	x[len(x)-1] = x[len(x)-2]
	x[len(x)-2] = toInsert
	return x
}

func prepend[T any](list []T, toPrepend T) []T {
	x := append(list, *new(T))
	copy(x[1:], x)
	x[0] = toPrepend
	return x
}

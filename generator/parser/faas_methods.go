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

func AddDBOperationsToMethods(path string, parsedPackage ParsedPackage) {
	set := token.NewFileSet()
	packs, err := parser.ParseDir(set, path, nil, 0)
	assertDirParsed(err)

	IsTypeNewCtorImplemented := make(map[string]bool)
	IsTypeReNewCtorImplemented := make(map[string]bool)

	for _, pack := range packs {
		for filePath, f := range pack.Files {
			fileModified := false
			for _, d := range f.Decls {
				if fn, isFn := d.(*ast.FuncDecl); isFn {

					if fn.Recv == nil {

						if strings.HasPrefix(fn.Name.Name, ConstructorPrefix) {
							typeName := strings.TrimPrefix(fn.Name.Name, ConstructorPrefix)
							if parsedPackage.IsNobjectInOrginalPackage[typeName] {
								if !areDBOperationsAlreadyAddedToNewCtor(fn.Body, set) {
									idFieldName := getIdFieldNameOfType(typeName, parsedPackage.TypesWithCustomId)
									stmtsToInsert, err := getNewCtorStmts(fn, typeName, idFieldName)
									if err != nil {
										fmt.Println("wrong constructor definition of ", fn.Name.Name, ": ", err)
										continue
									}
									fileModified = true
									fn.Body.List = appendListBeforeLastElem(fn.Body.List, stmtsToInsert)
								}

								IsTypeNewCtorImplemented[typeName] = true
							}
						}
						if strings.HasPrefix(fn.Name.Name, ReConstructorPrefix) {
							IsTypeReNewCtorImplemented[strings.TrimPrefix(fn.Name.Name, ConstructorPrefix)] = true
						}

					} else if fn.Name.Name != NobjectImplementationMethod && f.Name.Name != CustomIdImplementationMethod {
						typeName := getFunctionReceiverTypeAsString(fn.Recv)

						if isNobject := parsedPackage.IsNobjectInOrginalPackage[typeName]; isNobject {

							if strings.HasPrefix(fn.Name.Name, GetPrefix) {

								fieldName := strings.TrimPrefix(fn.Name.Name, GetPrefix)
								if _, fieldExist := parsedPackage.TypeFields[typeName][fieldName]; fieldExist {
									if !isGetFieldStmtAlreadyAdded(fn.Body, set) {
										idFieldName := getIdFieldNameOfType(typeName, parsedPackage.TypesWithCustomId)
										stmtsToInsert := getGetterDBStmts(fn, getDBStmtsParam{
											idFieldName:          idFieldName,
											typeName:             typeName,
											fieldName:            fieldName,
											fieldType:            getFirstFunctionReturnTypeAsString(fn),
											receiverVariableName: fn.Recv.List[0].Names[0].Name,
										})
										fn.Body.List = prependList(fn.Body.List, stmtsToInsert)
										fileModified = true
									}
									continue
								}

							} else if strings.HasPrefix(fn.Name.Name, SetPrefix) {
								fieldName := strings.TrimPrefix(fn.Name.Name, SetPrefix)
								if _, fieldExists := parsedPackage.TypeFields[typeName][fieldName]; fieldExists {

									if !isSetFieldStmtAlreadyAdded(fn.Body, set) {
										idFieldName := getIdFieldNameOfType(typeName, parsedPackage.TypesWithCustomId)
										stmtsToInsert := getSetterDBStmts(fn, getDBStmtsParam{
											idFieldName:          idFieldName,
											typeName:             typeName,
											fieldName:            fieldName,
											fieldType:            getFirstFunctionReturnTypeAsString(fn),
											receiverVariableName: fn.Recv.List[0].Names[0].Name,
										})
										fn.Body.List = appendListBeforeLastElem(fn.Body.List, stmtsToInsert)
										fileModified = true
									}
									continue
								}
							}

							if !isFunctionStateless(fn.Recv) && retParamsVerifier.Check(fn) && !isDBGetOperationAlreadyAddedToMethod(fn.Body, set) {
								fileModified = true

								SaveExpr := getUpsertInLibExpr(fn, parsedPackage.TypesWithCustomId)
								ReadFromLibExpr, isPointerReceiver := getReadFromLibExpr(fn, parsedPackage.TypesWithCustomId)
								ErrorCheck := getErrorCheckExpr(fn, LibErrorVariableName)

								if !isPointerReceiver {
									pointerStms := getPointerAssignStmt(fn.Recv.List[0].Names[0].Name)
									fn.Body.List = prepend[ast.Stmt](fn.Body.List, &pointerStms)
								}
								fn.Body.List = prepend[ast.Stmt](fn.Body.List, &ErrorCheck)
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

	printWarningIfCtorMissing(IsTypeNewCtorImplemented, IsTypeReNewCtorImplemented, parsedPackage.IsNobjectInOrginalPackage)
}

func printWarningIfCtorMissing(isTypeNewCtorImpl, isTypeReNewCtorImpl, isNobjectInOrgPkg map[string]bool) {
	for typeName, isNobject := range isNobjectInOrgPkg {
		if isNobject {
			if !isTypeNewCtorImpl[typeName] {
				fmt.Println("missing constructor definition for new instances of object type ", typeName)
			}

			if !isTypeReNewCtorImpl[typeName] {
				fmt.Println("missing constructor definition for existing instances of object type ", typeName)
			}
		}
	}
}

func areDBOperationsAlreadyAddedToNewCtor(funcBlock *ast.BlockStmt, set *token.FileSet) bool {
	if funcBlock != nil && funcBlock.List != nil && len(funcBlock.List) > 3 {
		assign, _ := funcBlock.List[len(funcBlock.List)-4].(*ast.AssignStmt)
		secLastElem, _ := getFunctionBodyStmtAsString(set, assign)
		return strings.Contains(secLastElem, "lib.Insert")
	}
	return false
}

func isDBGetOperationAlreadyAddedToMethod(funcBlock *ast.BlockStmt, set *token.FileSet) bool {
	if funcBlock != nil && funcBlock.List != nil && len(funcBlock.List) > 2 {
		assign, _ := funcBlock.List[0].(*ast.AssignStmt)
		secLastElem, _ := getFunctionBodyStmtAsString(set, assign)
		return strings.Contains(secLastElem, "lib.Get")
	}
	return false
}

func isGetFieldStmtAlreadyAdded(blockStmt *ast.BlockStmt, set *token.FileSet) bool {
	if blockStmt != nil && blockStmt.List != nil && len(blockStmt.List) > 0 {
		assign, _ := blockStmt.List[0].(*ast.AssignStmt)
		firstStmtString, _ := getFunctionBodyStmtAsString(set, assign)
		return strings.Contains(firstStmtString, "lib.GetField")
	}
	return false
}

func isSetFieldStmtAlreadyAdded(blockStmt *ast.BlockStmt, set *token.FileSet) bool {
	if blockStmt != nil && blockStmt.List != nil && len(blockStmt.List) > 2 {
		assign, _ := blockStmt.List[len(blockStmt.List)-3].(*ast.AssignStmt)
		stmtString, _ := getFunctionBodyStmtAsString(set, assign)
		return strings.Contains(stmtString, "lib.SetField")
	}
	return false
}

func getIdFieldNameOfType(typeName string, typesWithCustomId map[string]string) string {
	if idField, isPresent := typesWithCustomId[typeName]; isPresent {
		return idField
	}

	return "Id"
}

func prependBeforeLastElem[T any](stmtList []T, toInsert T) []T {
	x := append(stmtList, *new(T))
	x[len(x)-1] = x[len(x)-2]
	x[len(x)-2] = toInsert
	return x
}

func appendListBeforeLastElem[T any](stmtList []T, toInsert []T) []T {
	x := make([]T, len(stmtList)+len(toInsert))
	stmtNum := len(stmtList)
	copy(x[:], stmtList[0:stmtNum-1])
	copy(x[stmtNum-1:], toInsert[:])
	x[len(x)-1] = stmtList[stmtNum-1]
	return x
}

func prepend[T any](list []T, toPrepend T) []T {
	x := append(list, *new(T))
	copy(x[1:], x)
	x[0] = toPrepend
	return x
}

func prependList[T any](list []T, toPrepend []T) []T {
	toPrependNum := len(toPrepend)
	x := make([]T, len(list)+toPrependNum)
	copy(x[:], toPrepend[:])
	copy(x[toPrependNum:], list[:])
	return x
}

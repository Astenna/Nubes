package parser

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

func (t TypeSpecParser) adjustMethods() {

	for path, detectedFunctionsList := range t.detectedFunctions {
		for _, detectedFunction := range detectedFunctionsList {

			fn := detectedFunction.Function

			if fn.Recv == nil {
				continue
			}

			typeName := getFunctionReceiverTypeAsString(fn.Recv)
			if isNobject := t.Output.IsNobjectInOrginalPackage[typeName]; isNobject {
				isGetter := t.addDBOperationsIfGetter(fn, path)
				if !isGetter {
					isSetter := t.addDBOperationsIfSetter(fn, path)
					if !isSetter {
						if !isFunctionStateless(fn.Recv) && retParamsVerifier.Check(fn) {
							// add retrieve from DB if it doesn't exist
							if !isInitFieldCheckAlreadyAddedAsFirstStmt(fn.Body) {
								t.fileChanged[path] = true
								retrieveStateIfInitialized := getNobjectStateConditionalRetrieval(fn, t.Output)
								fn.Body.List = prepend[ast.Stmt](fn.Body.List, &retrieveStateIfInitialized)
							}

							// add invocation of a function that save changes in DB
							// only if the error returned is nil
							astutil.Apply(fn, nil, func(c *astutil.Cursor) bool {
								n := c.Node()

								if x, ok := n.(*ast.ReturnStmt); ok {

									if isErrorToBeReturnedNil(*x) {
										c.InsertBefore(invokeSaveChangesMethodForType(fn, t.Output))
										if len(x.Results) > 1 {
											ident := x.Results[1].(*ast.Ident)
											ident.Name = UpsertLibErrorVariableName
										} else {
											ident := x.Results[0].(*ast.Ident)
											ident.Name = UpsertLibErrorVariableName
										}
									}
									t.fileChanged[path] = true
								}
								return true
							})
						}
					}
				}
			}
		}
	}
}

func (t TypeSpecParser) addDBOperationsIfSetter(fn *ast.FuncDecl, path string) bool {
	typeName := getFunctionReceiverTypeAsString(fn.Recv)
	if strings.HasPrefix(fn.Name.Name, SetPrefix) {
		fieldName := strings.TrimPrefix(fn.Name.Name, SetPrefix)
		if fieldType, fieldExists := t.Output.TypeFields[typeName][fieldName]; fieldExists {
			if !isInitFieldCheckAlreadyAddedAsSecondLastStmt(fn.Body) {
				idFieldName := getIdFieldNameOfType(typeName, t.Output.TypesWithCustomId)
				if strings.Contains(fieldType, LibraryReferenceNavigationList) {
					returnErrorIfNotInitialized := getReferenceNavigationListDBStmts(fn, getDBStmtsParam{
						idFieldName:          idFieldName,
						typeName:             typeName,
						fieldName:            fieldName,
						fieldType:            fieldType,
						receiverVariableName: fn.Recv.List[0].Names[0].Name,
					})
					fn.Body.List = prepend[ast.Stmt](fn.Body.List, &returnErrorIfNotInitialized)
					t.fileChanged[path] = true
				} else {
					saveInDbIfInitialized := getSetterDBStmts(fn, getDBStmtsParam{
						idFieldName:          idFieldName,
						typeName:             typeName,
						fieldName:            fieldName,
						fieldType:            fieldType,
						receiverVariableName: fn.Recv.List[0].Names[0].Name,
					})
					fn.Body.List = appendBeforeLastElem[ast.Stmt](fn.Body.List, &saveInDbIfInitialized)
				}

				t.fileChanged[path] = true
			}
			return true
		}
	}

	return false
}

func (t TypeSpecParser) addDBOperationsIfGetter(fn *ast.FuncDecl, path string) bool {
	typeName := getFunctionReceiverTypeAsString(fn.Recv)

	if strings.HasPrefix(fn.Name.Name, GetPrefix) {
		fieldName := strings.TrimPrefix(fn.Name.Name, GetPrefix)
		fieldName = strings.TrimSuffix(fieldName, "Ids")

		if fieldType, fieldExist := t.Output.TypeFields[typeName][fieldName]; fieldExist {

			idFieldName := getIdFieldNameOfType(typeName, t.Output.TypesWithCustomId)
			if !isInitFieldCheckAlreadyAddedAsFirstStmt(fn.Body) {
				if strings.Contains(fieldType, LibraryReferenceNavigationList) {
					returnErrorIfNotInitialized := getReferenceNavigationListDBStmts(fn, getDBStmtsParam{
						idFieldName:          idFieldName,
						typeName:             typeName,
						fieldName:            fieldName,
						fieldType:            fieldType,
						receiverVariableName: fn.Recv.List[0].Names[0].Name,
					})
					fn.Body.List = prepend[ast.Stmt](fn.Body.List, &returnErrorIfNotInitialized)
					t.fileChanged[path] = true
				} else {
					retrieveFromDbIfInitialized := getGetterDBStmts(fn, getDBStmtsParam{
						idFieldName:          idFieldName,
						typeName:             typeName,
						fieldName:            fieldName,
						fieldType:            fieldType,
						receiverVariableName: fn.Recv.List[0].Names[0].Name,
					})
					fn.Body.List = prepend[ast.Stmt](fn.Body.List, &retrieveFromDbIfInitialized)
				}

				t.fileChanged[path] = true
			}
			return true
		}
	}

	return false
}

func isErrorToBeReturnedNil(x ast.ReturnStmt) bool {
	return (len(x.Results) > 1 && types.ExprString(x.Results[1]) == "nil") || (len(x.Results) == 1 && types.ExprString(x.Results[0]) == "nil")
}

func isInitFieldCheckAlreadyAddedAsFirstStmt(funcBlock *ast.BlockStmt) bool {
	if funcBlock != nil && funcBlock.List != nil && len(funcBlock.List) > 0 {
		ifStmt, _ := funcBlock.List[0].(*ast.IfStmt)
		if ifStmt != nil {
			ifConditionAsString := types.ExprString(ifStmt.Cond)
			return strings.Contains(ifConditionAsString, IsInitializedFieldName)
		}
	}
	return false
}

func isInitFieldCheckAlreadyAddedAsSecondLastStmt(funcBlock *ast.BlockStmt) bool {
	if funcBlock != nil && funcBlock.List != nil && len(funcBlock.List) > 1 {
		ifStmt, _ := funcBlock.List[len(funcBlock.List)-2].(*ast.IfStmt)
		if ifStmt != nil {
			ifConditionAsString := types.ExprString(ifStmt.Cond)
			return strings.Contains(ifConditionAsString, IsInitializedFieldName)
		}
	}
	return false
}

func getIdFieldNameOfType(typeName string, typesWithCustomId map[string]string) string {
	if idField, isPresent := typesWithCustomId[typeName]; isPresent {
		return idField
	}

	return "Id"
}

func appendBeforeLastElem[T any](stmtList []T, toInsert T) []T {
	x := append(stmtList, *new(T))
	x[len(x)-1] = x[len(x)-2]
	x[len(x)-2] = toInsert
	return x
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

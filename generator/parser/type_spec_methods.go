package parser

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

func (t TypeSpecParser) modifyAstMethods() {

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
						if !isFunctionStateless(fn.Recv) {
							// add retrieve from DB if it doesn't exist
							if !isInvocationDepthIncrementedInFirstStmt(fn.Body) {
								functionProlog := getNobjectFunctionProlog(fn, t.Output)
								fn.Body.List = prependList(fn.Body.List, functionProlog)
								t.fileChanged[path] = true
							}

							astutil.Apply(fn, nil, func(c *astutil.Cursor) bool {
								n := c.Node()

								if x, ok := n.(*ast.ReturnStmt); ok {

									// add invocation of a function that save changes in DB
									// (if function uses a pointer receiver - i.e. is not readonly)
									// decrement invocationDepth
									if !isFuncReadonly(fn.Recv) && isErrorToBeReturnedNil(*x) {
										c.InsertBefore(invokeSaveChangesMethodForType(fn, t.Output))
										if len(x.Results) > 1 {
											ident := x.Results[1].(*ast.Ident)
											ident.Name = UpsertLibErrorVariableName
										} else {
											ident := x.Results[0].(*ast.Ident)
											ident.Name = UpsertLibErrorVariableName
										}
										t.fileChanged[path] = true
									}

									// insert invocation depth decrement
									if !isInvocationDepthDecreasedBefore(c) {
										c.InsertBefore(getInvocationDepthDecremntStmt(fn.Recv.List[0].Names[0].Name))
										t.fileChanged[path] = true
									}
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

func isFuncReadonly(fields *ast.FieldList) bool {
	if fields.List == nil || fields.List[0].Names == nil {
		return true
	}
	_, ok := fields.List[0].Type.(*ast.StarExpr)
	return !ok
}

func isFunctionStateless(fields *ast.FieldList) bool {
	return fields.List == nil || fields.List[0].Names == nil || fields.List[0].Names[0].Name == ""
}

func isInvocationDepthDecreasedBefore(c *astutil.Cursor) bool {
	stmts := c.Parent().(*ast.BlockStmt)
	if index := c.Index(); index > 0 {
		decStmt, casted := stmts.List[index-1].(*ast.IncDecStmt)
		if decStmt == nil || !casted {
			return false
		}

		return decStmt.Tok == token.DEC && strings.Contains(types.ExprString(decStmt.X), InvocationDepthFieldName)
	}

	return false
}

func (t TypeSpecParser) addDBOperationsIfSetter(fn *ast.FuncDecl, path string) bool {
	typeName := getFunctionReceiverTypeAsString(fn.Recv)
	if strings.HasPrefix(fn.Name.Name, SetPrefix) {
		fieldName := strings.TrimPrefix(fn.Name.Name, SetPrefix)
		if fieldType, fieldExists := t.Output.TypeFields[typeName][fieldName]; fieldExists {
			if !isInitFieldCheckAlreadyAddedAsSecondLastStmt(fn.Body) {
				idFieldName := getIdFieldNameOfType(typeName, t.Output.TypesWithCustomId)
				if strings.Contains(fieldType, LibraryReferenceNavigationList) {
					returnErrorIfNotInitialized := getReferenceNavigationListDBStmts(fn)
					fn.Body.List = prependElem[ast.Stmt](fn.Body.List, &returnErrorIfNotInitialized)
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
					returnErrorIfNotInitialized := getReferenceNavigationListDBStmts(fn)
					fn.Body.List = prependElem[ast.Stmt](fn.Body.List, &returnErrorIfNotInitialized)
					t.fileChanged[path] = true
				} else {
					retrieveFromDbIfInitialized := getGetterDBStmts(fn, getDBStmtsParam{
						idFieldName:          idFieldName,
						typeName:             typeName,
						fieldName:            fieldName,
						fieldType:            fieldType,
						receiverVariableName: fn.Recv.List[0].Names[0].Name,
					})
					fn.Body.List = prependElem[ast.Stmt](fn.Body.List, &retrieveFromDbIfInitialized)
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

func isInvocationDepthIncrementedInFirstStmt(funcBlock *ast.BlockStmt) bool {
	if funcBlock != nil && funcBlock.List != nil && len(funcBlock.List) > 0 {
		incrementStmt, _ := funcBlock.List[0].(*ast.IncDecStmt)
		if incrementStmt != nil && incrementStmt.X != nil {
			incrementedVariable := types.ExprString(incrementStmt.X)
			return strings.Contains(incrementedVariable, InvocationDepthFieldName)
		}
	}
	return false
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

func prependElem[T any](list []T, toPrepend T) []T {
	x := append(list, *new(T))
	copy(x[1:], x)
	x[0] = toPrepend
	return x
}

func prependList[T any](list []T, toPrepend []T) []T {
	x := append(list, toPrepend...)
	copy(x[len(toPrepend):], x)
	copy(x[:len(toPrepend)], toPrepend)
	return x
}

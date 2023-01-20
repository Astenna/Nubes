package parser

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

func (t TypeSpecParser) detectAndAdjustMethods(isTypeNewCtorImplemented map[string]bool, isTypeReNewCtorImplemented map[string]bool, isTypeDestructorImplemented map[string]bool) {
	for _, pack := range t.packs {
		for path, f := range pack.Files {
			fileModified := false
			for _, d := range f.Decls {
				if fn, isFn := d.(*ast.FuncDecl); isFn {

					t.detectedFunctions[path] = append(t.detectedFunctions[path], detectedFunction{
						Function: fn,
						Imports:  f.Imports,
					})

					if fn.Recv == nil {

						ctorDetected := addDBOperationsIfCtor(fn, t.Output, t.tokenSet, isTypeNewCtorImplemented, isTypeReNewCtorImplemented, &fileModified)
						if ctorDetected {
							continue
						}

						destructorDetected := addDbOperationsIfDestructor(fn, t.Output, t.tokenSet, isTypeDestructorImplemented, &fileModified)
						if destructorDetected {
							continue
						}

					} else if fn.Name.Name != NobjectImplementationMethod && f.Name.Name != CustomIdImplementationMethod {

						typeName := getFunctionReceiverTypeAsString(fn.Recv)
						if isNobject := t.Output.IsNobjectInOrginalPackage[typeName]; isNobject {

							isGetterOrSetter := addDBOperationsIfGetterOrSetter(fn, t.Output, t.tokenSet, &fileModified)
							if isGetterOrSetter {
								continue
							}

							if !isFunctionStateless(fn.Recv) && retParamsVerifier.Check(fn) && !isDBGetOperationAlreadyAddedToMethod(fn.Body, t.tokenSet) {
								fileModified = true
								addDBOperationsToStateChangingMethod(fn, t.Output)
							}
						}
					}
				}
			}

			if !t.fileChanged[path] {
				t.fileChanged[path] = fileModified
			}
		}
	}
}

func addDBOperationsToStateChangingMethod(fn *ast.FuncDecl, parsedPackage ParsedPackage) {
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

func addDBOperationsIfGetterOrSetter(fn *ast.FuncDecl, parsedPackage ParsedPackage, set *token.FileSet, fileModified *bool) bool {
	var isGetterOrSetter bool
	typeName := getFunctionReceiverTypeAsString(fn.Recv)

	if strings.HasPrefix(fn.Name.Name, GetPrefix) {

		fieldName := strings.TrimPrefix(fn.Name.Name, GetPrefix)
		if fieldType, fieldExist := parsedPackage.TypeFields[typeName][fieldName]; fieldExist {
			isGetterOrSetter = true

			if !isGetFieldStmtAlreadyAdded(fn.Body, set) {
				idFieldName := getIdFieldNameOfType(typeName, parsedPackage.TypesWithCustomId)
				stmtsToInsert := getGetterDBStmts(fn, getDBStmtsParam{
					idFieldName:          idFieldName,
					typeName:             typeName,
					fieldName:            fieldName,
					fieldType:            fieldType,
					receiverVariableName: fn.Recv.List[0].Names[0].Name,
				})
				fn.Body.List = prependList(fn.Body.List, stmtsToInsert)
				*fileModified = true
			}
		}

	} else if strings.HasPrefix(fn.Name.Name, SetPrefix) {
		fieldName := strings.TrimPrefix(fn.Name.Name, SetPrefix)
		if fieldType, fieldExists := parsedPackage.TypeFields[typeName][fieldName]; fieldExists {
			isGetterOrSetter = true

			if !isSetFieldStmtAlreadyAdded(fn.Body, set) {
				idFieldName := getIdFieldNameOfType(typeName, parsedPackage.TypesWithCustomId)
				stmtsToInsert := getSetterDBStmts(fn, getDBStmtsParam{
					idFieldName:          idFieldName,
					typeName:             typeName,
					fieldName:            fieldName,
					fieldType:            fieldType,
					receiverVariableName: fn.Recv.List[0].Names[0].Name,
				})
				fn.Body.List = appendListBeforeLastElem(fn.Body.List, stmtsToInsert)
				*fileModified = true
			}
		}
	}

	return isGetterOrSetter
}

func addDBOperationsIfCtor(fn *ast.FuncDecl, parsedPackage ParsedPackage, set *token.FileSet, IsTypeNewCtorImplemented map[string]bool, IsTypeReNewCtorImplemented map[string]bool, fileModified *bool) bool {
	var ctorDetected bool

	if strings.HasPrefix(fn.Name.Name, ConstructorPrefix) {
		typeName := strings.TrimPrefix(fn.Name.Name, ConstructorPrefix)
		if parsedPackage.IsNobjectInOrginalPackage[typeName] {
			if !areDBOperationsAlreadyAddedToNewCtor(fn.Body, set) {
				idFieldName := getIdFieldNameOfType(typeName, parsedPackage.TypesWithCustomId)
				stmtsToInsert, err := getNewCtorStmts(fn, typeName, idFieldName)
				if err != nil {
					fmt.Println("wrong constructor definition of ", fn.Name.Name, ": ", err)
					return true
				}
				*fileModified = true
				fn.Body.List = appendListBeforeLastElem(fn.Body.List, stmtsToInsert)
			}

			IsTypeNewCtorImplemented[typeName] = true
			ctorDetected = true
		}
	}
	if strings.HasPrefix(fn.Name.Name, ReConstructorPrefix) {
		IsTypeReNewCtorImplemented[strings.TrimPrefix(fn.Name.Name, ConstructorPrefix)] = true
		ctorDetected = true
	}

	return ctorDetected
}

func addDbOperationsIfDestructor(fn *ast.FuncDecl, parsedPackage ParsedPackage, set *token.FileSet, isTypeDestructorImplemented map[string]bool, fileModified *bool) bool {
	var destructorDetected bool

	if strings.HasPrefix(fn.Name.Name, DestructorPrefix) {
		typeName := strings.TrimPrefix(fn.Name.Name, DestructorPrefix)
		if parsedPackage.IsNobjectInOrginalPackage[typeName] {
			if !areDBOperationsAlreadyAddedToDestructor(fn.Body, set) {
				idFieldName := getIdFieldNameOfType(typeName, parsedPackage.TypesWithCustomId)
				stmtsToInsert, err := getNewDestructorStmts(fn, typeName, idFieldName)
				if err != nil {
					fmt.Println("wrong destructor definition of ", fn.Name.Name, ": ", err)
					return true
				}
				*fileModified = true
				fn.Body.List = appendListBeforeLastElem(fn.Body.List, stmtsToInsert)
			}

			isTypeDestructorImplemented[typeName] = true
			destructorDetected = true
		}
	}

	return destructorDetected
}

func areDBOperationsAlreadyAddedToNewCtor(funcBlock *ast.BlockStmt, set *token.FileSet) bool {
	if funcBlock != nil && funcBlock.List != nil && len(funcBlock.List) > 3 {
		assign, _ := funcBlock.List[len(funcBlock.List)-4].(*ast.AssignStmt)
		secLastElem, _ := getFunctionBodyStmtAsString(set, assign)
		return strings.Contains(secLastElem, "lib.Insert")
	}
	return false
}

func areDBOperationsAlreadyAddedToDestructor(funcBlock *ast.BlockStmt, set *token.FileSet) bool {
	if funcBlock != nil && funcBlock.List != nil && len(funcBlock.List) > 2 {
		assign, _ := funcBlock.List[len(funcBlock.List)-3].(*ast.AssignStmt)
		secLastElem, _ := getFunctionBodyStmtAsString(set, assign)
		return strings.Contains(secLastElem, "lib.Delete")
	}
	return false
}

func isDBGetOperationAlreadyAddedToMethod(funcBlock *ast.BlockStmt, set *token.FileSet) bool {
	if funcBlock != nil && funcBlock.List != nil && len(funcBlock.List) > 2 {
		assign, _ := funcBlock.List[0].(*ast.AssignStmt)
		secLastElem, _ := getFunctionBodyStmtAsString(set, assign)
		return strings.Contains(secLastElem, "lib."+LibraryGetObjectStateMethod)
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

package parser

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"
)

func (t TypeSpecParser) adjustMethods(isTypeNewCtorImplemented map[string]bool, isTypeReNewCtorImplemented map[string]bool, isTypeDestructorImplemented map[string]bool) {

	for path, detectedFunctionsList := range t.detectedFunctions {
		for _, detectedFunction := range detectedFunctionsList {

			fn := detectedFunction.Function
			if fn.Name.Name == "Init" {
				continue
			}

			if fn.Recv == nil {

				// check if constructor, destructor

			} else if fn.Name.Name != NobjectImplementationMethod && fn.Name.Name != CustomIdImplementationMethod {

				typeName := getFunctionReceiverTypeAsString(fn.Recv)
				if isNobject := t.Output.IsNobjectInOrginalPackage[typeName]; isNobject {
					isGetter := t.addDBOperationsIfGetter(fn, path)
					if !isGetter {
						isSetter := t.addDBOperationsIfSetter(fn, path)
						if !isSetter {
							if !isFunctionStateless(fn.Recv) && retParamsVerifier.Check(fn) && !isInitFieldCheckAlreadyAddedAsFirstStmt(fn.Body, t.tokenSet) {
								t.fileChanged[path] = true
								t.addDBOperationsToStateChangingMethod(fn)
							}
						}
					}
				}
			}
		}
	}
}

func (t TypeSpecParser) addDBOperationsToStateChangingMethod(fn *ast.FuncDecl) {
	retrieveStateIfInitialized := getNobjectStateConditionalRetrieval(fn, t.Output)
	saveStateIfInitialized := getNobjectStateConditionalUpsert(fn, t.Output)
	fn.Body.List = prepend[ast.Stmt](fn.Body.List, &retrieveStateIfInitialized)
	fn.Body.List = prependBeforeLastElem[ast.Stmt](fn.Body.List, &saveStateIfInitialized)
}

func (t TypeSpecParser) addDBOperationsIfSetter(fn *ast.FuncDecl, path string) bool {
	typeName := getFunctionReceiverTypeAsString(fn.Recv)
	if strings.HasPrefix(fn.Name.Name, SetPrefix) {
		fieldName := strings.TrimPrefix(fn.Name.Name, SetPrefix)
		if fieldType, fieldExists := t.Output.TypeFields[typeName][fieldName]; fieldExists {
			if !isInitFieldCheckAlreadyAddedAsSecondLastStmt(fn.Body, t.tokenSet) {
				idFieldName := getIdFieldNameOfType(typeName, t.Output.TypesWithCustomId)
				saveInDbIfInitialized := getSetterDBStmts(fn, getDBStmtsParam{
					idFieldName:          idFieldName,
					typeName:             typeName,
					fieldName:            fieldName,
					fieldType:            fieldType,
					receiverVariableName: fn.Recv.List[0].Names[0].Name,
				})
				fn.Body.List = appendBeforeLastElem[ast.Stmt](fn.Body.List, &saveInDbIfInitialized)

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
		if fieldType, fieldExist := t.Output.TypeFields[typeName][fieldName]; fieldExist {

			if !isInitFieldCheckAlreadyAddedAsFirstStmt(fn.Body, t.tokenSet) {
				idFieldName := getIdFieldNameOfType(typeName, t.Output.TypesWithCustomId)
				retrieveFromDbIfInitialized := getGetterDBStmts(fn, getDBStmtsParam{
					idFieldName:          idFieldName,
					typeName:             typeName,
					fieldName:            fieldName,
					fieldType:            fieldType,
					receiverVariableName: fn.Recv.List[0].Names[0].Name,
				})
				fn.Body.List = prepend[ast.Stmt](fn.Body.List, &retrieveFromDbIfInitialized)
				t.fileChanged[path] = true
			}
			return true
		}
	}

	return false
}

func isInitFieldCheckAlreadyAddedAsFirstStmt(funcBlock *ast.BlockStmt, set *token.FileSet) bool {
	if funcBlock != nil && funcBlock.List != nil && len(funcBlock.List) > 0 {
		ifStmt, _ := funcBlock.List[0].(*ast.IfStmt)
		if ifStmt != nil {
			ifConditionAsString := types.ExprString(ifStmt.Cond)
			return strings.Contains(ifConditionAsString, IsInitializedFieldName)
		}
	}
	return false
}

func isInitFieldCheckAlreadyAddedAsSecondLastStmt(funcBlock *ast.BlockStmt, set *token.FileSet) bool {
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

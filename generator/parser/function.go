package parser

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"go/types"
	"strconv"
	"strings"
)

func isErrorTypeReturned(f *ast.FuncDecl) bool {
	return types.ExprString(f.Type.Results.List[len(f.Type.Results.List)-1].Type) == "error"
}

func getFunctionReceiverTypeAsString(fieldList *ast.FieldList) string {
	return strings.TrimPrefix(types.ExprString(fieldList.List[0].Type), "*")
}

func isFunctionStateless(fields *ast.FieldList) bool {
	return fields.List == nil || fields.List[0].Names == nil || fields.List[0].Names[0].Name == ""
}

func getFunctionReturnTypesAsString(results *ast.FieldList, isNobjectInOrgPkg map[string]bool) (string, error) {

	if len(results.List) > 2 {
		return "", fmt.Errorf("maximum allowed number of non-error return parameters is 1, found " + strconv.Itoa(len(results.List)))
	}
	if types.ExprString(results.List[len(results.List)-1].Type) != "error" {
		return "", fmt.Errorf("the last return parameter type of repository function must be error")
	}

	if len(results.List) == 1 {
		return "error", nil
	}

	// given the prev. conditions now we're sure
	// the case is: (optionalReturnType, error)
	returnParam := types.ExprString(results.List[0].Type)

	if _, isPresent := isNobjectInOrgPkg[returnParam]; isPresent {
		return "( " + OrginalPackageAlias + "." + returnParam + " , error)", nil
	}

	return "( " + returnParam + " , error)", nil
}

func getFunctionBodyStmtAsString(fset *token.FileSet, stmt *ast.AssignStmt) (string, error) {
	if stmt == nil {
		return "", fmt.Errorf("statement is nil")
	}
	var buf bytes.Buffer
	err := printer.Fprint(&buf, fset, stmt)
	if err != nil {
		return "", fmt.Errorf("error occurred when parsing the function body")
	}
	return buf.String(), nil
}

func getFunctionBodyAsString(fset *token.FileSet, body *ast.BlockStmt) (string, error) {
	var buf bytes.Buffer
	err := printer.Fprint(&buf, fset, body)
	if err != nil {
		return "", fmt.Errorf("error occurred when parsing the function body")
	}
	return buf.String(), nil
}

func getFirstFunctionReturnTypeAsString(fn *ast.FuncDecl) string {

	if fn.Type.Results == nil || fn.Type.Results.List == nil || len(fn.Type.Results.List) == 0 {
		return ""
	}

	return types.ExprString(fn.Type.Results.List[0].Type)
}

type returnParamsVerifier struct {
	errorsPrinted bool
}

var retParamsVerifier returnParamsVerifier

func init() {
	retParamsVerifier = *new(returnParamsVerifier)
}

func (v *returnParamsVerifier) Check(f *ast.FuncDecl) bool {

	if v.errorsPrinted {
		return f.Type.Results != nil && f.Type.Results.List != nil && isErrorTypeReturned(f) && len(f.Type.Results.List) <= 2
	} else {
		v.errorsPrinted = true
		if f.Type.Results == nil && f.Type.Results.List == nil && !isErrorTypeReturned(f) {
			fmt.Println("error type must be defined as the last return type from type's method. Handler generation for " + f.Name.Name + "skipped")
			return false
		}
		if len(f.Type.Results.List) > 2 {
			fmt.Println("Maximum allowed number of non-error return parameters is 1. Handler generation for " + f.Name.Name + "skipped")
			return false
		}
	}

	return true
}

package parser

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"go/types"
	"strings"
)

func isErrorTypeReturned(f *ast.FuncDecl) bool {
	return len(f.Type.Results.List) > 0 && types.ExprString(f.Type.Results.List[len(f.Type.Results.List)-1].Type) == "error"
}

func getFunctionReceiverTypeAsString(fieldList *ast.FieldList) string {
	return strings.TrimPrefix(types.ExprString(fieldList.List[0].Type), "*")
}

func isFunctionStateless(fields *ast.FieldList) bool {
	return fields.List == nil || fields.List[0].Names == nil || fields.List[0].Names[0].Name == ""
}

func getFunctionBodyAsString(fset *token.FileSet, body *ast.BlockStmt) (string, error) {
	var buf bytes.Buffer
	err := printer.Fprint(&buf, fset, body)
	if err != nil {
		return "", fmt.Errorf("error occurred when parsing the function body")
	}
	return buf.String(), nil
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
		if f.Type.Results == nil || f.Type.Results.List == nil || !isErrorTypeReturned(f) {
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

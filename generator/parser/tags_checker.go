package parser

import (
	"go/ast"
	"strings"
)

func isReadonly(field *ast.Field) bool {
	return field.Tag != nil && strings.Contains(field.Tag.Value, ReadonlyTag) && strings.Contains(field.Tag.Value, TagKey)
}

func isIndex(field *ast.Field) bool {
	return field.Tag != nil && strings.Contains(field.Tag.Value, IndexTag) && strings.Contains(field.Tag.Value, TagKey)
}

package parser

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"github.com/fatih/structtag"
)

func isReadonly(field *ast.Field) bool {
	return field.Tag != nil && strings.Contains(field.Tag.Value, ReadonlyTag) && strings.Contains(field.Tag.Value, NubesTagKey)
}

func getParsedTags(field *ast.Field) (*structtag.Tags, error) {
	if field.Tag != nil && field.Tag.Kind == token.STRING {
		unquotedTag, err := strconv.Unquote(field.Tag.Value)
		if err != nil {
			return nil, err
		}
		tags, err := structtag.Parse(unquotedTag)
		if err != nil {
			return nil, err
		}
		return tags, nil
	}
	return nil, nil
}

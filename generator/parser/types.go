package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"os"
)

func GetNobjectsDefinedInPack(path string) map[string]struct{} {
	set := token.NewFileSet()
	packs, err := parser.ParseDir(set, path, nil, 0)
	AssertDirParsed(err)

	nobjectTypes := make(map[string]struct{})

	for _, pack := range packs {
		for _, f := range pack.Files {
			for _, d := range f.Decls {
				if fn, isFn := d.(*ast.FuncDecl); isFn {
					if fn.Recv != nil && fn.Name.Name == GetTypeName {
						ownerType := types.ExprString(fn.Recv.List[0].Type)
						nobjectTypes[ownerType] = struct{}{}
					}
				}
			}
		}
	}

	return nobjectTypes
}

func AssertDirParsed(err error) {
	if err != nil {
		fmt.Println("Failed to parse files in the directory: %w", err)
		os.Exit(1)
	}
}

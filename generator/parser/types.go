package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
)

type TypeDeclaration struct {
}

func PrepareTypes(path string) {
	set := token.NewFileSet()
	packs, err := parser.ParseDir(set, path, nil, 0)
	AssertDirParsed(err)

	structs := []*ast.StructType{}

	for _, pack := range packs {
		for _, f := range pack.Files {
			ast.Inspect(f, func(n ast.Node) bool {
				if n, ok := n.(*ast.StructType); ok {
					structs = append(structs, n)
				}
				return true
			})
		}
	}

	for _, i := range structs {
		fmt.Println()
		fmt.Println("NEXT STRUCT")
		printer.Fprint(os.Stdout, set, i)
	}

	detectedTypes := []*ast.TypeSpec{}
	for _, pack := range packs {
		for _, f := range pack.Files {
			ast.Inspect(f, func(n ast.Node) bool {
				if n, ok := n.(*ast.TypeSpec); ok {
					detectedTypes = append(detectedTypes, n)
				}
				return true
			})
		}
	}

	for _, i := range detectedTypes {
		printer.Fprint(os.Stdout, set, i)
	}
}

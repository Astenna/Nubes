package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
)

type HandlerFunc struct {
	Imports           string
	FunctionSignature string
	TypeInstantiation string
	FunctionBody      string
	TypeWrite         string
}

type TypeDeclaration struct {
}

func ParseTypes(path string) {
	set := token.NewFileSet()
	packs, err := parser.ParseDir(set, path, nil, 0)
	if err != nil {
		fmt.Println("Failed to parse package:", err)
		os.Exit(1)
	}
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
		// printer.Fprint(os.Stdout, set, i.Fields.List[0].Type)
		// printer.Fprint(os.Stdout, set, i.Fields.List[0].Names[0].Name)
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
		fmt.Println()
		fmt.Println("NEXT TYPE")
		printer.Fprint(os.Stdout, set, i)
		// printer.Fprint(os.Stdout, set, i.Fields.List[0].Type)
		// printer.Fprint(os.Stdout, set, i.Fields.List[0].Names[0].Name)
	}
}

func PrepareFunctions(path string) []HandlerFunc {
	set := token.NewFileSet()
	packs, err := parser.ParseDir(set, path, nil, 0)
	if err != nil {
		fmt.Println("Failed to parse package:", err)
		os.Exit(1)
	}

	funcs := []*ast.FuncDecl{}
	for _, pack := range packs {
		for _, f := range pack.Files {
			for _, d := range f.Decls {
				if fn, isFn := d.(*ast.FuncDecl); isFn {
					funcs = append(funcs, fn)
				}
			}
		}
	}

	for _, i := range funcs {
		i.Name.Name = i.Name.Name + "_MODIFIED"
		fmt.Println()
		fmt.Println("NEXT FUNC")
		fmt.Println("RECEIVER")
		fmt.Println()
		if i.Recv == nil {
			fmt.Println("Receiver was null")
		} else {
			printer.Fprint(os.Stdout, set, i.Recv.List[0].Type)
			if i.Recv.List[0].Names != nil {
				printer.Fprint(os.Stdout, set, i.Recv.List[0].Names[0].Name)
				fmt.Println(i.Recv.List[0].Names[0].Name)
				name := i.Recv.List[0].Names[0].Name
				_ = name
				fmt.Println()
			} else {
				fmt.Println("Names was null")
			}
		}
		fmt.Println()
		//printer.Fprint(os.Stdout, set, i)
	}
	return nil
}

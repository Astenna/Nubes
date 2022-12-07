package parser

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"go/types"
	"os"
	"strings"
)

type HandlerFunc struct {
	Imports       string
	Signature     string
	Prolog        string
	Body          string
	Epilog        string
	HandlerName   string
	OwnerType     string
	OwnerTypeName string
}

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
		printer.Fprint(os.Stdout, set, i)
	}
}

func PrepareHandlerFunctions(path string) []HandlerFunc {
	set := token.NewFileSet()
	packs, err := parser.ParseDir(set, path, nil, 0)
	AssertDirParsed(err)

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

	handlerFuncs := []HandlerFunc{}
	for _, f := range funcs {
		if f.Recv == nil || f.Name.Name == "GetTypeName" {
			continue
		}

		f.Name.Name = f.Name.Name + "Handler"
		newHandler := HandlerFunc{
			HandlerName: f.Name.Name,
			OwnerType:   types.ExprString(f.Recv.List[0].Type),
			Signature:   "func " + f.Name.Name + types.ExprString(f.Type),
		}

		if f.Recv.List[0].Names != nil {
			newHandler.OwnerTypeName = f.Recv.List[0].Names[0].Name
		}
		f.Recv = nil

		buf := new(bytes.Buffer)
		printer.Fprint(buf, set, f.Body)
		newHandler.Body = strings.Trim(buf.String(), "{}")
		newHandler.Prolog = "lib.Get[types.Shop](id)"     // TODO
		newHandler.Epilog = "lib.Write[types.Shop](shop)" // TODO

		handlerFuncs = append(handlerFuncs, newHandler)
	}
	return handlerFuncs
}

func AssertDirParsed(err error) {
	if err != nil {
		fmt.Println("Failed to parse files in the directory", err)
		os.Exit(1)
	}
}

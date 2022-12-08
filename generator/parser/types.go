package parser

import (
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
	OrginalPackage       string
	Imports              string
	Signature            string
	OwnerVariableName    string
	OwnerType            string
	ReturnFromInvocation string
	Invocation           string
	HandlerName          string
	Stateless            bool
	OptionalReturnVar    string
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

	funcsMap := make(map[string][]*ast.FuncDecl)

	for packageName, pack := range packs {
		for _, f := range pack.Files {
			for _, d := range f.Decls {
				if fn, isFn := d.(*ast.FuncDecl); isFn {

					if funcsMap[packageName] == nil {
						funcsMap[packageName] = []*ast.FuncDecl{}
					}
					funcsMap[packageName] = append(funcsMap[packageName], fn)
				}
			}
		}
	}

	handlerFuncs := []HandlerFunc{}
	for packageName, funcs := range funcsMap {
		for _, f := range funcs {
			if f.Recv == nil || f.Name.Name == "GetTypeName" {
				continue
			}

			signature := strings.SplitAfter("(id string, "+strings.TrimPrefix(types.ExprString(f.Type), "func("), ")")[0]
			newHandler := HandlerFunc{
				OrginalPackage: packageName,
				HandlerName:    f.Name.Name + "Handler",
				Signature:      "func " + f.Name.Name + "Handler" + signature,
			}

			// 4 cases:
			// C1: no return parameters
			// C2: 1 return: error
			// C3: 1 return: non-error
			// C4: 2 return: non-error, error
			errorTypeFound := false
			if f.Type.Results == nil {
				// C1
				newHandler.Signature += " error"
			} else {
				errorTypeFound = types.ExprString(f.Type.Results.List[len(f.Type.Results.List)-1].Type) == "error"
				if !errorTypeFound && len(f.Type.Results.List) >= 1 {
					fmt.Println("Maximum allowed number of non-error return parameters is 1. Handler generation for " + f.Name.Name + "skipped")
					continue
				} else if !errorTypeFound {
					// C3
					newHandler.Signature += "(" + types.ExprString(f.Type.Results.List[0].Type) + ", error)"
					newHandler.ReturnFromInvocation = "result :="
					newHandler.OptionalReturnVar = "result"
				} else {
					if len(f.Type.Results.List) == 1 {
						// C2
						newHandler.Signature += " error"
						newHandler.ReturnFromInvocation = "err :="
					} else {
						// C4
						newHandler.Signature += "(" + types.ExprString(f.Type.Results.List[0].Type) + " ,error)"
						newHandler.ReturnFromInvocation = "result, err :="
						newHandler.OptionalReturnVar = "result"
					}
				}
			}
			_ = errorTypeFound

			newHandler.OwnerType = strings.TrimPrefix(types.ExprString(f.Recv.List[0].Type), "*")
			var ownerTypeName string
			if f.Recv.List[0].Names == nil {
				// stateless method, instance will be created just to invoke the method
				newHandler.Stateless = true
				newHandler.OwnerVariableName = "typeInstance"
				ownerTypeName = "typeInstance"
			} else {
				newHandler.Stateless = false
				// stateful method, create instance to invoke the method and then save state changes
				ownerTypeName = f.Recv.List[0].Names[0].Name
				newHandler.OwnerVariableName = f.Recv.List[0].Names[0].Name
			}
			newHandler.Invocation = ownerTypeName + "." + f.Name.Name + "(" + GetPareterNames(f.Type.Params) + ")"

			handlerFuncs = append(handlerFuncs, newHandler)
		}
	}

	return handlerFuncs
}

func GetPareterNames(params *ast.FieldList) string {
	var names []string
	for _, param := range params.List {
		for _, name := range param.Names {
			names = append(names, name.Name)
		}
	}
	return strings.Join(names, ", ")
}

func AssertDirParsed(err error) {
	if err != nil {
		fmt.Println("Failed to parse files in the directory", err)
		os.Exit(1)
	}
}

package main

import {
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
}

func main() {
	srcPath := os.Args[1]

	// Parse the program into an AST
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, srcPath, nil, 0)
	if err != nil {
		panic(err)
	}

	// Loop over all the declarations in the program
	for _, f := range node.Decls {
		fn, ok := f.(*ast.FuncDecl)
		if !ok {
			continue
		}

		a2 := ast.ExprStmt{ // Create print statement to insert
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.Ident{Name: "fmt"},
					Sel: &ast.Ident{Name: "Printf"},
				},
				Args: []ast.Expr{
					&ast.BasicLit{Kind: token.STRING, Value: "\"" + fn.Name.Name + "\""},
				},
			},
		}

		// Append print statement to function body
		fn.Body.List = append([]ast.Stmt{&a2}, fn.Body.List...)
	}

	// Write our changed AST to a new file in the same directory
	f, err := os.Create("program-altered.go")
	defer f.Close()

	if err := printer.Fprint(f, fset, node); err != nil {
		log.Fatal(err)
	}
}
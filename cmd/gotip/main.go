package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	rootDir := "./testdata"
	return filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || !strings.HasSuffix(path, "_test.go") {
			return nil
		}
		processFile(path)
		return nil
	})
}

func processFile(path string) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		fmt.Printf("Error parsing file %s: %v\n", path, err)
		return
	}
	for _, decl := range node.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || !strings.HasPrefix(fn.Name.Name, "Test") {
			continue
		}
		fmt.Printf("Test Function: %s (in %s)\n", fn.Name.Name, path)

		ast.Inspect(fn.Body, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok || sel.Sel.Name != "Run" {
				return true
			}
			subtestName := "<?>"
			if arg, ok := call.Args[0].(*ast.BasicLit); ok {
				subtestName = strings.Trim(arg.Value, `"`)
			}
			fmt.Printf("  - %s\n", subtestName)

			return true
		})
	}
}

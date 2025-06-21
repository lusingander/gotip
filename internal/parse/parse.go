package parse

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type TestFunction struct {
	path string
	name string
	subs []*SubTest
}

type SubTest struct {
	name string
	subs []*SubTest
}

func ProcessFile(path string) ([]*TestFunction, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file %s: %w", path, err)
	}
	testFunctions := make([]*TestFunction, 0)
	for _, decl := range node.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || !strings.HasPrefix(fn.Name.Name, "Test") {
			continue
		}
		subs, err := parseTestFunction(fn)
		if err != nil {
			return nil, fmt.Errorf("failed to parse test function %s: %w", fn.Name.Name, err)
		}
		testFunction := &TestFunction{
			name: fn.Name.Name,
			path: path,
			subs: subs,
		}
		testFunctions = append(testFunctions, testFunction)
	}
	return testFunctions, nil
}

func parseTestFunction(fn *ast.FuncDecl) ([]*SubTest, error) {
	subs := make([]*SubTest, 0)
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
		subtest := &SubTest{
			name: subtestName,
			subs: nil, // todo: parse nested subtests if needed
		}
		subs = append(subs, subtest)
		return true
	})
	return subs, nil
}

func PrintTestFunctions(tests map[string][]*TestFunction) {
	for path, testFunctions := range tests {
		fmt.Printf("# %s\n", path)
		for _, tf := range testFunctions {
			printTestFunction(tf)
		}
		fmt.Println()
	}
}

func printTestFunction(tf *TestFunction) {
	fmt.Printf("- %s\n", tf.name)
	if len(tf.subs) > 0 {
		printSubTests(tf.subs, "  ")
	}
}

func printSubTests(subs []*SubTest, indent string) {
	for _, sub := range subs {
		fmt.Printf("%s- %s\n", indent, sub.name)
		if len(sub.subs) > 0 {
			printSubTests(sub.subs, indent+"  ")
		}
	}
}

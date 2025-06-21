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
		testFunction := &TestFunction{
			name: fn.Name.Name,
			path: path,
			subs: parseTestFunction(fn.Body),
		}
		testFunctions = append(testFunctions, testFunction)
	}
	return testFunctions, nil
}

func parseTestFunction(fnBody *ast.BlockStmt) []*SubTest {
	subs := make([]*SubTest, 0)
	for _, stmt := range fnBody.List {
		subs = append(subs, findSubTestsInStmt(stmt)...)
	}
	return subs
}

func findSubTestsInStmt(stmt ast.Stmt) []*SubTest {
	subs := make([]*SubTest, 0)
	switch s := stmt.(type) {
	case *ast.ExprStmt:
		call, ok := s.X.(*ast.CallExpr)
		if !ok {
			return nil
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != "Run" || len(call.Args) < 2 {
			return nil
		}
		subtestName := "<unknown>"
		if arg, ok := call.Args[0].(*ast.BasicLit); ok {
			subtestName = strings.Trim(arg.Value, `"`)
		}
		subtest := &SubTest{
			name: subtestName,
			subs: nil,
		}
		if fnLit, ok := call.Args[1].(*ast.FuncLit); ok {
			subtest.subs = parseTestFunction(fnLit.Body)
		}
		subs = append(subs, subtest)
	case *ast.BlockStmt:
		for _, innerStmt := range s.List {
			subs = append(subs, findSubTestsInStmt(innerStmt)...)
		}
	case *ast.ForStmt:
		for _, innerStmt := range s.Body.List {
			subs = append(subs, findSubTestsInStmt(innerStmt)...)
		}
	case *ast.RangeStmt:
		for _, innerStmt := range s.Body.List {
			subs = append(subs, findSubTestsInStmt(innerStmt)...)
		}
	}
	return subs
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
		printSubTests(tf.subs, 1)
	}
}

func printSubTests(subs []*SubTest, indentLevel int) {
	indent := strings.Repeat("  ", indentLevel)
	for _, sub := range subs {
		fmt.Printf("%s- %s\n", indent, sub.name)
		if len(sub.subs) > 0 {
			printSubTests(sub.subs, indentLevel+1)
		}
	}
}

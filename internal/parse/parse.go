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
		testFunctions = append(testFunctions, processTestFunction(fn, path))
	}
	return testFunctions, nil
}

func processTestFunction(fn *ast.FuncDecl, path string) *TestFunction {
	unresolvedSubTests := findSubTests(fn.Body.List)

	subs := make([]*SubTest, len(unresolvedSubTests))
	for i, sub := range unresolvedSubTests {
		subs[i] = sub.resolve()
	}

	return &TestFunction{
		name: fn.Name.Name,
		path: path,
		subs: subs,
	}
}

func findSubTests(stmts []ast.Stmt) []*unresolvedSubTest {
	subs := make([]*unresolvedSubTest, 0)
	for _, stmt := range stmts {
		switch s := stmt.(type) {
		case *ast.ExprStmt:
			call, ok := s.X.(*ast.CallExpr)
			if !ok {
				continue
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok || sel.Sel.Name != "Run" || len(call.Args) < 2 {
				continue
			}
			subs = append(subs, findSubTest(call.Args))
		case *ast.BlockStmt:
			subs = append(subs, findSubTests(s.List)...)
		case *ast.ForStmt:
			subs = append(subs, findSubTests(s.Body.List)...)
		case *ast.RangeStmt:
			subs = append(subs, findSubTests(s.Body.List)...)
		}
	}
	return subs
}

func findSubTest(exprs []ast.Expr) *unresolvedSubTest {
	var name unresolvedSubTestName = &unknownSubTestName{}

	switch e := exprs[0].(type) {
	case *ast.BasicLit:
		if e.Kind == token.STRING {
			name = &literalSubTestName{
				name: strings.Trim(e.Value, `"`),
			}
		}
	case *ast.SelectorExpr:
		name = &selectorSubTestName{
			receiver: e.X.(*ast.Ident).Name,
			field:    e.Sel.Name,
		}
	case *ast.Ident:
		name = &identSubTestName{
			name: e.Name,
		}
	case *ast.BinaryExpr:
		// If it's a binary expression, we can't resolve it to a specific name without more context.
		// This might happen in cases like `t.Run("test"+string(i), ...)`.
	case *ast.CallExpr:
		// If it's a call expression, we can't resolve it to a specific name without more context
		// This might happen in cases like `t.Run(fmt.Sprintf("test%d", i), ...)`.
	}

	var subs []*unresolvedSubTest
	if fnLit, ok := exprs[1].(*ast.FuncLit); ok {
		subs = findSubTests(fnLit.Body.List)
	}

	return &unresolvedSubTest{
		name: name,
		subs: subs,
	}
}

type unresolvedSubTest struct {
	name unresolvedSubTestName
	subs []*unresolvedSubTest
}

func (t *unresolvedSubTest) resolve() *SubTest {
	subTests := make([]*SubTest, len(t.subs))
	for i, sub := range t.subs {
		subTests[i] = sub.resolve()
	}
	return &SubTest{
		name: t.name.resolveTestName(),
		subs: subTests,
	}
}

type unresolvedSubTestName interface {
	resolveTestName() string
}

type literalSubTestName struct {
	name string
}

func (l *literalSubTestName) resolveTestName() string {
	return l.name
}

type selectorSubTestName struct {
	receiver string
	field    string
}

func (s *selectorSubTestName) resolveTestName() string {
	// todo: resolve the selector to a specific name if possible
	return fmt.Sprintf("%s.%s", s.receiver, s.field)
}

type identSubTestName struct {
	name string
}

func (i *identSubTestName) resolveTestName() string {
	// todo: resolve the identifier to a specific name if possible
	return i.name
}

type unknownSubTestName struct{}

func (u *unknownSubTestName) resolveTestName() string {
	return "<unknown>"
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

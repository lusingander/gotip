package parse

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type TestFunction struct {
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
		testFunctions = append(testFunctions, processTestFunction(fn))
	}
	return testFunctions, nil
}

func processTestFunction(fn *ast.FuncDecl) *TestFunction {
	unresolvedSubTests := findSubTests(fn.Body.List)

	subs := make([]*SubTest, 0)
	for _, sub := range unresolvedSubTests {
		subs = append(subs, sub.resolve()...)
	}

	return &TestFunction{
		name: fn.Name.Name,
		subs: subs,
	}
}

func findSubTests(stmts []ast.Stmt, cs ...subTestContext) []*unresolvedSubTest {
	newCs := append([]subTestContext{}, cs...)
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
			subs = append(subs, findSubTest(call.Args, newCs...))
		case *ast.BlockStmt:
			subs = append(subs, findSubTests(s.List, newCs...)...)
		case *ast.ForStmt:
			subs = append(subs, findSubTests(s.Body.List, newCs...)...)
		case *ast.RangeStmt:
			if c := forRangeContextFromRangeStmt(s); c != nil {
				newCs = append(newCs, c)
			}
			subs = append(subs, findSubTests(s.Body.List, newCs...)...)
		case *ast.AssignStmt:
			if c := findStructSliceLiteralDeclaration(s); c != nil {
				newCs = append(newCs, c)
			}
		}
	}
	return subs
}

func findStructSliceLiteralDeclaration(assign *ast.AssignStmt) *structSliceLiteralDeclarationContext {
	if len(assign.Lhs) != 1 || len(assign.Rhs) != 1 {
		return nil
	}
	ident, ok := assign.Lhs[0].(*ast.Ident)
	if !ok {
		return nil
	}
	compLit, ok := assign.Rhs[0].(*ast.CompositeLit)
	if !ok {
		return nil
	}
	if _, ok := compLit.Type.(*ast.ArrayType); !ok {
		return nil
	}
	return &structSliceLiteralDeclarationContext{
		ident:   ident.Name,
		compLit: compLit,
	}
}

func findSubTest(exprs []ast.Expr, cs ...subTestContext) *unresolvedSubTest {
	var name unresolvedSubTestName = &unknownSubTestName{}

	switch e := exprs[0].(type) {
	case *ast.BasicLit:
		if e.Kind == token.STRING {
			name = &literalSubTestName{
				name: strings.Trim(e.Value, `"`),
			}
		}
	case *ast.SelectorExpr:
		n := &selectorSubTestName{
			receiver: e.X.(*ast.Ident).Name,
			field:    e.Sel.Name,
		}
		for _, c := range cs {
			if forRangeCtx, ok := c.(*forRangeContext); ok {
				if n.receiver == forRangeCtx.valueIdent {
					for _, c := range cs {
						if structSliceCtx, ok := c.(*structSliceLiteralDeclarationContext); ok {
							if structSliceCtx.ident == forRangeCtx.iterIdent {
								n.cases = structSliceCtx.extractTestCaseName(n.field)
							}
						}
					}
				}
			}
		}
		name = n
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

func (t *unresolvedSubTest) resolve() []*SubTest {
	subTests := make([]*SubTest, 0)
	for _, sub := range t.subs {
		subTests = append(subTests, sub.resolve()...)
	}
	tests := make([]*SubTest, 0)
	for _, n := range t.name.resolveTestName() {
		test := &SubTest{
			name: n,
			subs: subTests,
		}
		tests = append(tests, test)
	}
	return tests
}

type unresolvedSubTestName interface {
	resolveTestName() []string
}

type literalSubTestName struct {
	name string
}

func (l *literalSubTestName) resolveTestName() []string {
	return []string{l.name}
}

type selectorSubTestName struct {
	receiver string
	field    string
	cases    []string
}

func (s *selectorSubTestName) resolveTestName() []string {
	if len(s.cases) == 0 {
		return []string{"<unknown>"}
	}
	return s.cases
}

type identSubTestName struct {
	name string
}

func (i *identSubTestName) resolveTestName() []string {
	// todo: resolve the identifier to a specific name if possible
	return []string{"<unknown>"}
}

type unknownSubTestName struct{}

func (u *unknownSubTestName) resolveTestName() []string {
	return []string{"<unknown>"}
}

type subTestContext interface{}

type structSliceLiteralDeclarationContext struct {
	ident   string
	compLit *ast.CompositeLit
}

func (c *structSliceLiteralDeclarationContext) extractTestCaseName(name string) []string {
	caseFieldIdx := -1
	if st, ok := c.compLit.Type.(*ast.ArrayType).Elt.(*ast.StructType); ok {
		for i, field := range st.Fields.List {
			if len(field.Names) == 1 && field.Names[0].Name == name {
				caseFieldIdx = i
				break
			}
		}
	}
	ns := make([]string, 0)
	for _, elt := range c.compLit.Elts {
		st, ok := elt.(*ast.CompositeLit)
		if !ok {
			continue
		}
		for i, elt := range st.Elts {
			switch e := elt.(type) {
			case *ast.BasicLit:
				if e.Kind == token.STRING && i == caseFieldIdx {
					n := strings.Trim(e.Value, `"`)
					ns = append(ns, n)
				}
			case *ast.KeyValueExpr:
				if keyIdent, ok := e.Key.(*ast.Ident); ok {
					if keyIdent.Name == name {
						if e, ok := e.Value.(*ast.BasicLit); ok {
							if e.Kind == token.STRING {
								n := strings.Trim(e.Value, `"`)
								ns = append(ns, n)
							}
						}
					}
				}
			}
		}
	}
	return ns
}

type forRangeContext struct {
	keyIdent   string
	valueIdent string
	iterIdent  string
}

func forRangeContextFromRangeStmt(stmt *ast.RangeStmt) *forRangeContext {
	if stmt == nil {
		return nil
	}
	var keyName, valueName, iterName string
	if stmt.Key != nil {
		if ident, ok := stmt.Key.(*ast.Ident); ok {
			keyName = ident.Name
		}
	}
	if stmt.Value != nil {
		if ident, ok := stmt.Value.(*ast.Ident); ok {
			valueName = ident.Name
		}
	}
	if stmt.X != nil {
		if ident, ok := stmt.X.(*ast.Ident); ok {
			iterName = ident.Name
		}
	}
	if keyName == "" || valueName == "" || iterName == "" {
		return nil
	}
	return &forRangeContext{
		keyIdent:   keyName,
		valueIdent: valueName,
		iterIdent:  iterName,
	}
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

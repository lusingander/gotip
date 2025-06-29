package parse

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"slices"
	"strings"
)

type TestFunction struct {
	Name string
	Subs []*SubTest
}

type SubTest struct {
	Name                 string
	Subs                 []*SubTest
	IsUnresolvedSubTests bool
}

var ignore = []string{
	"vendor",
	"testdata",
}

func ProcessFilesRecursively(rootDir string) (map[string][]*TestFunction, error) {
	tests := make(map[string][]*TestFunction)
	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() && slices.Contains(ignore, d.Name()) {
			return filepath.SkipDir
		}
		if err != nil || !strings.HasSuffix(path, "_test.go") {
			return nil
		}
		testFunctions, err := processFile(path)
		if err != nil {
			return fmt.Errorf("error processing file %s: %w", path, err)
		}
		tests[path] = testFunctions
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tests, nil
}

func processFile(path string) ([]*TestFunction, error) {
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
		Name: fn.Name.Name,
		Subs: subs,
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
	var name unresolvedSubTestName

	switch e := exprs[0].(type) {
	case *ast.BasicLit:
		name = findSubTestNameFromBasicLit(e)
	case *ast.SelectorExpr:
		name = findSubTestNameFromSelectorExpr(e, cs...)
	case *ast.Ident:
		name = findSubTestNameFromIdent(e)
	case *ast.BinaryExpr:
		// If it's a binary expression, we can't resolve it to a specific name without more context.
		// This might happen in cases like `t.Run("test"+string(i), ...)`.
	case *ast.CallExpr:
		// If it's a call expression, we can't resolve it to a specific name without more context
		// This might happen in cases like `t.Run(fmt.Sprintf("test%d", i), ...)`.
	}

	if name == nil {
		name = &unknownSubTestName{}
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

func findSubTestNameFromBasicLit(lit *ast.BasicLit) *literalSubTestName {
	if lit.Kind != token.STRING {
		return nil
	}
	return &literalSubTestName{
		name: strings.Trim(lit.Value, `"`),
	}
}

func findSubTestNameFromSelectorExpr(sel *ast.SelectorExpr, cs ...subTestContext) *selectorSubTestName {
	n := &selectorSubTestName{
		receiver: sel.X.(*ast.Ident).Name,
		field:    sel.Sel.Name,
	}
	for _, c := range cs {
		forRangeCtx, ok := c.(*forRangeContext)
		if !ok {
			continue
		}
		if n.receiver != forRangeCtx.valueIdent {
			continue
		}
		for _, c := range cs {
			structSliceCtx, ok := c.(*structSliceLiteralDeclarationContext)
			if !ok {
				continue
			}
			if structSliceCtx.ident != forRangeCtx.iterIdent {
				continue
			}
			n.cases = structSliceCtx.extractTestCaseName(n.field)
		}
	}
	return n
}

func findSubTestNameFromIdent(ident *ast.Ident) *identSubTestName {
	return &identSubTestName{
		name: ident.Name,
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
	ns, resolved := t.name.resolveTestName()
	for _, n := range ns {
		test := &SubTest{
			Name:                 n,
			Subs:                 subTests,
			IsUnresolvedSubTests: !resolved,
		}
		tests = append(tests, test)
	}
	return tests
}

type unresolvedSubTestName interface {
	resolveTestName() ([]string, bool)
}

type literalSubTestName struct {
	name string
}

func (l *literalSubTestName) resolveTestName() ([]string, bool) {
	return []string{l.name}, true
}

type selectorSubTestName struct {
	receiver string
	field    string
	cases    []string
}

func (s *selectorSubTestName) resolveTestName() ([]string, bool) {
	if len(s.cases) == 0 {
		return []string{"<unknown>"}, false
	}
	return s.cases, true
}

type identSubTestName struct {
	name string
}

func (i *identSubTestName) resolveTestName() ([]string, bool) {
	// todo: resolve the identifier to a specific name if possible
	return []string{"<unknown>"}, false
}

type unknownSubTestName struct{}

func (u *unknownSubTestName) resolveTestName() ([]string, bool) {
	return []string{"<unknown>"}, false
}

type subTestContext interface{}

type structSliceLiteralDeclarationContext struct {
	ident   string
	compLit *ast.CompositeLit
}

func (c *structSliceLiteralDeclarationContext) extractTestCaseName(name string) []string {
	caseFieldIdx := c.findCaseNameFieldIndex(name)
	ns := make([]string, 0)
	for _, elt := range c.compLit.Elts {
		st, ok := elt.(*ast.CompositeLit)
		if !ok {
			continue
		}
		for i, elt := range st.Elts {
			switch e := elt.(type) {
			case *ast.BasicLit:
				if i == caseFieldIdx && e.Kind == token.STRING {
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

func (c *structSliceLiteralDeclarationContext) findCaseNameFieldIndex(name string) int {
	if st, ok := c.compLit.Type.(*ast.ArrayType).Elt.(*ast.StructType); ok {
		for i, field := range st.Fields.List {
			if len(field.Names) == 1 && field.Names[0].Name == name {
				return i
			}
		}
	}
	return -1
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
	fmt.Printf("- %s\n", tf.Name)
	if len(tf.Subs) > 0 {
		printSubTests(tf.Subs, 1)
	}
}

func printSubTests(subs []*SubTest, indentLevel int) {
	indent := strings.Repeat("  ", indentLevel)
	for _, sub := range subs {
		fmt.Printf("%s- %s\n", indent, sub.Name)
		if len(sub.Subs) > 0 {
			printSubTests(sub.Subs, indentLevel+1)
		}
	}
}

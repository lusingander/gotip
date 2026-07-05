package tip

import "testing"

func TestFilterTestsByPackages(t *testing.T) {
	tests := map[string][]*TestFunction{
		"./internal/parse/parse_test.go": {
			{Name: "TestParse"},
		},
		"./internal/tip/model_test.go": {
			{Name: "TestTarget"},
		},
		"./main_test.go": {
			{Name: "TestMain"},
		},
	}

	got := FilterTestsByPackages(tests, []string{"./internal/parse"})

	if len(got) != 1 {
		t.Fatalf("filtered tests len = %d, want 1", len(got))
	}
	if _, ok := got["./internal/parse/parse_test.go"]; !ok {
		t.Fatal("filtered tests does not contain ./internal/parse/parse_test.go")
	}
}

func TestFilterTestsByPackages_normalizesPackageNames(t *testing.T) {
	tests := map[string][]*TestFunction{
		"./internal/parse/parse_test.go": {
			{Name: "TestParse"},
		},
	}

	got := FilterTestsByPackages(tests, []string{"internal/parse/"})

	if len(got) != 1 {
		t.Fatalf("filtered tests len = %d, want 1", len(got))
	}
	if _, ok := got["./internal/parse/parse_test.go"]; !ok {
		t.Fatal("filtered tests does not contain ./internal/parse/parse_test.go")
	}
}

func TestFilterTestsByPackages_includesSubpackages(t *testing.T) {
	tests := map[string][]*TestFunction{
		"./internal/parse/parse_test.go": {
			{Name: "TestParse"},
		},
		"./internal/tip/model_test.go": {
			{Name: "TestTarget"},
		},
		"./internalx/x_test.go": {
			{Name: "TestX"},
		},
	}

	got := FilterTestsByPackages(tests, []string{"./internal"})

	if len(got) != 2 {
		t.Fatalf("filtered tests len = %d, want 2", len(got))
	}
	if _, ok := got["./internal/parse/parse_test.go"]; !ok {
		t.Fatal("filtered tests does not contain ./internal/parse/parse_test.go")
	}
	if _, ok := got["./internal/tip/model_test.go"]; !ok {
		t.Fatal("filtered tests does not contain ./internal/tip/model_test.go")
	}
	if _, ok := got["./internalx/x_test.go"]; ok {
		t.Fatal("filtered tests contains ./internalx/x_test.go")
	}
}

func TestFilterTestsByPackages_supportsGoStylePackagePattern(t *testing.T) {
	tests := map[string][]*TestFunction{
		"./internal/parse/parse_test.go": {
			{Name: "TestParse"},
		},
		"./internal/tip/model_test.go": {
			{Name: "TestTarget"},
		},
		"./cmd/gotip/main_test.go": {
			{Name: "TestMain"},
		},
	}

	got := FilterTestsByPackages(tests, []string{"internal/..."})

	if len(got) != 2 {
		t.Fatalf("filtered tests len = %d, want 2", len(got))
	}
	if _, ok := got["./internal/parse/parse_test.go"]; !ok {
		t.Fatal("filtered tests does not contain ./internal/parse/parse_test.go")
	}
	if _, ok := got["./internal/tip/model_test.go"]; !ok {
		t.Fatal("filtered tests does not contain ./internal/tip/model_test.go")
	}
	if _, ok := got["./cmd/gotip/main_test.go"]; ok {
		t.Fatal("filtered tests contains ./cmd/gotip/main_test.go")
	}
}

func TestFilterTestsByPackages_allPackagePatternReturnsAllTests(t *testing.T) {
	tests := map[string][]*TestFunction{
		"./internal/parse/parse_test.go": {
			{Name: "TestParse"},
		},
		"./main_test.go": {
			{Name: "TestMain"},
		},
	}

	got := FilterTestsByPackages(tests, []string{"./..."})

	if len(got) != len(tests) {
		t.Fatalf("filtered tests len = %d, want %d", len(got), len(tests))
	}
	if _, ok := got["./internal/parse/parse_test.go"]; !ok {
		t.Fatal("filtered tests does not contain ./internal/parse/parse_test.go")
	}
	if _, ok := got["./main_test.go"]; !ok {
		t.Fatal("filtered tests does not contain ./main_test.go")
	}
}

func TestFilterTestsByPackages_emptyPackagesReturnsOriginalTests(t *testing.T) {
	tests := map[string][]*TestFunction{
		"./internal/parse/parse_test.go": {
			{Name: "TestParse"},
		},
	}

	got := FilterTestsByPackages(tests, nil)

	if len(got) != len(tests) {
		t.Fatalf("filtered tests len = %d, want %d", len(got), len(tests))
	}
	if _, ok := got["./internal/parse/parse_test.go"]; !ok {
		t.Fatal("filtered tests does not contain ./internal/parse/parse_test.go")
	}
}

func TestFilterHistoriesByPackages(t *testing.T) {
	histories := &Histories{
		ProjectDir: "/path/to/project",
		Histories: []*History{
			{
				Path:            "./internal/parse/parse_test.go",
				PackageName:     "./internal/parse",
				TestNamePattern: "TestParse",
			},
			{
				Path:            "./internal/tip/model_test.go",
				PackageName:     "./internal/tip",
				TestNamePattern: "TestTarget",
			},
		},
	}

	got := FilterHistoriesByPackages(histories, []string{"internal/parse"})

	if got == histories {
		t.Fatal("filtered histories points to original histories")
	}
	if got.ProjectDir != histories.ProjectDir {
		t.Fatalf("filtered histories project dir = %q, want %q", got.ProjectDir, histories.ProjectDir)
	}
	if len(got.Histories) != 1 {
		t.Fatalf("filtered histories len = %d, want 1", len(got.Histories))
	}
	if got.Histories[0].TestNamePattern != "TestParse" {
		t.Fatalf("filtered history test name = %q, want %q", got.Histories[0].TestNamePattern, "TestParse")
	}
}

func TestFilterHistoriesByPackages_includesSubpackages(t *testing.T) {
	histories := &Histories{
		Histories: []*History{
			{
				Path:            "./internal/parse/parse_test.go",
				PackageName:     "./internal/parse",
				TestNamePattern: "TestParse",
			},
			{
				Path:            "./internalx/x_test.go",
				PackageName:     "./internalx",
				TestNamePattern: "TestX",
			},
		},
	}

	got := FilterHistoriesByPackages(histories, []string{"./internal"})

	if len(got.Histories) != 1 {
		t.Fatalf("filtered histories len = %d, want 1", len(got.Histories))
	}
	if got.Histories[0].TestNamePattern != "TestParse" {
		t.Fatalf("filtered history test name = %q, want %q", got.Histories[0].TestNamePattern, "TestParse")
	}
}

func TestFilterHistoriesByPackages_supportsGoStylePackagePattern(t *testing.T) {
	histories := &Histories{
		Histories: []*History{
			{
				Path:            "./internal/parse/parse_test.go",
				PackageName:     "./internal/parse",
				TestNamePattern: "TestParse",
			},
			{
				Path:            "./cmd/gotip/main_test.go",
				PackageName:     "./cmd/gotip",
				TestNamePattern: "TestMain",
			},
		},
	}

	got := FilterHistoriesByPackages(histories, []string{"./internal/..."})

	if len(got.Histories) != 1 {
		t.Fatalf("filtered histories len = %d, want 1", len(got.Histories))
	}
	if got.Histories[0].TestNamePattern != "TestParse" {
		t.Fatalf("filtered history test name = %q, want %q", got.Histories[0].TestNamePattern, "TestParse")
	}
}

func TestFilterHistoriesByPackages_usesPathWhenPackageNameIsEmpty(t *testing.T) {
	histories := &Histories{
		Histories: []*History{
			{
				Path:            "./internal/parse/parse_test.go",
				TestNamePattern: "TestParse",
			},
		},
	}

	got := FilterHistoriesByPackages(histories, []string{"./internal/parse"})

	if len(got.Histories) != 1 {
		t.Fatalf("filtered histories len = %d, want 1", len(got.Histories))
	}
}

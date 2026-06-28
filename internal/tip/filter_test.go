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

func TestFilterTestsByPackages_requiresExactMatch(t *testing.T) {
	tests := map[string][]*TestFunction{
		"./internal/parse/parse_test.go": {
			{Name: "TestParse"},
		},
	}

	got := FilterTestsByPackages(tests, []string{"./internal"})

	if len(got) != 0 {
		t.Fatalf("filtered tests len = %d, want 0", len(got))
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

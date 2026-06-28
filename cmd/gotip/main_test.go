package main

import "testing"

func TestParseArgs_list(t *testing.T) {
	got, err := parseArgs([]string{"gotip", "list", "--format=json", "--package=./internal/parse", "--skip-subtests"})
	if err != nil {
		t.Fatalf("parseArgs() error = %v", err)
	}
	if got.Command != "list" {
		t.Errorf("command = %q, want %q", got.Command, "list")
	}
	if got.ListOptions.Format != "json" {
		t.Errorf("format = %q, want %q", got.ListOptions.Format, "json")
	}
	if !got.ListOptions.SkipSubtests {
		t.Error("list skip-subtests = false, want true")
	}
	wantPackages := []string{"./internal/parse"}
	if len(got.ListOptions.Packages) != len(wantPackages) {
		t.Fatalf("list packages len = %d, want %d", len(got.ListOptions.Packages), len(wantPackages))
	}
	for i := range wantPackages {
		if got.ListOptions.Packages[i] != wantPackages[i] {
			t.Errorf("list package %d = %q, want %q", i, got.ListOptions.Packages[i], wantPackages[i])
		}
	}
}

func TestParseArgs_listMultiplePackages(t *testing.T) {
	got, err := parseArgs([]string{"gotip", "list", "-p", "internal/parse", "-p", "internal/tip"})
	if err != nil {
		t.Fatalf("parseArgs() error = %v", err)
	}
	wantPackages := []string{"internal/parse", "internal/tip"}
	if len(got.ListOptions.Packages) != len(wantPackages) {
		t.Fatalf("list packages len = %d, want %d", len(got.ListOptions.Packages), len(wantPackages))
	}
	for i := range wantPackages {
		if got.ListOptions.Packages[i] != wantPackages[i] {
			t.Errorf("list package %d = %q, want %q", i, got.ListOptions.Packages[i], wantPackages[i])
		}
	}
}

func TestParseArgs_defaultExecutionArgs(t *testing.T) {
	got, err := parseArgs([]string{"gotip", "--package", "./internal/parse", "--skip-subtests", "--", "-v", "-count=1"})
	if err != nil {
		t.Fatalf("parseArgs() error = %v", err)
	}
	if got.Command != "" {
		t.Errorf("command = %q, want empty", got.Command)
	}
	if !got.Options.SkipSubtests {
		t.Error("root skip-subtests = false, want true")
	}
	wantPackages := []string{"./internal/parse"}
	if len(got.Options.Packages) != len(wantPackages) {
		t.Fatalf("root packages len = %d, want %d", len(got.Options.Packages), len(wantPackages))
	}
	for i := range wantPackages {
		if got.Options.Packages[i] != wantPackages[i] {
			t.Errorf("root package %d = %q, want %q", i, got.Options.Packages[i], wantPackages[i])
		}
	}
	wantArgs := []string{"-v", "-count=1"}
	if len(got.TestArgs) != len(wantArgs) {
		t.Fatalf("test args len = %d, want %d", len(got.TestArgs), len(wantArgs))
	}
	for i := range wantArgs {
		if got.TestArgs[i] != wantArgs[i] {
			t.Errorf("test arg %d = %q, want %q", i, got.TestArgs[i], wantArgs[i])
		}
	}
}

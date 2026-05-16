package main

import "testing"

func TestParseArgs_list(t *testing.T) {
	got, err := parseArgs([]string{"gotip", "list", "--format=json", "--skip-subtests"})
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
}

func TestParseArgs_defaultExecutionArgs(t *testing.T) {
	got, err := parseArgs([]string{"gotip", "--skip-subtests", "--", "-v", "-count=1"})
	if err != nil {
		t.Fatalf("parseArgs() error = %v", err)
	}
	if got.Command != "" {
		t.Errorf("command = %q, want empty", got.Command)
	}
	if !got.Options.SkipSubtests {
		t.Error("root skip-subtests = false, want true")
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

package main

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"slices"
	"strings"

	"github.com/lusingander/gotip/internal/parse"
	"github.com/lusingander/gotip/internal/ui"
)

var ignore = []string{
	"vendor",
	"testdata",
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	rootDir := "."
	tests := make(map[string][]*parse.TestFunction)

	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() && slices.Contains(ignore, d.Name()) {
			return filepath.SkipDir
		}
		if err != nil || !strings.HasSuffix(path, "_test.go") {
			return nil
		}
		testFunctions, err := parse.ProcessFile(path)
		if err != nil {
			return fmt.Errorf("error processing file %s: %w", path, err)
		}
		tests[path] = testFunctions
		return nil
	})
	if err != nil {
		return err
	}

	target, err := ui.Start(tests)
	if err != nil {
		return err
	}
	if target != nil {
		// todo
		fmt.Printf("Selected test: %s %s (unresolved: %t)\n", target.Path, target.Name, target.IsUnresolved)
	}
	return nil
}

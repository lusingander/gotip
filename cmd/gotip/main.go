package main

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"strings"

	"github.com/lusingander/gotip/internal/parse"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	rootDir := "./testdata"
	tests := make(map[string][]*parse.TestFunction)

	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
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

	parse.PrintTestFunctions(tests)

	return nil
}

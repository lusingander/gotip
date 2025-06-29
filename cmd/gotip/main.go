package main

import (
	"fmt"
	"log"

	"github.com/lusingander/gotip/internal/parse"
	"github.com/lusingander/gotip/internal/ui"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	tests, err := parse.ProcessFilesRecursively(".")
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

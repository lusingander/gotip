package main

import (
	"log"

	"github.com/lusingander/gotip/internal/command"
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
	if target == nil {
		return nil
	}
	return command.Test(target)
}

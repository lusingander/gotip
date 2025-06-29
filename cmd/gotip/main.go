package main

import (
	"log"
	"os"
	"slices"

	"github.com/lusingander/gotip/internal/command"
	"github.com/lusingander/gotip/internal/parse"
	"github.com/lusingander/gotip/internal/ui"
)

func main() {
	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func parseArgs(args []string) ([]string, []string) {
	i := slices.Index(args, "--")
	if i == -1 {
		return args, nil
	}
	return args[:i], args[i+1:]
}

func run(args []string) error {
	_, testArgs := parseArgs(args)

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

	return command.Test(target, testArgs)
}

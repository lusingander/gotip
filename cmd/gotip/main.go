package main

import (
	"log"
	"os"
	"slices"

	"github.com/lusingander/gotip/internal/command"
	"github.com/lusingander/gotip/internal/parse"
	"github.com/lusingander/gotip/internal/tip"
	"github.com/lusingander/gotip/internal/ui"
)

func main() {
	code, err := run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(code)
}

func parseArgs(args []string) ([]string, []string) {
	i := slices.Index(args, "--")
	if i == -1 {
		return args, nil
	}
	return args[:i], args[i+1:]
}

func run(args []string) (int, error) {
	_, testArgs := parseArgs(args)
	conf, err := tip.LoadConfig()
	if err != nil {
		return 1, err
	}

	tests, err := parse.ProcessFilesRecursively(".")
	if err != nil {
		return 1, err
	}

	target, err := ui.Start(tests)
	if err != nil {
		return 1, err
	}
	if target == nil {
		return 0, nil
	}

	return command.Test(target, testArgs, conf)
}

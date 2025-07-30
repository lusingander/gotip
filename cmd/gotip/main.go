package main

import (
	"fmt"
	"os"
	"slices"

	"github.com/jessevdk/go-flags"
	"github.com/lusingander/gotip/internal/command"
	"github.com/lusingander/gotip/internal/parse"
	"github.com/lusingander/gotip/internal/tip"
	"github.com/lusingander/gotip/internal/ui"
)

type options struct {
	View    string `short:"v" long:"view" description:"Default view" choice:"all" choice:"history" default:"all"`
	Filter  string `short:"f" long:"filter" description:"Default filter type" choice:"fuzzy" choice:"exact" default:"fuzzy"`
	Version bool   `short:"V" long:"version" description:"Print version"`
}

func main() {
	code, err := run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
	os.Exit(code)
}

func parseArgs(args []string) (*options, []string, error) {
	var cliArgs, testArgs []string
	if i := slices.Index(args, "--"); i != -1 {
		cliArgs = args[:i]
		testArgs = args[i+1:]
	} else {
		cliArgs = args
		testArgs = nil
	}
	var opts options
	if _, err := flags.ParseArgs(&opts, cliArgs); err != nil {
		return nil, nil, err
	}
	return &opts, testArgs, nil
}

func run(args []string) (int, error) {
	opt, testArgs, err := parseArgs(args)
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			return 0, nil
		}
		return 1, nil
	}
	if opt.Version {
		fmt.Fprintf(os.Stderr, "gotip %s\n", tip.AppVersion)
		return 0, nil
	}
	conf, err := tip.LoadConfig()
	if err != nil {
		return 1, err
	}
	histories, err := tip.LoadHistories(".")
	if err != nil {
		return 1, err
	}

	tests, err := parse.ProcessFilesRecursively(".")
	if err != nil {
		return 1, err
	}

	target, err := ui.Start(tests, histories, conf, opt.View, opt.Filter)
	if err != nil {
		return 1, err
	}
	if target == nil {
		return 0, nil
	}

	code, err := command.Test(target, testArgs, conf)
	if err != nil {
		return 1, err
	}

	histories.Add(target, conf.History.Limit)
	if err := tip.SaveHistories(".", histories); err != nil {
		return 1, err
	}

	return code, nil
}

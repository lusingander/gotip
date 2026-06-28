package main

import (
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/jessevdk/go-flags"
	"github.com/lusingander/gotip/internal/command"
	"github.com/lusingander/gotip/internal/listfmt"
	"github.com/lusingander/gotip/internal/parse"
	"github.com/lusingander/gotip/internal/tip"
	"github.com/lusingander/gotip/internal/ui"
)

type options struct {
	View         string   `short:"v" long:"view" description:"Default view" choice:"all" choice:"history" default:"all"`
	Filter       string   `short:"f" long:"filter" description:"Default filter type" choice:"fuzzy" choice:"exact" default:"fuzzy"`
	Packages     []string `short:"p" long:"package" value-name:"PACKAGE" description:"Filter by package name"`
	SkipSubtests bool     `short:"s" long:"skip-subtests" description:"Skip subtest detection"`
	Rerun        bool     `short:"r" long:"rerun" description:"Rerun the last test without showing the UI"`
	Version      bool     `short:"V" long:"version" description:"Print version"`
}

type listOptions struct {
	Packages     []string `short:"p" long:"package" value-name:"PACKAGE" description:"Filter by package name"`
	SkipSubtests bool     `short:"s" long:"skip-subtests" description:"Skip subtest detection"`
	Format       string   `long:"format" description:"Output format" choice:"text" choice:"json" default:"text"`
}

type parsedArgs struct {
	Options     *options
	ListOptions *listOptions
	Command     string
	TestArgs    []string
}

func main() {
	code, err := run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
	os.Exit(code)
}

func parseArgs(args []string) (*parsedArgs, error) {
	if len(args) > 0 {
		args = args[1:]
	}
	var cliArgs, testArgs []string
	if i := slices.Index(args, "--"); i != -1 {
		cliArgs = args[:i]
		testArgs = args[i+1:]
	} else {
		cliArgs = args
		testArgs = nil
	}
	var opts options
	var listOpts listOptions
	parser := flags.NewNamedParser("gotip", flags.Default)
	if _, err := parser.AddGroup("Application Options", "", &opts); err != nil {
		return nil, err
	}
	if _, err := parser.AddCommand("list", "List discovered tests", "List discovered tests without launching the UI", &listOpts); err != nil {
		return nil, err
	}
	parser.SubcommandsOptional = true
	if _, err := parser.ParseArgs(cliArgs); err != nil {
		return nil, err
	}
	command := ""
	if parser.Active != nil {
		command = parser.Active.Name
	}
	return &parsedArgs{
		Options:     &opts,
		ListOptions: &listOpts,
		Command:     command,
		TestArgs:    testArgs,
	}, nil
}

func run(args []string) (int, error) {
	parsed, err := parseArgs(args)
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			return 0, nil
		}
		return 1, nil
	}
	opt := parsed.Options

	if opt.Version {
		fmt.Fprintf(os.Stderr, "gotip %s\n", tip.AppVersion)
		return 0, nil
	}

	conf, err := tip.LoadConfig(".")
	if err != nil {
		return 1, err
	}

	if parsed.Command == "list" {
		if len(parsed.TestArgs) > 0 {
			return 1, errors.New("list does not accept test arguments after --")
		}
		skipSubtests := opt.SkipSubtests || parsed.ListOptions.SkipSubtests
		tests, err := parse.ProcessFilesRecursively(".", conf.Ignore, skipSubtests)
		if err != nil {
			return 1, err
		}
		packages := append([]string{}, opt.Packages...)
		packages = append(packages, parsed.ListOptions.Packages...)
		tests = tip.FilterTestsByPackages(tests, packages)
		switch parsed.ListOptions.Format {
		case "text":
			if err := listfmt.WriteText(os.Stdout, tests); err != nil {
				return 1, err
			}
		case "json":
			if err := listfmt.WriteJSON(os.Stdout, tests); err != nil {
				return 1, err
			}
		}
		return 0, nil
	}

	histories, err := tip.LoadHistories(".")
	if err != nil {
		return 1, err
	}

	if opt.Rerun {
		if len(histories.Histories) == 0 {
			fmt.Fprintln(os.Stderr, "No test history found.")
			return 1, nil
		}
		code, err := command.Test(histories.Histories[0].ToTarget(), parsed.TestArgs, conf)
		if err != nil {
			return 1, err
		}
		return code, nil
	}

	tests, err := parse.ProcessFilesRecursively(".", conf.Ignore, opt.SkipSubtests)
	if err != nil {
		return 1, err
	}
	tests = tip.FilterTestsByPackages(tests, opt.Packages)

	target, err := ui.Start(tests, histories, conf, opt.View, opt.Filter)
	if err != nil {
		return 1, err
	}
	if target == nil {
		return 0, nil
	}

	code, err := command.Test(target, parsed.TestArgs, conf)
	if err != nil {
		return 1, err
	}

	histories.Add(target, conf.History.Limit)
	if err := tip.SaveHistories(".", histories); err != nil {
		return 1, err
	}

	return code, nil
}

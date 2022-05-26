package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/IBM/gauge/cmd/gauge/cli"
	"github.com/peterbourgon/ff/v3/ffcli"
)

var (
	rootFlagSet = flag.NewFlagSet("gauge", flag.ExitOnError)
)

func main() {
	root := &ffcli.Command{
		ShortUsage: "gauge [flags] <subcommand>",
		FlagSet:    rootFlagSet,
		Subcommands: []*ffcli.Command{
			cli.Package(),
			cli.SBOM(),
			cli.Version()},
		Exec: func(context.Context, []string) error {
			return flag.ErrHelp
		},
	}

	if err := root.Parse(os.Args[1:]); err != nil {
		printErrAndExit(err)
	}

	if err := root.Run(context.Background()); err != nil {
		printErrAndExit(err)
	}
}

func printErrAndExit(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}

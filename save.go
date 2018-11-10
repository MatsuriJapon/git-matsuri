package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/subcommands"
	"strconv"
)

type saveCmd struct{}

func (*saveCmd) Name() string     { return "save" }
func (*saveCmd) Synopsis() string { return "save current work on GitHub" }
func (*saveCmd) Usage() string {
	return `save <ISSUE>:
	Save current work on GitHub in the correct branch
	`
}
func (p *saveCmd) SetFlags(_ *flag.FlagSet) {}
func (p *saveCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if len(f.Args()) == 1 {
		issue, err := strconv.Atoi(f.Args()[0])
		if err != nil {
			f.Usage()
			return subcommands.ExitUsageError
		}
		fmt.Printf("git matsuri save %v\n", issue)
		return subcommands.ExitSuccess
	}
	f.Usage()
	return subcommands.ExitUsageError
}

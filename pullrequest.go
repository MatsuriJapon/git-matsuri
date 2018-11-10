package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/subcommands"
	"strconv"
)

type prCmd struct {
	noclose bool
}

func (*prCmd) Name() string     { return "pr" }
func (*prCmd) Synopsis() string { return "open a pull request for ISSUE" }
func (*prCmd) Usage() string {
	return `pr [-noclose] <ISSUE>:
	Open a pull request for ISSUE, adding a mention to $ISSUE in the message to link the PR to the issue. Add '-noclose' to override the closing of the issue.
	`
}
func (p *prCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&p.noclose, "noclose", false, "do not close issue on merge")
}
func (p *prCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if len(f.Args()) == 1 {
		issue, err := strconv.Atoi(f.Args()[0])
		if err != nil {
			f.Usage()
			return subcommands.ExitUsageError
		}
		if p.noclose {
			fmt.Printf("git matsuri pr -noclose %v\n", issue)
		} else {
			fmt.Printf("git matsuri pr %v\n", issue)
		}
		return subcommands.ExitSuccess
	}
	f.Usage()
	return subcommands.ExitUsageError
}

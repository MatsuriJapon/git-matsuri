package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/subcommands"
	"strconv"
)

type fixCmd struct {
	noclose bool
}

func (*fixCmd) Name() string     { return "fix" }
func (*fixCmd) Synopsis() string { return "open a new PR to fix a bug in the original one" }
func (*fixCmd) Usage() string {
	return `fix [-noclose] <ISSUE>:
	Open a new PR to fix the original one. Add '-noclose' to override the closing of the issue.
	`
}

func (p *fixCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&p.noclose, "noclose", false, "do not close issue on merge")
}
func (p *fixCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if len(f.Args()) != 1 {
		f.Usage()
		return subcommands.ExitUsageError
	}
	issueNum, err := strconv.Atoi(f.Args()[0])
	if err != nil || !IsValidIssue(ctx, issueNum) {
		f.Usage()
		return subcommands.ExitUsageError
	}
	fmt.Printf("Creating a fix PR for ISSUE-%d...\n", issueNum)
	pr, err := CreateFixPRForIssueNumber(ctx, issueNum, p.noclose)
	if pr != nil {
		fmt.Printf("Pull Request created: %s\n", pr.GetHTMLURL())
	}
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

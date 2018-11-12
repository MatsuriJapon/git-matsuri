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
func (p *prCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if len(f.Args()) != 1 {
		f.Usage()
		return subcommands.ExitUsageError
	}
	issueNum, err := strconv.Atoi(f.Args()[0])
	if err != nil || !IsValidIssue(ctx, issueNum) {
		f.Usage()
		return subcommands.ExitUsageError
	}
	pr, cardErr := CreatePRForIssueNumber(ctx, issueNum, p.noclose)
	// we might succeed at creating the PR but fail at placing it in the To Do column
	if pr != nil {
		fmt.Printf("Pull Request created: %s\n", pr.GetHTMLURL())
	}
	if cardErr != nil {
		fmt.Println(cardErr)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

package matsuri

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/subcommands"
	"strconv"
)

// FixCmd is a git-matsuri subcommand
type FixCmd struct {
	noclose bool
}

// Name returns the subcommand name
func (*FixCmd) Name() string { return "fix" }

// Synopsis returns the subcommand synopsis
func (*FixCmd) Synopsis() string { return "open a new PR to fix a bug in the original one" }

// Usage returns the subcommand usage
func (*FixCmd) Usage() string {
	return `fix [-noclose] <ISSUE>:
	Open a new PR to fix the original one. Add '-noclose' to override the closing of the issue.
	`
}

// SetFlags sets the subcommand flags
func (p *FixCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&p.noclose, "noclose", false, "do not close issue on merge")
}

// Execute runs the subcommand
func (p *FixCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if len(f.Args()) != 1 {
		f.Usage()
		return subcommands.ExitUsageError
	}
	issueNum, err := strconv.Atoi(f.Args()[0])
	if err != nil || !IsExistingIssue(ctx, issueNum) {
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

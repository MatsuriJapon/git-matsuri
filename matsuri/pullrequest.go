package matsuri

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/subcommands"
	"os/exec"
	"strconv"
)

// PrCmd is a git-matsuri subcommand
type PrCmd struct {
	noclose bool
}

// Name returns the subcommand name
func (*PrCmd) Name() string { return "pr" }

// Synopsis returns the subcommand synopsis
func (*PrCmd) Synopsis() string { return "open a pull request for ISSUE" }

// Usage returns the subcommand usage
func (*PrCmd) Usage() string {
	return `pr [-noclose] <ISSUE>:
	Open a pull request for ISSUE, adding a mention to $ISSUE in the message to link the PR to the issue. Add '-noclose' to override the closing of the issue.
	`
}

// SetFlags sets the subcommand flags
func (p *PrCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&p.noclose, "noclose", false, "do not close issue on merge")
}

// Execute runs the subcommand
func (p *PrCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if len(f.Args()) != 1 {
		f.Usage()
		return subcommands.ExitUsageError
	}
	issueNum, err := strconv.Atoi(f.Args()[0])
	if err != nil || !IsValidIssue(ctx, issueNum) {
		f.Usage()
		return subcommands.ExitUsageError
	}
	cmd := exec.Command("git", "matsuri", "save", f.Args()[0])
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}
	fmt.Println(string(out))

	fmt.Printf("Creating a PR for ISSUE-%d...\n", issueNum)
	pr, err := CreatePRForIssueNumber(ctx, issueNum, p.noclose)
	// we might succeed at creating the PR but fail at placing it in the To Do column
	if pr != nil {
		fmt.Printf("Pull Request created: %s\n", pr.GetHTMLURL())
	}
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

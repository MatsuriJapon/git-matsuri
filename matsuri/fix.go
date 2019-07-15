package matsuri

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/subcommands"
	"os/exec"
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
	cmd := exec.Command("git", "matsuri", "save", f.Args()[0])
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}
	fmt.Println(string(out))

	fmt.Printf("Creating a fix PR for ISSUE-%d...\n", issueNum)
	pr, err := CreateFixPRForIssueNumber(ctx, issueNum, p.noclose)
	if pr != nil {
		fmt.Printf("Pull Request created: %s\n", pr.GetHTMLURL())
	}
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}
	// reopen Issue if it has been closed
	err = ReopenIssue(ctx, issueNum)
	if err != nil {
		fmt.Printf("Created PR but could not reopen Issue #%d\n", issueNum)
	}
	return subcommands.ExitSuccess
}

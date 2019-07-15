package matsuri

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/subcommands"
	"os/exec"
	"strconv"
)

// SaveCmd is a git-matsuri subcommand
type SaveCmd struct{}

// Name returns the subcommand name
func (*SaveCmd) Name() string { return "save" }

// Synopsis returns the subcommand synopsis
func (*SaveCmd) Synopsis() string { return "save current work on GitHub" }

// Usage returns the subcommand usage
func (*SaveCmd) Usage() string {
	return `save <ISSUE>:
	Save current work on GitHub in the correct branch
	`
}

// SetFlags sets the subcommand flags
func (p *SaveCmd) SetFlags(_ *flag.FlagSet) {}

// Execute runs the subcommand
func (p *SaveCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if len(f.Args()) != 1 {
		f.Usage()
		return subcommands.ExitUsageError
	}
	issue, err := strconv.Atoi(f.Args()[0])
	if err != nil || !IsExistingIssue(ctx, issue) {
		f.Usage()
		return subcommands.ExitUsageError
	}
	fmt.Println("Pushing your changes to GitHub...")
	branches := fmt.Sprintf("ISSUE-%d:ISSUE-%d", issue, issue)
	cmd := exec.Command("git", "push", "-u", "origin", branches)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("There was an issue pushing the branch")
		return subcommands.ExitFailure
	}
	fmt.Println(string(out))
	return subcommands.ExitSuccess
}

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/subcommands"
	"os/exec"
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
func (p *saveCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if len(f.Args()) != 1 {
		f.Usage()
		return subcommands.ExitUsageError
	}
	issue, err := strconv.Atoi(f.Args()[0])
	if err != nil || !IsValidIssue(ctx, issue) {
		f.Usage()
		return subcommands.ExitUsageError
	}
	branches := fmt.Sprintf("ISSUE-%d:ISSUE-%d", issue, issue)
	cmd := exec.Command("git", "push", "-u", "origin", branches)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("There was an issue pushing the branch")
		return subcommands.ExitFailure
	}
	fmt.Print(string(out))
	return subcommands.ExitSuccess
}

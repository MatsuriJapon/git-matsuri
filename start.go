package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/subcommands"
	"os/exec"
	"strconv"
)

type startCmd struct{}

func (*startCmd) Name() string     { return "start" }
func (*startCmd) Synopsis() string { return "start working on an open issue" }
func (*startCmd) Usage() string {
	return `start [<ISSUE>]:
	Create a named branch for the chosen issue and change its status to 'Doing'
	`
}
func (p *startCmd) SetFlags(_ *flag.FlagSet) {}
func (p *startCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	var issue int
	// weird logic to get "current" Matsuri project year
	currentYear := GetCurrentProjectYear()
	if len(f.Args()) == 1 {
		issueNumber, err := strconv.Atoi(f.Args()[0])
		if err != nil {
			f.Usage()
			return subcommands.ExitUsageError
		}
		if !IsValidIssue(ctx, issueNumber) {
			fmt.Println("Error: an invalid Issue was provided")
			return subcommands.ExitFailure
		}
		issue = issueNumber
	} else {
		issues, err := GetIssuesForProject(ctx, currentYear)
		if err != nil {
			fmt.Println(err)
			return subcommands.ExitFailure
		}
		PrintIssues(issues)
		fmt.Print("\n\nEnter an Issue number: ")
		var input string
		_, err2 := fmt.Scanln(&input)
		issueNumber, err3 := strconv.Atoi(input)
		if err2 != nil || err3 != nil || !IsValidIssue(ctx, issueNumber) {
			fmt.Println("Error: an invalid Issue was provided")
			return subcommands.ExitFailure
		}
		issue = issueNumber
	}
	if IsMainRepo() {
		// move project to Doing, or fail
		moveErr := MoveProjectCardForProject(ctx, issue, currentYear)
		if moveErr != nil {
			fmt.Println(moveErr)
			return subcommands.ExitFailure
		}
	}
	// checkout branch
	branchName := fmt.Sprintf("ISSUE-%d", issue)
	cmd := exec.Command("git", "checkout", "-b", branchName)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("There was an issue creating the git branch")
		return subcommands.ExitFailure
	}
	fmt.Print(string(out))
	return subcommands.ExitSuccess
}

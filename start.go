package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/subcommands"
	"os/exec"
	"regexp"
	"strconv"
)

type startCmd struct{}

func prepareCheckout(ctx context.Context) (err error) {
	// status
	fmt.Println("Checking status of current branch...")
	cmd := exec.Command("git", "status")
	out, err := cmd.Output()
	fmt.Println(string(out))
	if err != nil {
		return
	}
	r := regexp.MustCompile("nothing to commit")
	match := r.Match(out)
	if !match {
		err = fmt.Errorf("Error: there might be unsaved changes in the current repository.\nResolve them before creating a new branch")
		return
	}

	// checkout default
	fmt.Println("Checking out default branch...")
	defaultBranch, _ := GetDefaultBranch(ctx)
	cmd = exec.Command("git", "checkout", *defaultBranch)
	out, err = cmd.Output()
	fmt.Println(string(out))
	if err != nil {
		return
	}

	// pull
	fmt.Println("Pulling changes...")
	cmd = exec.Command("git", "pull")
	out, err = cmd.Output()
	fmt.Println(string(out))
	return
}

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
	// pre-checkout checks
	if err := prepareCheckout(ctx); err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	if IsMainRepo() {
		// move project to Doing, or fail
		if err := MoveProjectCardForProject(ctx, issue, currentYear); err != nil {
			fmt.Println(err)
			return subcommands.ExitFailure
		}
	}
	// checkout branch
	fmt.Println("Checking out topic branch...")
	branchName := fmt.Sprintf("ISSUE-%d", issue)
	cmd := exec.Command("git", "checkout", "-b", branchName)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("There was an issue creating the git branch")
		return subcommands.ExitFailure
	}
	fmt.Println(string(out))
	fmt.Printf("You are now working in branch %s\n", branchName)
	return subcommands.ExitSuccess
}

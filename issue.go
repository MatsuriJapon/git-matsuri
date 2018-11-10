package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/subcommands"
	"strconv"
)

type issueCmd struct{}

func (*issueCmd) Name() string     { return "issues" }
func (*issueCmd) Synopsis() string { return "list all opened issues" }
func (*issueCmd) Usage() string {
	return `issues [<YEAR>]:
	List all opened issues. If YEAR is provided, only show those for the current year.
	`
}
func (p *issueCmd) SetFlags(_ *flag.FlagSet) {}
func (p *issueCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if len(f.Args()) == 1 {
		year, err := strconv.Atoi(f.Args()[0])
		if err != nil {
			f.Usage()
			return subcommands.ExitUsageError
		}
		issues, err := GetIssuesForProject(ctx, year)
		if err != nil {
			fmt.Println(err)
			return subcommands.ExitFailure
		}
		PrintIssues(issues)
	} else {
		issues, _, _ := client.Issues.ListByRepo(ctx, owner, repo, nil)
		PrintIssues(issues)
	}
	return subcommands.ExitSuccess
}

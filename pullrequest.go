package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/go-github/v18/github"
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
	issue, _, _ := client.Issues.Get(ctx, owner, repo, issueNum)
	title := fmt.Sprintf("ISSUE-%d: %s", issue.GetNumber(), issue.GetTitle())
	head := fmt.Sprintf("ISSUE-%d", issue.GetNumber())
	base := fmt.Sprintf("v%d", GetCurrentProjectYear())
	body := fmt.Sprintf("Closes #%d\n", issue.GetNumber())
	if p.noclose {
		body = fmt.Sprintf("Related to #%d\n", issue.GetNumber())
	}
	newPr := &github.NewPullRequest{
		Title: github.String(title),
		Head:  github.String(head),
		Base:  github.String(base),
		Body:  github.String(body),
	}
	pr, _, prErr := client.PullRequests.Create(ctx, owner, repo, newPr)
	if prErr != nil {
		fmt.Println(prErr)
		return subcommands.ExitFailure
	}
	fmt.Printf("Pull Request created: %s\n", pr.GetHTMLURL())
	return subcommands.ExitSuccess
}

package cmd

import (
	"github.com/MatsuriJapon/git-matsuri/matsuri"
	"github.com/spf13/cobra"
)

func completeIssues(_ *cobra.Command, args []string, toComplete string, issueGetter matsuri.IssueGetterFunc) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	openIssues, err := issueGetter()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return matsuri.GetIssueNumbersStartingWith(openIssues, toComplete), cobra.ShellCompDirectiveNoFileComp
}

func completeOpenIssuesForProject(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return completeIssues(cmd, args, toComplete, matsuri.GetOpenIssuesForProject)
}

func completeInProgressIssuesForProject(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return completeIssues(cmd, args, toComplete, matsuri.GetInProgressIssues)
}

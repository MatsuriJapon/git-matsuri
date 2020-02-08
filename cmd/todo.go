package cmd

import (
	"strconv"

	"github.com/MatsuriJapon/git-matsuri/matsuri"
	"github.com/spf13/cobra"
)

var (
	todoCmd = &cobra.Command{
		Use:   "todo",
		Short: "list opened issues",
		Long:  "If YEAR is provided, list Issues for the current year's Kanban. Otherwise, all open Issues are listed",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runTodo,
	}
)

func runTodo(cmd *cobra.Command, args []string) error {
	if len(args) == 1 {
		year, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		issues, err := matsuri.GetIssuesForProject(year)
		if err != nil {
			return err
		}
		matsuri.PrintIssues(issues)
	} else {
		issues, err := matsuri.GetRepoIssues()
		if err != nil {
			return err
		}
		matsuri.PrintIssues(issues)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(todoCmd)
}

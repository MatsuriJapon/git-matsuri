package cmd

import (
	"github.com/MatsuriJapon/git-matsuri/matsuri"
	"github.com/spf13/cobra"
)

var (
	todoCmd = &cobra.Command{
		Use:   "todo",
		Short: "list opened issues",
		RunE:  runTodo,
	}
	showOnlyCurrentRepo bool
)

func runTodo(cmd *cobra.Command, args []string) error {
	issues, err := matsuri.GetIssues(showOnlyCurrentRepo)
	if err != nil {
		return err
	}

	matsuri.PrintIssues(issues)
	return nil
}

func init() {
	todoCmd.Flags().BoolVarP(&showOnlyCurrentRepo, "current", "c", false, "show only issues for the current repo")
	rootCmd.AddCommand(todoCmd)
}

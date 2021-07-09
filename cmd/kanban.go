package cmd

import (
	"github.com/MatsuriJapon/git-matsuri/matsuri"
	"github.com/spf13/cobra"
)

var (
	kanbanCmd = &cobra.Command{
		Use:   "kanban",
		Short: "show the Kanban for the current year",
		RunE:  runKanban,
	}
)

func runKanban(cmd *cobra.Command, args []string) (err error) {
	project, err := matsuri.GetProject()
	if err != nil {
		return
	}
	matsuri.PrintProjectKanban(project)
	return
}

func init() {
	rootCmd.AddCommand(kanbanCmd)
}

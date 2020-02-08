package cmd

import (
	"github.com/MatsuriJapon/git-matsuri/matsuri"
	"github.com/spf13/cobra"
	"strconv"
)

var (
	kanbanCmd = &cobra.Command{
		Use:   "kanban",
		Short: "show the Kanban for the current year",
		Args:  cobra.ExactArgs(1),
		RunE:  runKanban,
	}
)

func runKanban(cmd *cobra.Command, args []string) (err error) {
	year, err := strconv.Atoi(args[0])
	if err != nil {
		return
	}
	project, err := matsuri.GetProjectForYear(year)
	if err != nil {
		return
	}
	matsuri.PrintProjectKanban(project)
	return
}

func init() {
	rootCmd.AddCommand(kanbanCmd)
}

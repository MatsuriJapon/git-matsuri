package cmd

import (
	"errors"
	"fmt"
	"github.com/MatsuriJapon/git-matsuri/matsuri"
	"github.com/spf13/cobra"
	"os/exec"
	"strconv"
)

var (
	saveCmd = &cobra.Command{
		Use:               "save",
		Short:             "save current work on GitHub",
		Args:              cobra.ExactArgs(1),
		RunE:              runSave,
		ValidArgsFunction: completeInProgressIssuesForProject,
	}
)

func runSave(cmd *cobra.Command, args []string) (err error) {
	issue, err := strconv.Atoi(args[0])
	if err != nil {
		return
	}
	if !matsuri.IsExistingIssue(issue) {
		err = errors.New("the provided Issue doesn't exist")
		return
	}
	cmd.Println("Pushing your changes to GitHub...")
	branches := fmt.Sprintf("ISSUE-%d:ISSUE-%d", issue, issue)
	pushCmd := exec.Command("git", "push", "-u", "origin", branches)
	out, err := pushCmd.Output()
	if err != nil {
		err = fmt.Errorf("there was a problem pushing the branch: %s", err.Error())
		return
	}
	cmd.Println(string(out))
	return
}

func init() {
	rootCmd.AddCommand(saveCmd)
}

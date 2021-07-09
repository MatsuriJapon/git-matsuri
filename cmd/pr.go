package cmd

import (
	"errors"
	"github.com/MatsuriJapon/git-matsuri/matsuri"
	"github.com/spf13/cobra"
	"os/exec"
	"strconv"
)

var (
	noCloseAfterPR bool
	prCmd          = &cobra.Command{
		Use:               "pr",
		Short:             "open a pull request for ISSUE",
		Long:              "Open a pull request for ISSUE, adding a mention to $ISSUE in the message to link the PR to the issue. Add '-noclose' to override the closing of the issue",
		Args:              cobra.ExactArgs(1),
		RunE:              runPR,
		ValidArgsFunction: completeInProgressIssuesForProject,
	}
)

func runPR(cmd *cobra.Command, args []string) (err error) {
	issueNum, err := strconv.Atoi(args[0])
	if err != nil {
		return
	}
	if !matsuri.IsValidIssue(issueNum) {
		err = errors.New("an invalid Issue number was provided")
		return
	}
	pushCmd := exec.Command("git", "matsuri", "save", args[0])
	out, err := pushCmd.Output()
	if err != nil {
		return
	}
	cmd.Println(string(out))
	cmd.Printf("Creating a PR for ISSUE-%d...\n", issueNum)
	pr, err := matsuri.CreatePRForIssueNumber(issueNum, noCloseAfterPR)
	// we might succeed at creating the PR but fail at placing it in the To Do column
	if pr != nil {
		cmd.Printf("Pull Request created: %s\n", pr.GetHTMLURL())
	}
	return
}

func init() {
	prCmd.Flags().BoolVar(&noCloseAfterPR, "noclose", false, "do not close Issue on merge")
	rootCmd.AddCommand(prCmd)
}

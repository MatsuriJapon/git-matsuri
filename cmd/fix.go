package cmd

import (
	"errors"
	"github.com/MatsuriJapon/git-matsuri/matsuri"
	"github.com/spf13/cobra"
	"os/exec"
	"strconv"
)

var (
	noCloseAfterFix bool
	fixCmd          = &cobra.Command{
		Use:   "fix",
		Short: "open a new PR to fix a bug in the original one",
		Long:  "Open a new PR to fix the original one. Add '-noclose' to override the closing of the issue",
		Args:  cobra.ExactArgs(1),
		RunE:  runFix,
	}
)

func runFix(cmd *cobra.Command, args []string) (err error) {
	issueNum, err := strconv.Atoi(args[0])
	if err != nil {
		return
	}
	if !matsuri.IsExistingIssue(issueNum) {
		err = errors.New("an invalid Issue was provided")
		return
	}
	pushCmd := exec.Command("git", "matsuri", "save", args[0])
	out, err := pushCmd.Output()
	if err != nil {
		return
	}
	cmd.Println(string(out))
	cmd.Printf("Creating a fix PR for ISSUE-%d...\n", issueNum)
	pr, err := matsuri.CreateFixPRForIssueNumber(issueNum, noCloseAfterFix)
	if pr != nil {
		cmd.Printf("Pull Request created: %s\n", pr.GetHTMLURL())
	}
	if err != nil {
		return
	}
	// reopen Issue if it has been closed
	err = matsuri.ReopenIssue(issueNum)
	return
}

func init() {
	fixCmd.Flags().BoolVar(&noCloseAfterFix, "noclose", false, "do not close Issue on merge")
	rootCmd.AddCommand(fixCmd)
}

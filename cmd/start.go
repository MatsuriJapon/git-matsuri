package cmd

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/MatsuriJapon/git-matsuri/matsuri"
	"github.com/spf13/cobra"
)

var (
	startCmd = &cobra.Command{
		Use:               "start ISSUE_NUMBER",
		Short:             "start working on an open issue",
		Args:              cobra.ExactArgs(1),
		RunE:              runStart,
		ValidArgsFunction: completeOpenIssuesForProject,
	}
)

func prepareCheckout(cmd *cobra.Command) (err error) {
	// status
	cmd.Println("Checking status of current branch...")
	statusCmd := exec.Command("git", "status")
	out, err := statusCmd.Output()
	cmd.Println(string(out))
	if err != nil {
		return
	}
	r := regexp.MustCompile("nothing to commit")
	match := r.Match(out)
	if !match {
		err = errors.New("there might be unsaved changes in the current repository.\nResolve them before creating a new branch")
		return
	}

	// checkout default
	cmd.Println("Checking out default branch...")
	defaultBranch, _ := matsuri.GetDefaultBranch()
	checkoutCmd := exec.Command("git", "checkout", *defaultBranch)
	out, err = checkoutCmd.Output()
	cmd.Println(string(out))
	if err != nil {
		return
	}

	// pull
	cmd.Println("Pulling changes...")
	pullCmd := exec.Command("git", "pull")
	out, err = pullCmd.Output()
	cmd.Println(string(out))
	return
}

func runStart(cmd *cobra.Command, args []string) (err error) {
	issueNumber, err := strconv.Atoi(args[0])
	if err != nil {
		return
	}

	if !matsuri.IsValidIssue(issueNumber) {
		err = errors.New("invalid Issue provided")
		return
	}
	err = prepareCheckout(cmd)
	if err != nil {
		return
	}
	// Some Issues may not be assigned to a Project, so we'll ignore errors here
	_ = matsuri.MoveProjectCardForProject(issueNumber)
	// checkout branch
	cmd.Println("Checking out topic branch...")
	branchName := fmt.Sprintf("ISSUE-%d", issueNumber)
	checkoutCmd := exec.Command("git", "checkout", "-b", branchName)
	out, err := checkoutCmd.Output()
	if err != nil {
		err = fmt.Errorf("there was an issue creating the git branch: %s", err.Error())
		return
	}
	cmd.Println(string(out))
	cmd.Printf("You are now working in branch %s\n", branchName)
	return
}

func init() {
	rootCmd.AddCommand(startCmd)
}

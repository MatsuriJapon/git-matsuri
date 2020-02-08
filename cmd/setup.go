package cmd

import (
	"errors"
	"github.com/MatsuriJapon/git-matsuri/matsuri"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

var (
	useHTTP  bool
	setupCmd = &cobra.Command{
		Use:   "setup",
		Short: "clones a Matsuri repository",
		Args:  cobra.ExactArgs(1),
		RunE:  runSetup,
	}
)

func runSetup(cmd *cobra.Command, args []string) (err error) {
	repoName := args[0]
	existingRepo, _ := matsuri.GetRepoName()
	if existingRepo == repoName {
		cmd.Println("you are already inside the target repository")
		return
	}
	if existingRepo != "" {
		err = errors.New("you are already inside a git repository, aborting")
		return
	}
	cloneURL, err := matsuri.GetRepoURL(repoName, useHTTP)
	if err != nil {
		return
	}
	cloneCmd := exec.Command("git", "clone", cloneURL)
	out, err := cloneCmd.CombinedOutput()
	if err != nil {
		return
	}
	cmd.Println(out)

	matsuriEmail, err := matsuri.GetMatsuriEmail()
	if err != nil {
		return
	}
	if matsuriEmail == "" {
		err = errors.New("a Matsuri email address was not found in this account, please add one in your GitHun profile")
		return
	}
	err = os.Chdir(repoName)
	if err != nil {
		return
	}
	configCmd := exec.Command("git", "config", "--local", "user.email", matsuriEmail)
	err = configCmd.Run()
	return
}

func init() {
	setupCmd.Flags().BoolVar(&useHTTP, "http", false, "clones the repository using the HTTP protocol")
	rootCmd.AddCommand(setupCmd)
}

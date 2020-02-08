package cmd

import (
	"errors"
	"github.com/MatsuriJapon/git-matsuri/matsuri"
	"github.com/spf13/cobra"
	"os"
)

var (
	// CurrentVersion is a build-time string representing the current version
	rootCmd = &cobra.Command{
		Use:               "git-matsuri",
		Short:             "git-matsuri provides useful git subcommands for matsuri workflows",
		PersistentPreRunE: sanity,
	}
)

func sanity(cmd *cobra.Command, args []string) (err error) {
	if _, err := matsuri.GetRepoName(); err != nil {
		cmd.Println("WARN: You are currently not in a git repository, some subcommands may not run.")
	}
	if token := os.Getenv(matsuri.TokenName); token == "" {
		err = errors.New("gitHub token not found.\nPlease create one at https://github.com/settings/tokens/new with 'repo' and 'user:email' permissions and save it to your system environment variables under the name MATSURI_TOKEN")
	}
	return
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		rootCmd.PrintErr(err)
		os.Exit(1)
	}
}

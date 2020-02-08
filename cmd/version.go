package cmd

import (
	"github.com/MatsuriJapon/git-matsuri/matsuri"
	"github.com/hashicorp/go-version"
	"github.com/spf13/cobra"
)

var (
	CurrentVersion string
	versionCmd     = &cobra.Command{
		Use:   "version",
		Short: "show the version number",
		Long:  "Show the version number for git-matsuri and checks if it currently is the latest version.string",
		Args:  cobra.NoArgs,
		RunE:  runVersion,
	}
)

func runVersion(cmd *cobra.Command, args []string) (err error) {
	currentVersion, err := version.NewVersion(CurrentVersion)
	if err != nil {
		return
	}
	latestVersion, err := matsuri.GetLatestVersion()
	if err != nil {
		return
	}
	cmd.Printf("git-matsuri version %s\n", currentVersion)
	if currentVersion.LessThan(latestVersion) {
		cmd.Printf("A new version is available: %s\n", latestVersion)
	} else {
		cmd.Println("You are using the latest version")
	}
	return
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

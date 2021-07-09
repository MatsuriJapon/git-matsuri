package cmd

import (
	"github.com/MatsuriJapon/git-matsuri/matsuri"
	"github.com/hashicorp/go-version"
	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "show the version number",
		Long:  "Show the version number for git-matsuri and checks if it currently is the latest version.string",
		Args:  cobra.NoArgs,
		RunE:  runVersion,
	}

	currentVersion string
	commit         string
	date           string
	builtBy        string
)

func runVersion(cmd *cobra.Command, args []string) (err error) {
	cv, err := version.NewVersion(currentVersion)
	if err != nil {
		return
	}
	lv, err := matsuri.GetLatestVersion()
	if err != nil {
		return
	}
	cmd.Printf("git-matsuri version %s\ncommit: %s\ndate: %s\nbuiltBy: %s\n\n", cv, commit, date, builtBy)
	if cv.LessThan(lv) {
		cmd.Printf("A new version is available: %s\n", lv)
	} else {
		cmd.Println("You are using the latest version")
	}
	return
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

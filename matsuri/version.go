package matsuri

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/subcommands"
	"github.com/hashicorp/go-version"
)

var CurrentVersion string

// VersionCmd is a git-matsuri subcommand
type VersionCmd struct{}

// Name returns the subcommand name
func (*VersionCmd) Name() string { return "version" }

// Synopsis returns the subcommand synopsis
func (*VersionCmd) Synopsis() string { return "show the version number" }

// Usage returns the subcommand usage
func (*VersionCmd) Usage() string {
	return `version:
	Show the version number for git-matsuri and checks if it currently is the latest version.string
	`
}

// SetFlags sets the subcommand flags
func (p *VersionCmd) SetFlags(_ *flag.FlagSet) {}

// Execute runs the subcommand
func (p *VersionCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if len(f.Args()) != 0 {
		f.Usage()
		return subcommands.ExitUsageError
	}
	currentVersion, err := version.NewVersion(CurrentVersion)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}
	latestVersion, err := GetLatestVersion(ctx)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}
	fmt.Printf("git-matsuri version %s\n", currentVersion)
	if currentVersion.LessThan(latestVersion) {
		fmt.Printf("A new version is available: %s\n", latestVersion)
	} else {
		fmt.Println("You are using the latest version")
	}
	return subcommands.ExitSuccess
}

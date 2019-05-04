package matsuri

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/subcommands"
	"os/exec"
)

// SetupCmd is a git-matsuri subcommand to clone a Matsuri repository
type SetupCmd struct {
	http bool
}

// Name returns the subcommand name
func (*SetupCmd) Name() string { return "setup" }

// Synopsis returns the subcommand synopsis
func (*SetupCmd) Synopsis() string { return "clones a Matsuri repository" }

// Usage returns the subcommand usage
func (*SetupCmd) Usage() string {
	return `setup [-http] <NAME>:
	Clones a MatsuriJapon repository having the provided NAME, if it exists. Cloning is done via SSH, although HTTP is provided as a fallback using the '-http' flag (not recommended).
	`
}

// SetFlags sets the subcommand flags
func (p *SetupCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&p.http, "http", false, "clones the repository using the HTTP protocol")
}

// Execute runs the subcommand
func (p *SetupCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if len(f.Args()) != 1 {
		f.Usage()
		return subcommands.ExitUsageError
	}
	repoName := f.Args()[0]
	// Check if we're not already in a git repository
	existingRepo, _ := GetRepoName()
	if existingRepo == repoName {
		fmt.Println("You are already inside the target repository.")
		return subcommands.ExitSuccess
	}
	if existingRepo != "" {
		fmt.Println("You are already inside a git repository. Aborting.")
		return subcommands.ExitFailure
	}

	cloneURL, err := GetRepoURL(ctx, repoName, p.http)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}
	cloneCmd := exec.Command("git", "clone", cloneURL)
	out, err := cloneCmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}
	fmt.Println(string(out))
	return subcommands.ExitSuccess
}
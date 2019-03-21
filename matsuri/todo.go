package matsuri

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/subcommands"
	"strconv"
)

// TodoCmd is a git-matsuri subcommand
type TodoCmd struct{}

// Name returns the subcommand name
func (*TodoCmd) Name() string { return "todo" }

// Synopsis returns the subcommand synopsis
func (*TodoCmd) Synopsis() string { return "list all opened issues" }

// Usage returns the subcommand usage
func (*TodoCmd) Usage() string {
	return `todo [<YEAR>]:
	List all opened issues. If YEAR is provided, only show those for the current year.
	`
}

// SetFlags sets the subcommand flags
func (p *TodoCmd) SetFlags(_ *flag.FlagSet) {}

// Execute runs the subcommand
func (p *TodoCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if len(f.Args()) == 1 {
		year, err := strconv.Atoi(f.Args()[0])
		if err != nil {
			f.Usage()
			return subcommands.ExitUsageError
		}
		issues, err := GetIssuesForProject(ctx, year)
		if err != nil {
			fmt.Println(err)
			return subcommands.ExitFailure
		}
		PrintIssues(issues)
	} else {
		issues, err := GetRepoIssues(ctx)
		if err != nil {
			fmt.Println(err)
			return subcommands.ExitFailure
		}
		PrintIssues(issues)
	}
	return subcommands.ExitSuccess
}

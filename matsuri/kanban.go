package matsuri

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/subcommands"
	"strconv"
)

// KanbanCmd is a git-matsuri subcommand
type KanbanCmd struct{}

// Name returns the subcommand name
func (*KanbanCmd) Name() string { return "kanban" }

// Synopsis returns the subcommand synopsis
func (*KanbanCmd) Synopsis() string { return "show the kanban, if any" }

// Usage returns the subcommand usage
func (*KanbanCmd) Usage() string {
	return `kanban [<YEAR>]:
	Show the kanban for the given YEAR if Projects is enabled in the given repository.
	`
}

// SetFlags sets the subcommand flags
func (p *KanbanCmd) SetFlags(_ *flag.FlagSet) {}

// Execute runs the subcommand
func (p *KanbanCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if len(f.Args()) != 1 {
		f.Usage()
		return subcommands.ExitUsageError
	}
	year, err := strconv.Atoi(f.Args()[0])
	if err != nil {
		f.Usage()
		return subcommands.ExitUsageError
	}
	project, err := GetProjectForYear(ctx, year)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}
	PrintProjectKanban(ctx, project)
	return subcommands.ExitSuccess
}

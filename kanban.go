package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/subcommands"
	"strconv"
)

type kanbanCmd struct{}

func (*kanbanCmd) Name() string     { return "kanban" }
func (*kanbanCmd) Synopsis() string { return "show the kanban, if any" }
func (*kanbanCmd) Usage() string {
	return `kanban [<YEAR>]:
	Show the kanban for the given YEAR if Projects is enabled in the given repository.
	`
}
func (p *kanbanCmd) SetFlags(_ *flag.FlagSet) {}
func (p *kanbanCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
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

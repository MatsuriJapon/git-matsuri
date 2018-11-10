package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/subcommands"
	"strconv"
)

type startCmd struct{}

func (*startCmd) Name() string     { return "do" }
func (*startCmd) Synopsis() string { return "start working on an open issue" }
func (*startCmd) Usage() string {
	return `do [<ISSUE>]:
	Create a named branch for the chosen issue and change its status to 'Doing'
	`
}
func (p *startCmd) SetFlags(_ *flag.FlagSet) {}
func (p *startCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if len(f.Args()) == 1 {
		issue, err := strconv.Atoi(f.Args()[0])
		if err != nil {
			f.Usage()
			return subcommands.ExitUsageError
		}
		fmt.Printf("git matsuri do %v\n", issue)
	} else {
		fmt.Println("git matsuri do")
	}
	return subcommands.ExitSuccess
}

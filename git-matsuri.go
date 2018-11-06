package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/google/subcommands"
)

type issueCmd struct{}

func (*issueCmd) Name() string     { return "issues" }
func (*issueCmd) Synopsis() string { return "list all opened issues" }
func (*issueCmd) Usage() string {
	return `issues [<YEAR>]:
	List all opened issues. If YEAR is provided, only show those for the current year.
	`
}
func (p *issueCmd) SetFlags(_ *flag.FlagSet) {}
func (p *issueCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if len(f.Args()) == 1 {
		year, err := strconv.Atoi(f.Args()[0])
		if err != nil {
			f.Usage()
			return subcommands.ExitUsageError
		}
		fmt.Printf("git matsuri issues %v\n", year)
	} else {
		fmt.Println("git matsuri issues")
	}
	return subcommands.ExitSuccess
}

type doCmd struct{}

func (*doCmd) Name() string     { return "do" }
func (*doCmd) Synopsis() string { return "start working on an open issue" }
func (*doCmd) Usage() string {
	return `do [<ISSUE>]:
	Create a named branch for the chosen issue and change its status to 'Doing'
	`
}
func (p *doCmd) SetFlags(_ *flag.FlagSet) {}
func (p *doCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
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

type saveCmd struct{}

func (*saveCmd) Name() string     { return "save" }
func (*saveCmd) Synopsis() string { return "save current work on GitHub" }
func (*saveCmd) Usage() string {
	return `save <ISSUE>:
	Save current work on GitHub in the correct branch
	`
}
func (p *saveCmd) SetFlags(_ *flag.FlagSet) {}
func (p *saveCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if len(f.Args()) == 1 {
		issue, err := strconv.Atoi(f.Args()[0])
		if err != nil {
			f.Usage()
			return subcommands.ExitUsageError
		}
		fmt.Printf("git matsuri save %v\n", issue)
		return subcommands.ExitSuccess
	}
	f.Usage()
	return subcommands.ExitUsageError
}

type prCmd struct {
	noclose bool
}

func (*prCmd) Name() string     { return "pr" }
func (*prCmd) Synopsis() string { return "open a pull request for ISSUE" }
func (*prCmd) Usage() string {
	return `pr [-noclose] <ISSUE>:
	Open a pull request for ISSUE, adding a mention to $ISSUE in the message to link the PR to the issue. Add '-noclose' to override the closing of the issue.
	`
}
func (p *prCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&p.noclose, "noclose", false, "do not close issue on merge")
}
func (p *prCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if len(f.Args()) == 1 {
		issue, err := strconv.Atoi(f.Args()[0])
		if err != nil {
			f.Usage()
			return subcommands.ExitUsageError
		}
		if p.noclose {
			fmt.Printf("git matsuri pr -noclose %v\n", issue)
		} else {
			fmt.Printf("git matsuri pr %v\n", issue)
		}
		return subcommands.ExitSuccess
	}
	f.Usage()
	return subcommands.ExitUsageError
}

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(&issueCmd{}, "")
	subcommands.Register(&doCmd{}, "")
	subcommands.Register(&saveCmd{}, "")
	subcommands.Register(&prCmd{}, "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}

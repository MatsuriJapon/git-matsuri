package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/MatsuriJapon/git-matsuri/matsuri"
	"github.com/google/go-github/v18/github"
	"github.com/google/subcommands"
	"golang.org/x/oauth2"
)

var client *github.Client

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(&matsuri.KanbanCmd{}, "")
	subcommands.Register(&matsuri.TodoCmd{}, "")
	subcommands.Register(&matsuri.StartCmd{}, "")
	subcommands.Register(&matsuri.SaveCmd{}, "")
	subcommands.Register(&matsuri.PrCmd{}, "")
	subcommands.Register(&matsuri.FixCmd{}, "")

	flag.Parse()
	ctx := context.Background()

	if _, err := matsuri.GetRepoName(); err != nil {
		fmt.Println("A git repository could not be found. You might not be inside a directory containing a git repository.")
		fmt.Println(err)
		return
	}

	token := os.Getenv(matsuri.TokenName)
	if token == "" {
		fmt.Printf("GitHub token not found.\nPlease create one at https://github.com/settings/tokens/new with 'repo' permissions and save it to your system environment variables under the name MATSURI_TOKEN\n")
		return
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)
	clientCtx := context.WithValue(ctx, matsuri.ContextKey("client"), client)

	os.Exit(int(subcommands.Execute(clientCtx)))
}

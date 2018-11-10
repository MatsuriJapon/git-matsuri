package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/go-github/v18/github"
	"github.com/google/subcommands"
	"golang.org/x/oauth2"
)

var client *github.Client
var owner = "MatsuriJapon"
var repo = "matsuri-japon"

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(&issueCmd{}, "")
	subcommands.Register(&startCmd{}, "")
	subcommands.Register(&saveCmd{}, "")
	subcommands.Register(&prCmd{}, "")

	flag.Parse()
	ctx := context.Background()

	// https://github.com/settings/tokens/new
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv("MATSURI_TOKEN")})
	tc := oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)

	os.Exit(int(subcommands.Execute(ctx)))
}

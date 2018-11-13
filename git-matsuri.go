package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"github.com/google/go-github/v18/github"
	"github.com/google/subcommands"
	"golang.org/x/oauth2"
)

var client *github.Client
var owner = "MatsuriJapon"
var repo = ""

func getRepoInfo() (err error) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	url, err := cmd.Output()
	if err != nil {
		return
	}
	r := regexp.MustCompile(`.+github\.com:MatsuriJapon\/(?P<repo>.+)\.git`)
	match := r.FindStringSubmatch(string(url))
	repo = match[1]
	return
}

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(&issueCmd{}, "")
	subcommands.Register(&startCmd{}, "")
	subcommands.Register(&saveCmd{}, "")
	subcommands.Register(&prCmd{}, "")

	flag.Parse()
	ctx := context.Background()

	err := getRepoInfo()
	if err != nil {
		fmt.Println(err)
	}

	// https://github.com/settings/tokens/new
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv("MATSURI_TOKEN")})
	tc := oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)

	os.Exit(int(subcommands.Execute(ctx)))
}

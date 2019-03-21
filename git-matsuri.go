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

const tokenName = "MATSURI_TOKEN"
const owner = "MatsuriJapon"

var repo string

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
	subcommands.Register(&kanbanCmd{}, "")
	subcommands.Register(&todoCmd{}, "")
	subcommands.Register(&startCmd{}, "")
	subcommands.Register(&saveCmd{}, "")
	subcommands.Register(&prCmd{}, "")
	subcommands.Register(&fixCmd{}, "")

	flag.Parse()
	ctx := context.Background()

	if err := getRepoInfo(); err != nil {
		fmt.Println(err)
		return
	}

	token := os.Getenv(tokenName)
	if token == "" {
		fmt.Printf("GitHub token not found.\nPlease create one at https://github.com/settings/tokens/new with 'repo' permissions and save it to your system environment variables under the name MATSURI_TOKEN\n")
		return
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)

	os.Exit(int(subcommands.Execute(ctx)))
}

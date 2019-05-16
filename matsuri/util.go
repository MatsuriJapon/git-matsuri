package matsuri

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"time"

	"github.com/google/go-github/v18/github"
)

// TokenName is the environment variable name for the GitHub token
const TokenName = "MATSURI_TOKEN"

const owner = "MatsuriJapon"

// ContextKey is a key to retrieve the value from a context
type ContextKey string

// GetRepoName gets the repository name from the current directory
func GetRepoName() (repo string, err error) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	url, err := cmd.Output()
	if err != nil {
		return
	}
	r := regexp.MustCompile(`.+github\.com[:\/]MatsuriJapon\/(?P<repo>.+)\.git`)
	match := r.FindStringSubmatch(string(url))
	repo = match[1]
	return
}

// GetRepoURL verifies that the given repository name matches a MatsuriJapon repository and returns its url
func GetRepoURL(ctx context.Context, name string, http bool) (url string, err error) {
	client, err := GetClient(ctx)
	if err != nil {
		return
	}
	repo, _, err := client.Repositories.Get(ctx, owner, name)
	if err != nil || repo == nil {
		return
	}
	if http {
		url = repo.GetCloneURL()
	} else {
		url = repo.GetSSHURL()
	}
	return
}

// GetClient retrieves a client from a context with value
func GetClient(ctx context.Context) (client *github.Client, err error) {
	v := ctx.Value(ContextKey("client"))
	t, ok := v.(*github.Client)
	if !ok {
		err = errors.New("Could not get a GitHub client")
		return
	}
	client = t
	return
}

// IsMainRepo checks whether the current repo is the MatsuriJapon/matsuri-japon repo
func IsMainRepo() bool {
	if repoName, err := GetRepoName(); err != nil || repoName != "matsuri-japon" {
		return false
	}
	return true
}

// IsMainRepoName checks whether the current repo is the MatsuriJapon/matsuri-japon repo
func IsMainRepoName(repoName string) bool {
	return repoName == "matsuri-japon"
}

// IsCardIssueOrPR checks whether the ProjectCard is an Issue Card
func IsCardIssueOrPR(c *github.ProjectCard) bool {
	repoName, err := GetRepoName()
	if err != nil {
		return false
	}
	base := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/\\d+", owner, repoName)
	re := regexp.MustCompile(base)
	return re.MatchString(c.GetContentURL())
}

// IsValidIssue verifies that an open Issue with the given number exists
func IsValidIssue(ctx context.Context, num int) bool {
	repoName, err := GetRepoName()
	if err != nil {
		return false
	}
	client, err := GetClient(ctx)
	if err != nil {
		return false
	}
	issue, _, err := client.Issues.Get(ctx, owner, repoName, num)
	if err != nil {
		return false
	}
	return issue.GetState() == "open" && !issue.IsPullRequest()
}

// IsExistingIssue verifies that the Issue exists and is not a pull request
func IsExistingIssue(ctx context.Context, num int) bool {
	repoName, err := GetRepoName()
	if err != nil {
		return false
	}
	client, err := GetClient(ctx)
	if err != nil {
		return false
	}
	issue, _, err := client.Issues.Get(ctx, owner, repoName, num)
	if err != nil {
		return false
	}
	return !issue.IsPullRequest()
}

// GetDefaultBranch gets the name of the default branch for the repo
func GetDefaultBranch(ctx context.Context) (branch *string, err error) {
	repoName, err := GetRepoName()
	if err != nil {
		return
	}
	client, err := GetClient(ctx)
	if err != nil {
		return
	}
	repo, _, err := client.Repositories.Get(ctx, owner, repoName)
	if err != nil {
		return
	}
	branch = repo.DefaultBranch
	return
}

// GetCurrentProjectYear gets the "current" Matsuri project year by using the default branch name. If this fails, default to the current year, determined by the current month
func GetCurrentProjectYear(ctx context.Context) (currentYear int, err error) {
	defaultBranch, err := GetDefaultBranch(ctx)
	if err != nil {
		return
	}
	r := regexp.MustCompile(`^v(?P<year>\d+)`)
	matches := r.FindStringSubmatch(*defaultBranch)
	if len(matches) == 2 {
		currentYear, _ = strconv.Atoi(matches[1])
	} else {
		currentYear = time.Now().Year()
		if time.Now().Month() > 3 {
			currentYear++
		}
	}
	return
}

// GetIssuesForProject retrieves Issues for a Project, specified by its year. If this is not the main repo, return all Issues
func GetIssuesForProject(ctx context.Context, year int) (issues []*github.Issue, err error) {
	repoName, err := GetRepoName()
	if err != nil {
		return
	}
	client, err := GetClient(ctx)
	if err != nil {
		return
	}
	if !IsMainRepoName(repoName) {
		issues, _, err = client.Issues.ListByRepo(ctx, owner, repoName, nil)
		return
	}
	project, err := GetProjectForYear(ctx, year)
	if err != nil {
		return
	}
	column, err := GetProjectColumnByName(ctx, project, "To Do")
	if err != nil {
		return
	}
	cards, _, _ := client.Projects.ListProjectCards(ctx, column.GetID(), nil)
	for i := 0; i < len(cards); i++ {
		if card := cards[i]; IsCardIssueOrPR(card) {
			num := GetIssueNumberFromCard(card)
			issue, _, _ := client.Issues.Get(ctx, owner, repoName, num)
			if issue.IsPullRequest() {
				continue
			}
			issues = append(issues, issue)
		}
	}
	return
}

// GetIssueNumberFromCard gets the Issue number from a Card
func GetIssueNumberFromCard(c *github.ProjectCard) (id int) {
	r := regexp.MustCompile(`(?:issues/)(?P<id>\d+)`)
	matches := r.FindStringSubmatch(c.GetContentURL())
	id, _ = strconv.Atoi(matches[1])
	return
}

// GetProjectForYear gets the project associated with the current Matsuri year
func GetProjectForYear(ctx context.Context, year int) (project *github.Project, err error) {
	repoName, err := GetRepoName()
	if err != nil {
		return
	}
	client, err := GetClient(ctx)
	if err != nil {
		return
	}
	if !IsMainRepoName(repoName) {
		err = fmt.Errorf("Error: this repository doesn't have a Project board")
		return
	}
	projectName := fmt.Sprintf("Matsuri %d", year)
	projects, _, _ := client.Repositories.ListProjects(ctx, owner, repoName, nil)
	for i := 0; i < len(projects); i++ {
		if projects[i].GetName() == projectName {
			project = projects[i]
			return
		}
	}
	err = fmt.Errorf("Error: Project %s was not found", projectName)
	return
}

// GetProjectColumnByName gets the column by its name
func GetProjectColumnByName(ctx context.Context, project *github.Project, columnName string) (column *github.ProjectColumn, err error) {
	client, err := GetClient(ctx)
	if err != nil {
		return
	}
	columns, _, _ := client.Projects.ListProjectColumns(ctx, project.GetID(), nil)
	for i := 0; i < len(columns); i++ {
		if columns[i].GetName() == columnName {
			column = columns[i]
			return
		}
	}
	err = fmt.Errorf("Error: there is no %s column for %s", columnName, project.GetName())
	return
}

// GetProjectCardInColumn gets the project card associated with the given issue number
func GetProjectCardInColumn(ctx context.Context, column *github.ProjectColumn, issueNumber int) *github.ProjectCard {
	client, err := GetClient(ctx)
	if err != nil {
		return nil
	}
	cards, _, _ := client.Projects.ListProjectCards(ctx, column.GetID(), nil)
	for i := 0; i < len(cards); i++ {
		if card := cards[i]; IsCardIssueOrPR(card) {
			num := GetIssueNumberFromCard(card)
			if num == issueNumber {
				return card
			}
		}
	}
	return nil
}

// GetRepoIssues gets the issues for the current repository
func GetRepoIssues(ctx context.Context) (issues []*github.Issue, err error) {
	repoName, err := GetRepoName()
	if err != nil {
		return
	}
	client, err := GetClient(ctx)
	if err != nil {
		return
	}
	issues, _, err = client.Issues.ListByRepo(ctx, owner, repoName, nil)
	return
}

func createPR(ctx context.Context, newPr *github.NewPullRequest) (pr *github.PullRequest, err error) {
	repoName, err := GetRepoName()
	if err != nil {
		return
	}
	client, err := GetClient(ctx)
	if err != nil {
		return
	}
	projectYear, _ := GetCurrentProjectYear(ctx)
	pr, _, err = client.PullRequests.Create(ctx, owner, repoName, newPr)
	if err != nil || !IsMainRepoName(repoName) {
		return
	}
	project, err := GetProjectForYear(ctx, projectYear)
	if err != nil {
		return
	}
	todo, err := GetProjectColumnByName(ctx, project, "To Do")
	if err != nil {
		return
	}
	cardOpt := &github.ProjectCardOptions{
		ContentID:   pr.GetID(),
		ContentType: "PullRequest",
	}
	_, _, err = client.Projects.CreateProjectCard(ctx, todo.GetID(), cardOpt)
	return
}

// CreatePRForIssueNumber creates a new PR for the given issue and returns the created card
func CreatePRForIssueNumber(ctx context.Context, issueNum int, noclose bool) (pr *github.PullRequest, err error) {
	repoName, err := GetRepoName()
	if err != nil {
		return
	}
	client, err := GetClient(ctx)
	if err != nil {
		return
	}
	issue, _, err := client.Issues.Get(ctx, owner, repoName, issueNum)
	if err != nil {
		return
	}
	title := fmt.Sprintf("ISSUE-%d: %s", issue.GetNumber(), issue.GetTitle())
	head := fmt.Sprintf("ISSUE-%d", issue.GetNumber())
	base, err := GetDefaultBranch(ctx)
	if err != nil {
		return
	}
	body := fmt.Sprintf("Closes #%d\n", issue.GetNumber())
	if noclose {
		body = fmt.Sprintf("Related to #%d\n", issue.GetNumber())
	}
	newPr := &github.NewPullRequest{
		Title: github.String(title),
		Head:  github.String(head),
		Base:  base,
		Body:  github.String(body),
	}
	return createPR(ctx, newPr)
}

// CreateFixPRForIssueNumber creates a fix PR for the provided issue
func CreateFixPRForIssueNumber(ctx context.Context, issueNum int, noclose bool) (pr *github.PullRequest, err error) {
	repoName, err := GetRepoName()
	if err != nil {
		return
	}
	client, err := GetClient(ctx)
	if err != nil {
		return
	}
	issue, _, err := client.Issues.Get(ctx, owner, repoName, issueNum)
	if err != nil {
		return
	}
	title := fmt.Sprintf("ISSUE-%d-fix: %s", issue.GetNumber(), issue.GetTitle())
	head := fmt.Sprintf("ISSUE-%d", issue.GetNumber())
	base, err := GetDefaultBranch(ctx)
	if err != nil {
		return
	}
	body := fmt.Sprintf("Fixes PR for #%d\n", issue.GetNumber())
	if !noclose {
		body += fmt.Sprintf("Closes #%d\n", issue.GetNumber())
	}
	newPr := &github.NewPullRequest{
		Title: github.String(title),
		Head:  github.String(head),
		Base:  base,
		Body:  github.String(body),
	}
	return createPR(ctx, newPr)
}

// MoveProjectCardForProject moves the Issue to the Doing column of the current Matsuri project year
func MoveProjectCardForProject(ctx context.Context, num int, year int) (err error) {
	project, err := GetProjectForYear(ctx, year)
	if err != nil {
		return
	}
	todo, err := GetProjectColumnByName(ctx, project, "To Do")
	if err != nil {
		return
	}
	doing, err := GetProjectColumnByName(ctx, project, "In progress")
	if err != nil {
		return err
	}
	card := GetProjectCardInColumn(ctx, todo, num)
	if card == nil {
		// handle the case where the Issue has already been moved to Doing
		card = GetProjectCardInColumn(ctx, doing, num)
		if card == nil {
			return fmt.Errorf("The specified Issue is not in %s's To Do or Doing columns", project.GetName())
		}
		return
	}
	opt := &github.ProjectCardMoveOptions{
		Position: "top",
		ColumnID: doing.GetID(),
	}
	client, err := GetClient(ctx)
	if err != nil {
		return
	}
	_, err = client.Projects.MoveProjectCard(ctx, card.GetID(), opt)
	return
}

// ReopenIssue reopens a closed Issue
func ReopenIssue(ctx context.Context, issueNum int) (err error) {
	repoName, err := GetRepoName()
	if err != nil {
		return
	}
	client, err := GetClient(ctx)
	if err != nil {
		return
	}
	issue, _, err := client.Issues.Get(ctx, owner, repoName, issueNum)
	if err != nil {
		return
	}
	if issue.GetState() == "open" || issue.IsPullRequest() {
		err = nil
		return
	}
	reopenRequest := &github.IssueRequest{
		State: github.String("open"),
	}
	_, _, err = client.Issues.Edit(ctx, owner, repoName, issueNum, reopenRequest)
	return
}

// PrintProjectKanban prints the project kanban
func PrintProjectKanban(ctx context.Context, project *github.Project) {
	repoName, err := GetRepoName()
	if err != nil {
		return
	}
	client, err := GetClient(ctx)
	if err != nil {
		return
	}
	columns, _, _ := client.Projects.ListProjectColumns(ctx, project.GetID(), nil)
	for i := 0; i < len(columns); i++ {
		column := columns[i]
		fmt.Println(column.GetName())
		cards, _, _ := client.Projects.ListProjectCards(ctx, column.GetID(), nil)
		for j := 0; j < len(cards); j++ {
			card := cards[j]
			num := GetIssueNumberFromCard(card)
			issue, _, _ := client.Issues.Get(ctx, owner, repoName, num)
			fmt.Printf("%d: %s\n", issue.GetNumber(), issue.GetTitle())
		}
		fmt.Println()
	}
}

// PrintIssues prints issues
func PrintIssues(issues []*github.Issue) {
	for i := 0; i < len(issues); i++ {
		issue := issues[i]
		// sanity check
		if issue.IsPullRequest() {
			continue
		}
		fmt.Printf("%d: %s\n", issue.GetNumber(), issue.GetTitle())
	}
}

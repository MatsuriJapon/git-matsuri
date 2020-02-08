package matsuri

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/go-github/v29/github"
	"github.com/hashicorp/go-version"
	"golang.org/x/oauth2"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"time"
)

// TokenName is the environment variable name for the GitHub token
const TokenName = "MATSURI_TOKEN"

const owner = "MatsuriJapon"

var ctx = context.Background()

// GetLatestVersion gets the release tag of the latest release version of git-matsuri
func GetLatestVersion() (v *version.Version, err error) {
	client := GetClient()

	release, _, err := client.Repositories.GetLatestRelease(ctx, owner, "git-matsuri")
	if err != nil {
		return
	}
	r := regexp.MustCompile(`^v(?P<version>\d+\.\d+\.\d+)$`)
	matches := r.FindStringSubmatch(*release.TagName)
	if len(matches) != 2 {
		err = errors.New("Could not retrieve version number")
		return
	}
	v, err = version.NewVersion(matches[1])
	return
}

// GetMatsuriEmail gets the festivaljapon.com email of the current user, if available
func GetMatsuriEmail() (email string, err error) {
	client := GetClient()
	userEmails, _, err := client.Users.ListEmails(ctx, nil)
	if err != nil {
		return
	}
	r := regexp.MustCompile(`^[\w._+-]+@festivaljapon.com$`)
	for _, userEmail := range userEmails {
		matches := r.FindStringSubmatch(userEmail.GetEmail())
		if len(matches) != 1 {
			continue
		}
		email = matches[0]
		return
	}
	return
}

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
func GetRepoURL(name string, http bool) (url string, err error) {
	client := GetClient()

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
func GetClient() (client *github.Client) {
	ctx := context.Background()
	token := os.Getenv(TokenName)
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)
	return
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
func IsValidIssue(num int) bool {
	repoName, err := GetRepoName()
	if err != nil {
		return false
	}
	client := GetClient()
	issue, _, err := client.Issues.Get(ctx, owner, repoName, num)
	if err != nil {
		return false
	}
	return issue.GetState() == "open" && !issue.IsPullRequest()
}

// IsExistingIssue verifies that the Issue exists and is not a pull request
func IsExistingIssue(num int) bool {
	repoName, err := GetRepoName()
	if err != nil {
		return false
	}
	client := GetClient()
	issue, _, err := client.Issues.Get(ctx, owner, repoName, num)
	if err != nil {
		return false
	}
	return !issue.IsPullRequest()
}

// GetDefaultBranch gets the name of the default branch for the repo
func GetDefaultBranch() (branch *string, err error) {
	repoName, err := GetRepoName()
	if err != nil {
		return
	}
	client := GetClient()
	repo, _, err := client.Repositories.Get(ctx, owner, repoName)
	if err != nil {
		return
	}
	branch = repo.DefaultBranch
	return
}

// GetCurrentProjectYear gets the "current" Matsuri project year by using the default branch name. If this fails, default to the current year, determined by the current month
func GetCurrentProjectYear() (currentYear int, err error) {
	defaultBranch, err := GetDefaultBranch()
	if err != nil {
		return
	}
	r := regexp.MustCompile(`^v(?P<year>\d+)`)
	matches := r.FindStringSubmatch(*defaultBranch)
	if len(matches) == 2 {
		currentYear, _ = strconv.Atoi(matches[1])
	} else {
		currentYear = time.Now().Year()
		if time.Now().Month() > 2 {
			currentYear++
		}
	}
	return
}

// GetIssuesForProject retrieves Issues for a Project, specified by its year
func GetIssuesForProject(year int) (issues []*github.Issue, err error) {
	repoName, err := GetRepoName()
	if err != nil {
		return
	}
	client := GetClient()

	project, err := GetProjectForYear(year)
	if err != nil {
		return
	}
	column, err := GetProjectColumnByName(project, "To do")
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

// GetRepoNameFromURL gets the Repository name from a URL
func GetRepoNameFromURL(url string) (repoName string) {
	r := regexp.MustCompile(`(?:MatsuriJapon/)(?P<repoName>[^/]+)`)
	matches := r.FindStringSubmatch(url)
	repoName = matches[1]
	return
}

// GetProjectForYear gets the project associated with the current Matsuri year
func GetProjectForYear(year int) (project *github.Project, err error) {
	client := GetClient()

	projectName := fmt.Sprintf("Matsuri %d", year)
	projects, _, err := client.Organizations.ListProjects(ctx, owner, nil)
	if err != nil {
		return
	}
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
func GetProjectColumnByName(project *github.Project, columnName string) (column *github.ProjectColumn, err error) {
	client := GetClient()

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
func GetProjectCardInColumn(column *github.ProjectColumn, issueNumber int) *github.ProjectCard {
	repoName, err := GetRepoName()
	if err != nil {
		return nil
	}
	client := GetClient()
	cards, _, _ := client.Projects.ListProjectCards(ctx, column.GetID(), nil)
	for i := 0; i < len(cards); i++ {
		if card := cards[i]; IsCardIssueOrPR(card) {
			num := GetIssueNumberFromCard(card)
			_, _, err := client.Issues.Get(ctx, owner, repoName, num)
			if err != nil {
				continue
			}
			if num == issueNumber {
				return card
			}
		}
	}
	return nil
}

// GetRepoIssues gets the issues for the current repository
func GetRepoIssues() (issues []*github.Issue, err error) {
	repoName, err := GetRepoName()
	if err != nil {
		return
	}
	client := GetClient()
	issues, _, err = client.Issues.ListByRepo(ctx, owner, repoName, nil)
	return
}

func createPR(newPr *github.NewPullRequest) (pr *github.PullRequest, err error) {
	repoName, err := GetRepoName()
	if err != nil {
		return
	}
	client := GetClient()

	projectYear, _ := GetCurrentProjectYear()
	pr, _, err = client.PullRequests.Create(ctx, owner, repoName, newPr)
	if err != nil {
		return
	}
	project, err := GetProjectForYear(projectYear)
	if err != nil {
		return
	}
	todo, err := GetProjectColumnByName(project, "To do")
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
func CreatePRForIssueNumber(issueNum int, noclose bool) (pr *github.PullRequest, err error) {
	repoName, err := GetRepoName()
	if err != nil {
		return
	}
	client := GetClient()
	issue, _, err := client.Issues.Get(ctx, owner, repoName, issueNum)
	if err != nil {
		return
	}
	title := fmt.Sprintf("ISSUE-%d: %s", issue.GetNumber(), issue.GetTitle())
	head := fmt.Sprintf("ISSUE-%d", issue.GetNumber())
	base, err := GetDefaultBranch()
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
	return createPR(newPr)
}

// CreateFixPRForIssueNumber creates a fix PR for the provided issue
func CreateFixPRForIssueNumber(issueNum int, noclose bool) (pr *github.PullRequest, err error) {
	repoName, err := GetRepoName()
	if err != nil {
		return
	}
	client := GetClient()
	issue, _, err := client.Issues.Get(ctx, owner, repoName, issueNum)
	if err != nil {
		return
	}
	title := fmt.Sprintf("ISSUE-%d-fix: %s", issue.GetNumber(), issue.GetTitle())
	head := fmt.Sprintf("ISSUE-%d", issue.GetNumber())
	base, err := GetDefaultBranch()
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
	return createPR(newPr)
}

// MoveProjectCardForProject moves the Issue to the Doing column of the current Matsuri project year
func MoveProjectCardForProject(num int, year int) (err error) {
	project, err := GetProjectForYear(year)
	if err != nil {
		return
	}
	todo, err := GetProjectColumnByName(project, "To do")
	if err != nil {
		return
	}
	doing, err := GetProjectColumnByName(project, "In progress")
	if err != nil {
		return
	}
	card := GetProjectCardInColumn(todo, num)
	if card == nil {
		// handle the case where the Issue has already been moved to Doing
		card = GetProjectCardInColumn(doing, num)
		if card == nil {
			return fmt.Errorf("The specified Issue is not in %s's To do or Doing columns", project.GetName())
		}
		return
	}
	opt := &github.ProjectCardMoveOptions{
		Position: "top",
		ColumnID: doing.GetID(),
	}
	client := GetClient()
	_, err = client.Projects.MoveProjectCard(ctx, card.GetID(), opt)
	return
}

// ReopenIssue reopens a closed Issue
func ReopenIssue(issueNum int) (err error) {
	repoName, err := GetRepoName()
	if err != nil {
		return
	}
	client := GetClient()

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
func PrintProjectKanban(project *github.Project) {
	client := GetClient()
	columns, _, _ := client.Projects.ListProjectColumns(ctx, project.GetID(), nil)
	for i := 0; i < len(columns); i++ {
		column := columns[i]
		fmt.Println(column.GetName())
		cards, _, _ := client.Projects.ListProjectCards(ctx, column.GetID(), nil)
		for j := 0; j < len(cards); j++ {
			card := cards[j]
			num := GetIssueNumberFromCard(card)
			repoName := GetRepoNameFromURL(card.GetContentURL())
			issue, _, _ := client.Issues.Get(ctx, owner, repoName, num)
			fmt.Printf("%d [%s]: %s\n", issue.GetNumber(), repoName, issue.GetTitle())
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
		repoName := GetRepoNameFromURL(issue.GetRepositoryURL())
		fmt.Printf("%d [%s]: %s\n", issue.GetNumber(), repoName, issue.GetTitle())
	}
}

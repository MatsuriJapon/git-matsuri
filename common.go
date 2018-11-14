package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/v18/github"
	"regexp"
	"strconv"
	"time"
)

// GetIssuesForProject retrieves Issues for a Project, specified by its year. If this is not the main repo, return all Issues
func GetIssuesForProject(ctx context.Context, year int) (issues []*github.Issue, err error) {
	if !IsMainRepo() {
		issues, _, err = client.Issues.ListByRepo(ctx, owner, repo, nil)
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
			issue, _, _ := client.Issues.Get(ctx, owner, repo, num)
			if issue.IsPullRequest() {
				continue
			}
			issues = append(issues, issue)
		}
	}
	return
}

// GetProjectForYear gets the project associated with the current Matsuri year
func GetProjectForYear(ctx context.Context, year int) (project *github.Project, err error) {
	if !IsMainRepo() {
		err = fmt.Errorf("Error: this repository doesn't have a Project board")
		return
	}
	projectName := fmt.Sprintf("Matsuri %d", year)
	projects, _, _ := client.Repositories.ListProjects(ctx, owner, repo, nil)
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
	_, err = client.Projects.MoveProjectCard(ctx, card.GetID(), opt)
	return
}

// CreatePRForIssueNumber creates a new PR for the given issue and returns the created card
func CreatePRForIssueNumber(ctx context.Context, issueNum int, noclose bool) (pr *github.PullRequest, err error) {
	issue, _, err := client.Issues.Get(ctx, owner, repo, issueNum)
	if err != nil {
		return
	}
	projectYear := GetCurrentProjectYear()
	title := fmt.Sprintf("ISSUE-%d: %s", issue.GetNumber(), issue.GetTitle())
	head := fmt.Sprintf("ISSUE-%d", issue.GetNumber())
	base := fmt.Sprintf("v%d", projectYear)
	if !IsMainRepo() {
		base = "master"
	}
	body := fmt.Sprintf("Closes #%d\n", issue.GetNumber())
	if noclose {
		body = fmt.Sprintf("Related to #%d\n", issue.GetNumber())
	}
	newPr := &github.NewPullRequest{
		Title: github.String(title),
		Head:  github.String(head),
		Base:  github.String(base),
		Body:  github.String(body),
	}
	pr, _, err = client.PullRequests.Create(ctx, owner, repo, newPr)
	if err != nil || !IsMainRepo() {
		return
	}
	project, err := GetProjectForYear(ctx, projectYear)
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

// IsCardIssueOrPR checks whether the ProjectCard is an Issue Card
func IsCardIssueOrPR(c *github.ProjectCard) bool {
	base := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/\\d+", owner, repo)
	re := regexp.MustCompile(base)
	return re.MatchString(c.GetContentURL())
}

// GetIssueNumberFromCard gets the Issue number from a Card
func GetIssueNumberFromCard(c *github.ProjectCard) (id int) {
	r := regexp.MustCompile(`(?:issues/)(?P<id>\d+)`)
	matches := r.FindStringSubmatch(c.GetContentURL())
	id, _ = strconv.Atoi(matches[1])
	return
}

// IsValidIssue verifies that an open Issue with the given number exists
func IsValidIssue(ctx context.Context, num int) bool {
	issue, _, err := client.Issues.Get(ctx, owner, repo, num)
	if err != nil {
		return false
	}
	return issue.GetState() == "open" && !issue.IsPullRequest()
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

// GetCurrentProjectYear weird logic to get "current" Matsuri project year
func GetCurrentProjectYear() (currentYear int) {
	currentYear = time.Now().Year()
	if time.Now().Month() > 8 {
		currentYear++
	}
	return
}

// IsMainRepo checks whether the current repo is the MatsuriJapon/matsuri-japon repo
func IsMainRepo() bool {
	return repo == "matsuri-japon"
}

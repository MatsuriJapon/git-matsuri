package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/v18/github"
	"regexp"
	"strconv"
)

// GetIssuesForProject retrieves Issues for a Project, specified by its year
func GetIssuesForProject(ctx context.Context, year int) ([]*github.Issue, error) {
	projectName := fmt.Sprintf("Matsuri %d", year)
	projects, _, _ := client.Repositories.ListProjects(ctx, owner, repo, nil)
	var project *github.Project
	for i := 0; i < len(projects); i++ {
		if projects[i].GetName() == projectName {
			project = projects[i]
			break
		}
	}
	if project == nil {
		return nil, fmt.Errorf("Error: Project %s was not found", projectName)
	}

	columns, _, _ := client.Projects.ListProjectColumns(ctx, project.GetID(), nil)
	var column *github.ProjectColumn
	for i := 0; i < len(columns); i++ {
		if columns[i].GetName() == "To Do" {
			column = columns[i]
			break
		}
	}
	if column == nil {
		return nil, fmt.Errorf("Error: there is no To Do column for %s", projectName)
	}
	cards, _, _ := client.Projects.ListProjectCards(ctx, column.GetID(), nil)
	var issues []*github.Issue
	for i := 0; i < len(cards); i++ {
		if card := cards[i]; IsCardIssue(card) {
			id := GetIssueIDFromCard(card)
			issue, _, _ := client.Issues.Get(ctx, owner, repo, id)
			issues = append(issues, issue)
		}
	}
	return issues, nil
}

// IsCardIssue checks whether the ProjectCard is an Issue Card
func IsCardIssue(c *github.ProjectCard) bool {
	base := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/\\d+", owner, repo)
	re := regexp.MustCompile(base)
	return re.MatchString(c.GetContentURL())
}

// GetIssueIDFromCard gets the Issue number from a Card
func GetIssueIDFromCard(c *github.ProjectCard) int {
	r := regexp.MustCompile(`(?:issues/)(?P<id>\d+)`)
	matches := r.FindStringSubmatch(c.GetContentURL())
	id, _ := strconv.Atoi(matches[1])
	return id
}

// PrintIssues prints issues
func PrintIssues(issues []*github.Issue) {
	for i := 0; i < len(issues); i++ {
		issue := issues[i]
		fmt.Printf("%d: %s\n", issue.GetNumber(), issue.GetTitle())
	}
}

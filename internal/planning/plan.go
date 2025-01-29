package planning

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cli/go-gh/v2/pkg/api"
)

// Manager handles GitHub planning
type Manager struct {
	client *api.RESTClient
}

// NewManager creates a new plan manager
func NewManager(client *api.RESTClient) *Manager {
	return &Manager{
		client: client,
	}
}

// Issue represents a GitHub issue
type Issue struct {
	Title    string   `json:"title"`
	Number   int      `json:"number"`
	State    string   `json:"state"`
	Labels   []Label  `json:"labels"`
	Assignee *User    `json:"assignee"`
}

// Label represents a GitHub label
type Label struct {
	Name string `json:"name"`
}

// User represents a GitHub user
type User struct {
	Login string `json:"login"`
}

// Milestone represents a GitHub milestone
type Milestone struct {
	Title       string     `json:"title"`
	DueOn       *time.Time `json:"due_on"`
	Description string     `json:"description"`
	Number      int        `json:"number"`
	State       string     `json:"state"`
}

// Options represents planning options
type Options struct {
	PlanningLabel string
	Categories    []string
	ExcludePR     bool
}

// DefaultOptions returns default planning options
func DefaultOptions() Options {
	return Options{
		PlanningLabel: "planning",
		Categories:    []string{"bug", "documentation", "enhancement"},
		ExcludePR:     true,
	}
}

// Update updates or creates a planning issue for a milestone
func (m *Manager) Update(ctx context.Context, owner, repo string, milestoneNumber int, opts Options) error {
	// Get milestone details
	var milestone Milestone
	err := m.client.Get(fmt.Sprintf("repos/%s/%s/milestones/%d", owner, repo, milestoneNumber), &milestone)
	if err != nil {
		return fmt.Errorf("failed to get milestone: %w", err)
	}

	// Skip if milestone is closed
	if milestone.State == "closed" {
		return fmt.Errorf("milestone #%d is closed", milestoneNumber)
	}

	// Get all issues for this milestone
	var allIssues []struct {
		Issue
		PullRequest interface{} `json:"pull_request"`
		HTMLURL     string      `json:"html_url"`
	}
	page := 1
	perPage := 30

	for {
		var issues []struct {
			Issue
			PullRequest interface{} `json:"pull_request"`
			HTMLURL     string      `json:"html_url"`
		}
		err := m.client.Get(fmt.Sprintf("repos/%s/%s/issues?milestone=%d&state=all&per_page=%d&page=%d", 
			owner, repo, milestoneNumber, perPage, page), &issues)
		if err != nil {
			return fmt.Errorf("failed to get issues: %w", err)
		}

		allIssues = append(allIssues, issues...)

		if len(issues) < perPage {
			break
		}
		page++
	}

	// Filter out pull requests if excludePR is true
	var issues []Issue
	for _, item := range allIssues {
		if opts.ExcludePR {
			isPR := item.PullRequest != nil || 
				strings.Contains(item.HTMLURL, "/pull/")
			if isPR {
				continue
			}
		}
		issues = append(issues, item.Issue)
	}

	// Generate planning content
	content := generatePlanningContent(milestone, issues, opts.Categories)

	// Find existing planning issue
	var existingIssues []struct {
		Number int    `json:"number"`
		Title  string `json:"title"`
	}
	err = m.client.Get(fmt.Sprintf("repos/%s/%s/issues?labels=%s&state=open", 
		owner, repo, opts.PlanningLabel), &existingIssues)
	if err != nil {
		return fmt.Errorf("failed to get existing issues: %w", err)
	}

	var planningIssue *struct {
		Number int    `json:"number"`
		Title  string `json:"title"`
	}
	for _, issue := range existingIssues {
		if strings.Contains(issue.Title, milestone.Title) {
			planningIssue = &issue
			break
		}
	}

	if planningIssue != nil {
		// Update existing issue
		body := map[string]string{"body": content}
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}

		err = m.client.Patch(fmt.Sprintf("repos/%s/%s/issues/%d", 
			owner, repo, planningIssue.Number), 
			bytes.NewReader(bodyBytes), nil)
		if err != nil {
			return fmt.Errorf("failed to update planning issue: %w", err)
		}
	} else {
		// Create new planning issue
		body := map[string]interface{}{
			"title":  fmt.Sprintf("Planning: %s", milestone.Title),
			"body":   content,
			"labels": []string{opts.PlanningLabel},
		}
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}

		err = m.client.Post(fmt.Sprintf("repos/%s/%s/issues", owner, repo), 
			bytes.NewReader(bodyBytes), nil)
		if err != nil {
			return fmt.Errorf("failed to create planning issue: %w", err)
		}
	}

	return nil
}

// generatePlanningContent generates the complete planning content
func generatePlanningContent(milestone Milestone, issues []Issue, categories []string) string {
	// Calculate statistics
	totalIssues := len(issues)
	var completedIssues int
	for _, issue := range issues {
		if issue.State == "closed" {
			completedIssues++
		}
	}
	inProgressIssues := totalIssues - completedIssues

	progressBar := generateProgressBar(completedIssues, totalIssues, 20)
	dueDate := formatDate(milestone.DueOn)

	var sb strings.Builder

	// Generate main content
	fmt.Fprintf(&sb, "# %s Planning\n\n", milestone.Title)

	fmt.Fprintf(&sb, "## Overview\n")
	fmt.Fprintf(&sb, "- Progress: %s\n", progressBar)
	fmt.Fprintf(&sb, "- Total Issues: %d\n", totalIssues)
	fmt.Fprintf(&sb, "  - Completed: %d\n", completedIssues)
	fmt.Fprintf(&sb, "  - In Progress: %d\n", inProgressIssues)
	fmt.Fprintf(&sb, "- Due Date: %s\n\n", dueDate)

	fmt.Fprintf(&sb, "## Description\n")
	if milestone.Description == "" {
		fmt.Fprintf(&sb, "No description provided.\n\n")
	} else {
		fmt.Fprintf(&sb, "%s\n\n", milestone.Description)
	}

	fmt.Fprintf(&sb, "## Tasks by Category\n")

	// Generate sections for each category
	for _, category := range categories {
		section := generateCategorySection(issues, category)
		if section != "" {
			fmt.Fprintf(&sb, "\n%s", section)
		}
	}

	// Generate uncategorized section
	uncategorizedSection := generateUncategorizedSection(issues, categories)
	if uncategorizedSection != "" {
		fmt.Fprintf(&sb, "\n%s", uncategorizedSection)
	}

	// Add contributors section (only for completed issues)
	var contributors []string
	contributorsMap := make(map[string]bool)
	for _, issue := range issues {
		if issue.State == "closed" && issue.Assignee != nil {
			contributorsMap[issue.Assignee.Login] = true
		}
	}
	for contributor := range contributorsMap {
		contributors = append(contributors, contributor)
	}

	if len(contributors) > 0 {
		fmt.Fprintf(&sb, "\n## Contributors\n")
		fmt.Fprintf(&sb, "Thanks to all our contributors for their efforts on completed issues:\n\n")
		for _, contributor := range contributors {
			fmt.Fprintf(&sb, "- @%s\n", contributor)
		}
	}

	// Add footer
	fmt.Fprintf(&sb, "\n---\n")
	fmt.Fprintf(&sb, "> Auto-generated by [OSP](https://github.com/elliotxx/osp). DO NOT EDIT.\n")
	fmt.Fprintf(&sb, "> Last Updated: %s\n", time.Now().Format("January 2, 2006 15:04 MST"))

	return sb.String()
}

func generateProgressBar(completed int, total int, length int) string {
	var sb strings.Builder
	percentage := float64(completed) / float64(total) * 100
	blocks := int(percentage / 100 * float64(length))
	for i := 0; i < length; i++ {
		if i < blocks {
			sb.WriteString("█")
		} else {
			sb.WriteString("░")
		}
	}
	return fmt.Sprintf("[%s] %.2f%%", sb.String(), percentage)
}

func formatDate(date *time.Time) string {
	if date == nil {
		return "No due date"
	}
	return date.Format("January 2, 2006")
}

func generateCategorySection(issues []Issue, category string) string {
	var sb strings.Builder
	var issuesInCategory []Issue
	for _, issue := range issues {
		for _, label := range issue.Labels {
			if label.Name == category {
				issuesInCategory = append(issuesInCategory, issue)
				break
			}
		}
	}
	if len(issuesInCategory) == 0 {
		return ""
	}
	fmt.Fprintf(&sb, "### %s\n", category)
	for _, issue := range issuesInCategory {
		fmt.Fprintf(&sb, "- [%s](https://github.com/%s/%s/issues/%d) - %s\n", 
			issue.Title, "your-username", "your-repo", issue.Number, issue.State)
	}
	return sb.String()
}

func generateUncategorizedSection(issues []Issue, categories []string) string {
	var sb strings.Builder
	var uncategorizedIssues []Issue
	for _, issue := range issues {
		isCategorized := false
		for _, category := range categories {
			for _, label := range issue.Labels {
				if label.Name == category {
					isCategorized = true
					break
				}
			}
			if isCategorized {
				break
			}
		}
		if !isCategorized {
			uncategorizedIssues = append(uncategorizedIssues, issue)
		}
	}
	if len(uncategorizedIssues) == 0 {
		return ""
	}
	fmt.Fprintf(&sb, "### Uncategorized\n")
	for _, issue := range uncategorizedIssues {
		fmt.Fprintf(&sb, "- [%s](https://github.com/%s/%s/issues/%d) - %s\n", 
			issue.Title, "your-username", "your-repo", issue.Number, issue.State)
	}
	return sb.String()
}

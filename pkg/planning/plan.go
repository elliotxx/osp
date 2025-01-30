package planning

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"text/template"
	"time"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/elliotxx/osp/pkg/log"
)

//go:embed templates/*.gotmpl
var templates embed.FS

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
	Title    string  `json:"title"`
	Number   int     `json:"number"`
	State    string  `json:"state"`
	Labels   []Label `json:"labels"`
	Assignee *User   `json:"assignee"`
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
	HTMLURL     string     `json:"html_url"`
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

// MilestoneStats represents milestone statistics
type MilestoneStats struct {
	TotalIssues     int
	CompletedIssues int
	Progress        float64
	Contributors    []string
}

// TemplateData represents the data passed to the template
type TemplateData struct {
	Milestone           Milestone
	Stats               MilestoneStats
	Categories          []string
	Issues              map[string][]Issue
	UncategorizedIssues []Issue
	ProgressBar         string
}

// Update updates or creates a planning issue for a milestone
func (m *Manager) Update(ctx context.Context, owner, repo string, milestoneNumber int, opts Options) error {
	log.Debug("Updating planning issue for milestone #%d in %s/%s", milestoneNumber, owner, repo)

	// Get milestone
	var milestone Milestone
	err := m.client.Get(fmt.Sprintf("repos/%s/%s/milestones/%d", owner, repo, milestoneNumber), &milestone)
	if err != nil {
		return fmt.Errorf("failed to get milestone: %w", err)
	}
	log.Debug("Found milestone: %s (#%d)", milestone.Title, milestone.Number)

	// Get all issues in the milestone
	var issues []Issue
	err = m.client.Get(fmt.Sprintf("repos/%s/%s/issues?milestone=%d&state=all", owner, repo, milestoneNumber), &issues)
	if err != nil {
		return fmt.Errorf("failed to get issues: %w", err)
	}
	log.Debug("Found %d issues in milestone", len(issues))

	// Filter out PRs if needed
	if opts.ExcludePR {
		var filteredIssues []Issue
		for _, issue := range issues {
			isPR := false
			for _, label := range issue.Labels {
				if strings.HasPrefix(label.Name, "pr/") {
					isPR = true
					break
				}
			}
			if !isPR {
				filteredIssues = append(filteredIssues, issue)
			}
		}
		issues = filteredIssues
	}

	// Prepare template data
	data := m.prepareTemplateData(milestone, issues, opts.Categories)

	// Generate content
	content, err := m.generatePlanningContent(data)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}
	log.Debug("Generated planning content with %d bytes", len(content))

	// Find existing planning issue
	var existingIssues []Issue
	err = m.client.Get(fmt.Sprintf("repos/%s/%s/issues?labels=%s&state=all", owner, repo, opts.PlanningLabel), &existingIssues)
	if err != nil {
		return fmt.Errorf("failed to get existing issues: %w", err)
	}
	log.Debug("Found %d existing issues with planning label", len(existingIssues))

	planningTitle := fmt.Sprintf("Planning: %s", milestone.Title)
	var planningIssue *Issue
	var minIssueNumber int = math.MaxInt32
	for _, issue := range existingIssues {
		if issue.Title == planningTitle {
			if planningIssue == nil || issue.Number < minIssueNumber {
				planningIssue = &issue
				minIssueNumber = issue.Number
				log.Debug("Found planning issue #%d with title '%s'", issue.Number, issue.Title)
			}
		}
	}

	// Create or update planning issue
	if planningIssue == nil {
		log.Info("Creating new planning issue for milestone '%s'", milestone.Title)
		// Create new issue
		body := map[string]interface{}{
			"title":  planningTitle,
			"body":   content,
			"labels": []string{opts.PlanningLabel},
		}
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}

		err = m.client.Post(fmt.Sprintf("repos/%s/%s/issues", owner, repo), bytes.NewReader(bodyBytes), nil)
		if err != nil {
			return fmt.Errorf("failed to create planning issue: %w", err)
		}
		log.Success("Successfully created planning issue for milestone '%s'", milestone.Title)
	} else {
		log.Info("Updating existing planning issue #%d", planningIssue.Number)
		// Update existing issue
		body := map[string]interface{}{
			"body": content,
		}
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}

		err = m.client.Patch(fmt.Sprintf("repos/%s/%s/issues/%d", owner, repo, planningIssue.Number), bytes.NewReader(bodyBytes), nil)
		if err != nil {
			return fmt.Errorf("failed to update planning issue: %w", err)
		}
		log.Success("Successfully updated planning issue #%d", planningIssue.Number)
	}

	return nil
}

// prepareTemplateData prepares data for the template
func (m *Manager) prepareTemplateData(milestone Milestone, issues []Issue, categories []string) TemplateData {
	// Calculate statistics
	totalIssues := len(issues)
	completedIssues := 0
	contributors := make(map[string]bool)

	// Count completed issues and collect contributors of completed issues
	for _, issue := range issues {
		if issue.State == "closed" {
			completedIssues++
			if issue.Assignee != nil {
				contributors[issue.Assignee.Login] = true
			}
		}
	}

	// Get unique contributors
	var contributorsList []string
	for contributor := range contributors {
		contributorsList = append(contributorsList, contributor)
	}

	// Group issues by category
	issuesByCategory := make(map[string][]Issue)
	var uncategorizedIssues []Issue

	for _, issue := range issues {
		categorized := false
		for _, category := range categories {
			for _, label := range issue.Labels {
				if strings.EqualFold(label.Name, category) {
					issuesByCategory[category] = append(issuesByCategory[category], issue)
					categorized = true
					break
				}
			}
			if categorized {
				break
			}
		}
		if !categorized {
			uncategorizedIssues = append(uncategorizedIssues, issue)
		}
	}

	// Calculate progress
	var progress float64
	if totalIssues > 0 {
		progress = float64(completedIssues) / float64(totalIssues) * 100
	}

	// Generate progress bar
	progressBar := generateProgressBar(completedIssues, totalIssues, 20)

	return TemplateData{
		Milestone: milestone,
		Stats: MilestoneStats{
			TotalIssues:     totalIssues,
			CompletedIssues: completedIssues,
			Progress:        progress,
			Contributors:    contributorsList,
		},
		Categories:          categories,
		Issues:              issuesByCategory,
		UncategorizedIssues: uncategorizedIssues,
		ProgressBar:         progressBar,
	}
}

// generatePlanningContent generates the complete planning content using the template
func (m *Manager) generatePlanningContent(data TemplateData) (string, error) {
	return m.generatePlanningContentWithTime(data, time.Now())
}

// generatePlanningContentWithTime generates the complete planning content using the template with a fixed time
func (m *Manager) generatePlanningContentWithTime(data TemplateData, now time.Time) (string, error) {
	// Define template functions
	funcMap := template.FuncMap{
		"now": func() string {
			return now.UTC().Format("January 2, 2006 15:04 MST")
		},
		"formatDate": func(date *time.Time) string {
			if date == nil {
				return "No due date"
			}
			return date.Format("January 2, 2006")
		},
		"sub": func(a, b int) int {
			return a - b
		},
	}

	// Load template with functions
	tmpl, err := template.New("planning.gotmpl").Funcs(funcMap).ParseFS(templates, "templates/planning.gotmpl")
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute template
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// generateProgressBar generates a progress bar string
func generateProgressBar(completed int, total int, length int) string {
	if total == 0 {
		return strings.Repeat("░", length) + " 0%"
	}

	progress := int(float64(completed) / float64(total) * float64(length))
	percentage := int(float64(completed) / float64(total) * 100)

	filled := strings.Repeat("█", progress)
	empty := strings.Repeat("░", length-progress)

	return filled + empty + fmt.Sprintf(" %d%%", percentage)
}

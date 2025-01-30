package planning

import (
	"bufio"
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/elliotxx/osp/pkg/log"
)

//go:embed templates/planning.gotmpl
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
	HTMLURL  string  `json:"html_url"`
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
	DryRun        bool     // If true, only show the planning content without updating
	AutoConfirm   bool     // If true, skip confirmation and update automatically
	Priorities    []string // Priority labels to sort issues by, from high to low
}

// DefaultOptions returns default planning options
func DefaultOptions() Options {
	return Options{
		PlanningLabel: "planning",
		Categories:    []string{"bug", "documentation", "enhancement"},
		ExcludePR:     true,
		DryRun:        false,
		AutoConfirm:   false,
		Priorities:    []string{"priority/high", "priority/medium", "priority/low"},
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
	HighPriorityIssues  []Issue
	ProgressBar         string
	Priorities          []string
}

// askForConfirmation asks the user for confirmation
func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		log.P("?").C(log.ColorBlue).N().Log("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Error("Error reading input: %v", err)
			return false
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}

// Update updates or creates a planning issue for a milestone
func (m *Manager) Update(ctx context.Context, owner, repo string, milestoneNumber int, opts Options) error {
	log.Debug("Updating planning issue for milestone #%d in %s/%s", milestoneNumber, owner, repo)

	// Get milestone
	var milestone Milestone
	path := fmt.Sprintf("repos/%s/%s/milestones/%d", owner, repo, milestoneNumber)
	err := m.client.Get(path, &milestone)
	if err != nil {
		return fmt.Errorf("failed to get milestone: %w", err)
	}
	log.Debug("Found milestone: %s (#%d)", milestone.Title, milestone.Number)

	// Get all issues in the milestone
	var issues []Issue
	path = fmt.Sprintf("repos/%s/%s/issues?milestone=%d&state=all", owner, repo, milestoneNumber)
	err = m.client.Get(path, &issues)
	if err != nil {
		return fmt.Errorf("failed to get issues: %w", err)
	}
	log.Debug("Found %d issues in milestone", len(issues))

	// Filter out pull requests if exclude_pr is true
	if opts.ExcludePR {
		var filtered []Issue
		for _, issue := range issues {
			if !strings.Contains(issue.HTMLURL, "/pull/") {
				filtered = append(filtered, issue)
			}
		}
		issues = filtered
	}

	// Prepare data for template
	data := m.prepareTemplateData(milestone, issues, opts)

	// Generate planning content
	content, err := m.generatePlanningContent(data)
	if err != nil {
		return fmt.Errorf("failed to generate planning content: %w", err)
	}
	log.Debug("Generated planning content with %d bytes", len(content))

	// Find existing planning issue
	path = fmt.Sprintf("repos/%s/%s/issues?labels=%s&state=all", owner, repo, opts.PlanningLabel)
	var existingIssues []Issue
	err = m.client.Get(path, &existingIssues)
	if err != nil {
		return fmt.Errorf("failed to get existing planning issues: %w", err)
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

	// Show preview
	if planningIssue == nil {
		log.Info("Creating new planning issue for milestone '%s'", milestone.Title)
	} else {
		log.Info("Updating existing planning issue #%d for milestone #%d (%s)", planningIssue.Number, milestone.Number, milestone.Title)
	}

	// Preview the content
	log.C(log.ColorBlue).P("↓").Log("Preview of the planning content:")
	log.C(log.ColorCyan).Log("%s", content)

	if !opts.DryRun {
		// Ask for confirmation if auto-confirm is not enabled
		if !opts.AutoConfirm {
			// Show update target
			if planningIssue == nil {
				log.Info("Will create a new planning issue with the above content")
			} else {
				issueURL := fmt.Sprintf("https://github.com/%s/%s/issues/%d", owner, repo, planningIssue.Number)
				log.Info("Will update existing planning issue (%s) with the above content", issueURL)
			}

			if !askForConfirmation("Do you want to proceed with the update?") {
				log.Info("Update cancelled")
				return nil
			}
		} else {
			log.C(log.ColorYellow).P("!").Log("Auto-confirm is enabled, skipping confirmation")
		}

		// Create or update the planning issue
		if planningIssue == nil {
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

			path := fmt.Sprintf("repos/%s/%s/issues", owner, repo)
			var response struct {
				Number int `json:"number"`
			}
			err = m.client.Post(path, bytes.NewReader(bodyBytes), &response)
			if err != nil {
				return fmt.Errorf("failed to create planning issue: %w", err)
			}
			issueURL := fmt.Sprintf("https://github.com/%s/%s/issues/%d", owner, repo, response.Number)
			log.Success("Successfully created planning issue for milestone '%s'", milestone.Title).
				L(1).P("→").Log("Planning issue URL: %s", issueURL)
		} else {
			// Update existing issue
			body := map[string]interface{}{
				"title": planningTitle,
				"body":  content,
			}
			bodyBytes, err := json.Marshal(body)
			if err != nil {
				return fmt.Errorf("failed to marshal request body: %w", err)
			}

			path := fmt.Sprintf("repos/%s/%s/issues/%d", owner, repo, planningIssue.Number)
			err = m.client.Patch(path, bytes.NewReader(bodyBytes), nil)
			if err != nil {
				return fmt.Errorf("failed to update planning issue: %w", err)
			}
			issueURL := fmt.Sprintf("https://github.com/%s/%s/issues/%d", owner, repo, planningIssue.Number)
			log.Success("Successfully updated planning issue #%d", planningIssue.Number).
				L(1).P("→").Log("Planning issue URL: %s", issueURL)
		}
	} else {
		log.C(log.ColorYellow).P("!").Log("Dry-run mode, skipping update")
	}

	return nil
}

// prepareTemplateData prepares data for the template
func (m *Manager) prepareTemplateData(milestone Milestone, issues []Issue, opts Options) TemplateData {
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
		for _, category := range opts.Categories {
			for _, label := range issue.Labels {
				if strings.EqualFold(label.Name, category) {
					issuesByCategory[category] = append(issuesByCategory[category], issue)
					categorized = true
				}
			}
		}
		if !categorized {
			uncategorizedIssues = append(uncategorizedIssues, issue)
		}
	}

	// Sort issues in each category by priority
	for category := range issuesByCategory {
		sort.Slice(issuesByCategory[category], func(i, j int) bool {
			iPriority := getPriorityLevel(issuesByCategory[category][i].Labels, opts.Priorities)
			jPriority := getPriorityLevel(issuesByCategory[category][j].Labels, opts.Priorities)
			if iPriority != jPriority {
				return iPriority < jPriority // Lower index means higher priority
			}
			// If priorities are equal, sort by issue number
			return issuesByCategory[category][i].Number < issuesByCategory[category][j].Number
		})
	}

	// Get high priority issues (top 2 priority levels)
	var highPriorityIssues []Issue
	if len(opts.Priorities) >= 2 {
		for _, issue := range issues {
			level := getPriorityLevel(issue.Labels, opts.Priorities)
			if level < 2 { // Only include top 2 priority levels
				highPriorityIssues = append(highPriorityIssues, issue)
			}
		}
		// Sort high priority issues by priority
		sort.Slice(highPriorityIssues, func(i, j int) bool {
			iPriority := getPriorityLevel(highPriorityIssues[i].Labels, opts.Priorities)
			jPriority := getPriorityLevel(highPriorityIssues[j].Labels, opts.Priorities)
			if iPriority != jPriority {
				return iPriority < jPriority
			}
			return highPriorityIssues[i].Number < highPriorityIssues[j].Number
		})
	}

	// Calculate progress
	var progress float64
	if totalIssues > 0 {
		progress = float64(completedIssues) / float64(totalIssues) * 100
	}

	return TemplateData{
		Milestone:           milestone,
		Stats:               MilestoneStats{TotalIssues: totalIssues, CompletedIssues: completedIssues, Progress: progress, Contributors: contributorsList},
		Categories:          opts.Categories,
		Issues:              issuesByCategory,
		UncategorizedIssues: uncategorizedIssues,
		HighPriorityIssues:  highPriorityIssues,
		ProgressBar:         generateProgressBar(completedIssues, totalIssues, 20),
		Priorities:          opts.Priorities,
	}
}

// getPriorityLevel returns the priority level of an issue based on its labels
// Returns the index of the highest priority label found, or len(priorities) if no priority label is found
func getPriorityLevel(labels []Label, priorities []string) int {
	for _, label := range labels {
		for i, priority := range priorities {
			if strings.EqualFold(label.Name, priority) {
				return i
			}
		}
	}
	return len(priorities)
}

// ListOpenMilestones returns a list of open milestones for the repository
func (m *Manager) ListOpenMilestones(ctx context.Context, owner, repo string) ([]Milestone, error) {
	var milestones []Milestone
	path := fmt.Sprintf("repos/%s/%s/milestones?state=open", owner, repo)

	err := m.client.Get(path, &milestones)
	if err != nil {
		return nil, fmt.Errorf("failed to list milestones: %w", err)
	}

	return milestones, nil
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
		"getPriorityLevel": func(labels []Label) int {
			for _, label := range labels {
				for i, priority := range data.Priorities {
					if strings.EqualFold(label.Name, priority) {
						return i
					}
				}
			}
			return len(data.Priorities)
		},
		"getPriorityMark": func(level int) string {
			if level >= len(data.Priorities) {
				return ""
			}
			return strings.Repeat("!", len(data.Priorities)-level)
		},
		"getTopTwoPriorities": func() string {
			if len(data.Priorities) == 0 {
				return ""
			}
			if len(data.Priorities) == 1 {
				return fmt.Sprintf("`%s`", data.Priorities[0])
			}
			return fmt.Sprintf("`%s` and `%s`", data.Priorities[0], data.Priorities[1])
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

package onboard

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/elliotxx/osp/pkg/config"
	"github.com/elliotxx/osp/pkg/log"
	"github.com/elliotxx/osp/pkg/util/prompt"
)

//go:embed templates/*.gotmpl
var templatesFS embed.FS

// Manager manages onboarding process
type Manager struct {
	state  *config.State
	client *api.RESTClient
}

// NewManager creates a new onboarding manager
func NewManager(client *api.RESTClient) (*Manager, error) {
	state, err := config.LoadState()
	if err != nil {
		return nil, fmt.Errorf("failed to load state: %w", err)
	}

	return &Manager{
		state:  state,
		client: client,
	}, nil
}

// OnboardIssue represents an issue suitable for new contributors
type OnboardIssue struct {
	Difficulty string `json:"difficulty"` // Easy, Medium, Hard
	Status     string `json:"status"`     // open, closed
	Assignee   string `json:"assignee,omitempty"`
	Number     int    `json:"number"` // Issue number for sorting
	Category   string `json:"category"`
}

// Options represents the options for onboarding
type Options struct {
	// Issue labels configuration
	OnboardLabels    []string // Labels for identifying suitable issues for community contributions
	DifficultyLabels []string // Labels indicating the difficulty of issues
	CategoryLabels   []string // Labels for classifying issues by type

	// Target issue configuration
	TargetLabel string // Label used to locate the issue where onboarding content will be updated
	TargetTitle string // Title of the target issue where onboarding content will be updated

	// Command behavior
	DryRun      bool // If true, only show preview without making changes
	AutoConfirm bool // If true, skip confirmation prompt
}

// DefaultOptions returns the default options
func DefaultOptions() Options {
	return Options{
		// Issue labels defaults
		OnboardLabels:    []string{"help wanted", "good first issue"},
		DifficultyLabels: []string{"good first issue", "help wanted"},
		CategoryLabels:   []string{"bug", "enhancement", "documentation"},

		// Target issue defaults
		TargetLabel: "onboarding",
		TargetTitle: "Onboarding: Getting Started with Contributing",

		// Command behavior defaults
		DryRun:      false,
		AutoConfirm: false,
	}
}

// Stats represents statistics about the issues
type Stats struct {
	TotalIssues      int      `json:"total_issues"`
	CompletedIssues  int      `json:"completed_issues"`
	InProgressIssues int      `json:"in_progress_issues"`
	UnassignedIssues int      `json:"unassigned_issues"`
	Contributors     []string `json:"contributors"`
}

// TemplateData represents the data passed to the template
type TemplateData struct {
	RepoName         string                               `json:"repo_name"`
	IssuesByCategory map[string]map[string][]OnboardIssue `json:"issues_by_category"`
	DifficultyLabels []string                             `json:"difficulty_labels"`
	CategoryLabels   []string                             `json:"category_labels"`
	Stats            Stats                                `json:"stats"`
	OnboardLabels    []string                             `json:"onboard_labels"`
}

// SearchOnboardIssues generates onboarding issues for new contributors
func (m *Manager) SearchOnboardIssues(_ context.Context, repoName string, opts Options) ([]OnboardIssue, error) {
	// Split owner and repo
	parts := strings.Split(repoName, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid repository name format, should be owner/repo")
	}

	// Query issues with help wanted labels
	var query string
	query = fmt.Sprintf("repo:%s is:issue", repoName)

	// Add help labels
	if len(opts.OnboardLabels) > 0 {
		query += " label:"
		for i, label := range opts.OnboardLabels {
			if i > 0 {
				query += ","
			}
			query += fmt.Sprintf("\"%s\"", label)
		}
	}

	// Add sorting parameters
	query += " sort:updated-desc"

	log.Debug("Search query: %s", query)

	// Make API request with pagination
	var allItems []struct {
		Title   string `json:"title"`
		Number  int    `json:"number"`
		HTMLURL string `json:"html_url"`
		Labels  []struct {
			Name string `json:"name"`
		} `json:"labels"`
		State    string `json:"state"`
		Assignee *struct {
			Login string `json:"login"`
		} `json:"assignee"`
	}

	page := 1
	for {
		var response struct {
			TotalCount        int  `json:"total_count"`
			IncompleteResults bool `json:"incomplete_results"`
			Items             []struct {
				Title   string `json:"title"`
				Number  int    `json:"number"`
				HTMLURL string `json:"html_url"`
				Labels  []struct {
					Name string `json:"name"`
				} `json:"labels"`
				State    string `json:"state"`
				Assignee *struct {
					Login string `json:"login"`
				} `json:"assignee"`
			} `json:"items"`
		}

		err := m.client.Get(fmt.Sprintf("search/issues?q=%s&page=%d&per_page=100", url.QueryEscape(query), page), &response)
		if err != nil {
			return nil, fmt.Errorf("failed to search issues: %w", err)
		}

		if len(response.Items) == 0 {
			break
		}

		allItems = append(allItems, response.Items...)
		log.Debug("Found %d issues on page %d", len(response.Items), page)

		if len(response.Items) < 100 {
			break
		}

		page++
	}

	log.Debug("Found %d issues in total", len(allItems))

	// Convert issues to onboard issues
	issues := make([]OnboardIssue, 0, len(allItems))
	for _, issue := range allItems {
		// Determine difficulty level
		difficulty := ""
		for _, label := range issue.Labels {
			for _, difficultyLabel := range opts.DifficultyLabels {
				if strings.EqualFold(label.Name, difficultyLabel) {
					difficulty = difficultyLabel
					break
				}
			}
			if difficulty != "" {
				break
			}
		}

		// Determine category
		category := ""
		for _, label := range issue.Labels {
			for _, categoryLabel := range opts.CategoryLabels {
				if strings.EqualFold(label.Name, categoryLabel) {
					category = categoryLabel
					break
				}
			}
			if category != "" {
				break
			}
		}

		onboardIssue := OnboardIssue{
			Difficulty: difficulty,
			Status:     issue.State,
			Number:     issue.Number,
			Category:   category,
			Assignee: func() string {
				if issue.Assignee != nil {
					return issue.Assignee.Login
				}
				return ""
			}(),
		}
		issues = append(issues, onboardIssue)
		log.Debug("Added issue: (Difficulty: %s, Status: %s)", onboardIssue.Difficulty, onboardIssue.Status)
	}

	return issues, nil
}

// GenerateContent generates the complete content using the template
func (m *Manager) GenerateContent(issues []OnboardIssue, repoName string, opts Options) (string, error) {
	// Load template
	log.Debug("Loading template...")
	tmpl := template.New("onboard.gotmpl").Funcs(template.FuncMap{
		"now": func() string {
			return time.Now().UTC().Format("January 2, 2006 15:04 MST")
		},
		"urlEncode":           url.QueryEscape,
		"generateProgressBar": generateProgressBar,
		"add": func(a, b int) int {
			return a + b
		},
		"hasUnspecifiedIssues": func(issuesByCategory map[string]map[string][]OnboardIssue) bool {
			if categoryMap, ok := issuesByCategory[""]; ok {
				for _, issues := range categoryMap {
					if len(issues) > 0 {
						return true
					}
				}
			}
			return false
		},
	})

	tmpl, err := tmpl.ParseFS(templatesFS, "templates/*.gotmpl")
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Group issues by difficulty and category
	issuesByDiffCategory := make(map[string]map[string][]OnboardIssue)
	uniqueIssues := make([]OnboardIssue, 0)

	// Initialize maps for each difficulty level
	for _, difficultyLabel := range opts.DifficultyLabels {
		issuesByDiffCategory[difficultyLabel] = make(map[string][]OnboardIssue)
	}
	// Initialize map for unspecified difficulty (empty string)
	issuesByDiffCategory[""] = make(map[string][]OnboardIssue)

	// Create a map to track unique issues
	uniqueIssueMap := make(map[int]struct{})

	// Group issues by difficulty and category
	for _, issue := range issues {
		// Skip if we've already processed this issue
		if _, ok := uniqueIssueMap[issue.Number]; ok {
			continue
		}
		uniqueIssueMap[issue.Number] = struct{}{}
		uniqueIssues = append(uniqueIssues, issue)

		var difficultyLabel string
		var categoryLabel string

		// Find difficulty label
		if issue.Difficulty != "" {
			for _, difficultyPrefix := range opts.DifficultyLabels {
				if issue.Difficulty == difficultyPrefix {
					difficultyLabel = issue.Difficulty
					break
				}
			}
		}

		// Find category label
		if issue.Category != "" {
			for _, categoryPrefix := range opts.CategoryLabels {
				if issue.Category == categoryPrefix {
					categoryLabel = issue.Category
					break
				}
			}
		}

		// Initialize category map if not exists
		if _, ok := issuesByDiffCategory[difficultyLabel][categoryLabel]; !ok {
			issuesByDiffCategory[difficultyLabel][categoryLabel] = make([]OnboardIssue, 0)
		}

		// Add issue to the appropriate category
		issuesByDiffCategory[difficultyLabel][categoryLabel] = append(issuesByDiffCategory[difficultyLabel][categoryLabel], issue)
	}

	// Sort issues within each category by status (open before closed) and number
	for difficulty := range issuesByDiffCategory {
		for category := range issuesByDiffCategory[difficulty] {
			sort.Slice(issuesByDiffCategory[difficulty][category], func(i, j int) bool {
				// First sort by status (open before closed)
				if issuesByDiffCategory[difficulty][category][i].Status != issuesByDiffCategory[difficulty][category][j].Status {
					return issuesByDiffCategory[difficulty][category][i].Status == "open"
				}
				// Then sort by number
				return issuesByDiffCategory[difficulty][category][i].Number < issuesByDiffCategory[difficulty][category][j].Number
			})
		}
	}

	// Calculate statistics
	stats := Stats{
		TotalIssues:      len(uniqueIssues),
		CompletedIssues:  0,
		InProgressIssues: 0,
		UnassignedIssues: 0,
	}

	// Use map to track unique contributors
	contributors := make(map[string]struct{})

	for _, issue := range uniqueIssues {
		switch {
		case issue.Status == "closed":
			stats.CompletedIssues++
			if issue.Assignee != "" {
				contributors[issue.Assignee] = struct{}{}
			}
		case issue.Assignee != "":
			stats.InProgressIssues++
		default:
			stats.UnassignedIssues++
		}
	}

	for contributor := range contributors {
		stats.Contributors = append(stats.Contributors, contributor)
	}
	sort.Strings(stats.Contributors)

	// Create a buffer to store the output
	var buf strings.Builder

	// Execute template
	log.Debug("Executing template...")
	err = tmpl.ExecuteTemplate(&buf, "onboard.gotmpl", TemplateData{
		RepoName:         repoName,
		IssuesByCategory: issuesByDiffCategory,
		DifficultyLabels: opts.DifficultyLabels, // 不包含空字符串，让模版决定何时显示未指定难度的 issue
		CategoryLabels:   opts.CategoryLabels,
		Stats:            stats,
		OnboardLabels:    opts.OnboardLabels,
	})
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// generateProgressBar generates a progress bar string based on completion percentage
func generateProgressBar(completed, total int) string {
	const width = 20 // Total width of the progress bar
	if total == 0 {
		return strings.Repeat("░", width) // Empty progress bar
	}

	percentage := float64(completed) / float64(total)
	filledWidth := int(percentage * float64(width))

	// Ensure at least one block is filled if there's any progress
	if completed > 0 && filledWidth == 0 {
		filledWidth = 1
	}

	// Ensure we don't exceed the width
	if filledWidth > width {
		filledWidth = width
	}

	filled := strings.Repeat("█", filledWidth)
	empty := strings.Repeat("░", width-filledWidth)

	return filled + empty + fmt.Sprintf(" %.1f%%", percentage*100)
}

// Update updates or creates an onboarding issue
func (m *Manager) Update(ctx context.Context, repoName string, opts Options) error {
	log.Debug("Updating onboarding issue in %s", repoName)

	// Generate onboarding content
	issues, err := m.SearchOnboardIssues(ctx, repoName, opts)
	if err != nil {
		return fmt.Errorf("failed to search onboarding issues: %w", err)
	}

	content, err := m.GenerateContent(issues, repoName, opts)
	if err != nil {
		return fmt.Errorf("failed to generate onboarding content: %w", err)
	}
	log.Debug("Generated onboarding content with %d bytes", len(content))

	// Find existing onboarding issues
	path := fmt.Sprintf("repos/%s/issues?labels=%s&state=all", repoName, opts.TargetLabel)
	var existingIssues []struct {
		Title  string `json:"title"`
		Number int    `json:"number"`
	}
	err = m.client.Get(path, &existingIssues)
	if err != nil {
		return fmt.Errorf("failed to get existing onboarding issues: %w", err)
	}
	log.Debug("Found %d existing issues with onboarding label", len(existingIssues))

	// Find the onboarding issue with the smallest number
	var onboardingIssue *struct {
		Title  string `json:"title"`
		Number int    `json:"number"`
	}
	if len(existingIssues) > 0 {
		onboardingIssue = &existingIssues[0]
		for i := 1; i < len(existingIssues); i++ {
			if existingIssues[i].Number < onboardingIssue.Number {
				onboardingIssue = &existingIssues[i]
			}
		}
		log.Debug("Found onboarding issue #%d", onboardingIssue.Number)
	}

	if len(existingIssues) > 1 {
		log.Warn("Found multiple onboarding issues, will update issue #%d", onboardingIssue.Number)
	}

	// Show preview
	if onboardingIssue == nil {
		log.Info("Creating new onboarding issue")
	} else {
		log.Info("Updating existing onboarding issue #%d", onboardingIssue.Number)
	}

	// Preview the content
	log.C(log.ColorBlue).P("↓").Log("Preview of the onboarding content:")
	log.C(log.ColorCyan).Log("%s", content)

	if !opts.DryRun {
		// Ask for confirmation if auto-confirm is not enabled
		if !opts.AutoConfirm {
			// Show update target
			if onboardingIssue == nil {
				log.Info("Will create a new onboarding issue with the above content")
			} else {
				issueURL := fmt.Sprintf("https://github.com/%s/issues/%d", repoName, onboardingIssue.Number)
				log.Info("Will update existing onboarding issue (%s) with the above content", issueURL)
			}

			confirmed, err := prompt.AskForConfirmation("Do you want to proceed with the update?")
			if err != nil {
				return err
			}
			if !confirmed {
				log.Info("Update cancelled")
				return nil
			}
		} else {
			log.Warn("Auto-confirm is enabled, skipping confirmation")
		}

		// Create or update the onboarding issue
		if onboardingIssue == nil {
			// Create new issue
			body := map[string]interface{}{
				"title":  opts.TargetTitle,
				"body":   content,
				"labels": []string{opts.TargetLabel},
			}
			bodyBytes, err := json.Marshal(body)
			if err != nil {
				return fmt.Errorf("failed to marshal request body: %w", err)
			}

			path := fmt.Sprintf("repos/%s/issues", repoName)
			var response struct {
				Number int `json:"number"`
			}
			err = m.client.Post(path, bytes.NewReader(bodyBytes), &response)
			if err != nil {
				return fmt.Errorf("failed to create onboarding issue: %w", err)
			}
			issueURL := fmt.Sprintf("https://github.com/%s/issues/%d", repoName, response.Number)
			log.Success("Successfully created onboarding issue").
				L(1).P("→").Log("Onboarding issue URL: %s", issueURL)
		} else {
			// Update existing issue
			body := map[string]interface{}{
				"title": opts.TargetTitle,
				"body":  content,
			}
			bodyBytes, err := json.Marshal(body)
			if err != nil {
				return fmt.Errorf("failed to marshal request body: %w", err)
			}

			path := fmt.Sprintf("repos/%s/issues/%d", repoName, onboardingIssue.Number)
			err = m.client.Patch(path, bytes.NewReader(bodyBytes), nil)
			if err != nil {
				return fmt.Errorf("failed to update onboarding issue: %w", err)
			}
			issueURL := fmt.Sprintf("https://github.com/%s/issues/%d", repoName, onboardingIssue.Number)
			log.Success("Successfully updated onboarding issue #%d", onboardingIssue.Number).
				L(1).P("→").Log("Onboarding issue URL: %s", issueURL)
		}
	} else {
		log.Warn("Dry-run mode, skipping update")
	}

	return nil
}

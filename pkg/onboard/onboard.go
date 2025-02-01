package onboard

import (
	"context"
	"embed"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"text/template"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/elliotxx/osp/pkg/config"
	"github.com/elliotxx/osp/pkg/log"
)

//go:embed templates/*.gotmpl
var templatesFS embed.FS

// Manager handles onboarding issue management
type Manager struct {
	cfg *config.Config
}

// OnboardIssue represents an issue suitable for new contributors
type OnboardIssue struct {
	Difficulty string `json:"difficulty"` // Easy, Medium, Hard
	Status     string `json:"status"`     // open, closed
	Assignee   string `json:"assignee,omitempty"`
	Number     int    `json:"number"` // Issue number for sorting
	Category   string `json:"category"`
}

// Options represents the options for generating onboarding issues
type Options struct {
	HelpLabels       []string `json:"help_labels"`
	DifficultyLabels []string `json:"difficulty_labels"`
	Categories       []string `json:"categories"`
}

var (
	// defaultHelpLabels is the default help labels
	defaultHelpLabels = []string{"good first issue", "help wanted"}
)

// TemplateData represents the data passed to the template
type TemplateData struct {
	RepoName         string                               `json:"repo_name"`
	IssuesByCategory map[string]map[string][]OnboardIssue `json:"issues_by_category"`
	DifficultyLabels []string                             `json:"difficulty_labels"`
	Categories       []string                             `json:"categories"`
}

// NewManager creates a new onboard manager
func NewManager(cfg *config.Config) *Manager {
	return &Manager{cfg: cfg}
}

// SearchOnboardIssues generates onboarding issues for new contributors
func (m *Manager) SearchOnboardIssues(_ context.Context, repoName string, opts Options) ([]OnboardIssue, error) {
	// Split owner and repo
	parts := strings.Split(repoName, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid repository name format, should be owner/repo")
	}

	// Create GitHub client
	client, err := api.DefaultRESTClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create GitHub client: %w", err)
	}

	// Query issues with help wanted labels
	var query string
	query = fmt.Sprintf("repo:%s is:issue", repoName)

	// Add help labels
	if len(opts.HelpLabels) > 0 {
		query += " label:"
		for i, label := range opts.HelpLabels {
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

		err = client.Get(fmt.Sprintf("search/issues?q=%s&page=%d&per_page=100", url.QueryEscape(query), page), &response)
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
			for _, categoryLabel := range opts.Categories {
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
	tmpl, err := template.ParseFS(templatesFS, "templates/*.gotmpl")
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Sort issues by number to ensure consistent ordering
	log.Debug("Sorting issues...")
	sort.Slice(issues, func(i, j int) bool {
		return issues[i].Number < issues[j].Number
	})

	// Group issues by difficulty and category, ensuring uniqueness
	log.Debug("Grouping issues...")
	// Use map to ensure issue uniqueness
	uniqueIssues := make(map[int]OnboardIssue)
	for _, issue := range issues {
		uniqueIssues[issue.Number] = issue
	}

	// Create a map of difficulty -> category -> issues
	issuesByDiffCategory := make(map[string]map[string][]OnboardIssue)
	for _, issue := range uniqueIssues {
		// Initialize maps if not exist
		if _, ok := issuesByDiffCategory[issue.Difficulty]; !ok {
			issuesByDiffCategory[issue.Difficulty] = make(map[string][]OnboardIssue)
		}
		issuesByDiffCategory[issue.Difficulty][issue.Category] = append(issuesByDiffCategory[issue.Difficulty][issue.Category], issue)
	}

	// Sort issues within each category
	for difficulty := range issuesByDiffCategory {
		for category := range issuesByDiffCategory[difficulty] {
			sort.Slice(issuesByDiffCategory[difficulty][category], func(i, j int) bool {
				return issuesByDiffCategory[difficulty][category][i].Number < issuesByDiffCategory[difficulty][category][j].Number
			})
		}
	}

	// Create a buffer to store the output
	var buf strings.Builder

	// Execute template
	log.Debug("Executing template...")
	err = tmpl.ExecuteTemplate(&buf, "onboard.gotmpl", TemplateData{
		RepoName:         repoName,
		IssuesByCategory: issuesByDiffCategory,
		DifficultyLabels: opts.DifficultyLabels,
		Categories:       opts.Categories,
	})
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

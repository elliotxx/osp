package plan

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/elliotxx/osp/internal/config"
)

// Manager handles project planning
type Manager struct {
	cfg    *config.Config
	client *http.Client
}

// Plan represents a project plan
type Plan struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Milestones  []Milestone `json:"milestones"`
}

// Milestone represents a project milestone
type Milestone struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	State       string    `json:"state"`
	Issues      []Issue   `json:"issues"`
}

// Issue represents a GitHub issue
type Issue struct {
	Number      int       `json:"number"`
	Title       string    `json:"title"`
	State       string    `json:"state"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Labels      []string  `json:"labels"`
	Assignees   []string  `json:"assignees"`
}

// NewManager creates a new plan manager
func NewManager(cfg *config.Config) *Manager {
	return &Manager{
		cfg:    cfg,
		client: http.DefaultClient,
	}
}

// Generate generates a project plan
func (m *Manager) Generate(ctx context.Context, repoName string) (*Plan, error) {
	// Get milestones
	milestones, err := m.getMilestones(ctx, repoName)
	if err != nil {
		return nil, fmt.Errorf("failed to get milestones: %w", err)
	}

	// Create plan based on milestones
	plan := &Plan{
		Title:       fmt.Sprintf("%s Project Plan", repoName),
		Description: "Automatically generated project plan",
		StartDate:   time.Now(),
		EndDate:     time.Now().AddDate(0, 3, 0), // Default to 3 months
		Milestones:  milestones,
	}

	return plan, nil
}

// Update updates the project plan
func (m *Manager) Update(ctx context.Context, repoName string, plan *Plan) error {
	// Update milestones on GitHub
	for _, milestone := range plan.Milestones {
		if err := m.updateMilestone(ctx, repoName, &milestone); err != nil {
			return fmt.Errorf("failed to update milestone: %w", err)
		}
	}

	return nil
}

// getMilestones fetches milestones from GitHub
func (m *Manager) getMilestones(ctx context.Context, repoName string) ([]Milestone, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/milestones?state=all", repoName)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	if m.cfg.Auth.Token != "" {
		req.Header.Set("Authorization", "token "+m.cfg.Auth.Token)
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var milestones []struct {
		Title       string    `json:"title"`
		Description string    `json:"description"`
		DueOn       time.Time `json:"due_on"`
		State       string    `json:"state"`
		Number      int       `json:"number"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&milestones); err != nil {
		return nil, err
	}

	result := make([]Milestone, 0, len(milestones))
	for _, milestone := range milestones {
		// Get issues for this milestone
		issues, err := m.getIssues(ctx, repoName, milestone.Number)
		if err != nil {
			return nil, err
		}

		result = append(result, Milestone{
			Title:       milestone.Title,
			Description: milestone.Description,
			DueDate:     milestone.DueOn,
			State:       milestone.State,
			Issues:      issues,
		})
	}

	return result, nil
}

// getIssues fetches issues for a milestone
func (m *Manager) getIssues(ctx context.Context, repoName string, milestoneNumber int) ([]Issue, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/issues?milestone=%d&state=all", repoName, milestoneNumber)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	if m.cfg.Auth.Token != "" {
		req.Header.Set("Authorization", "token "+m.cfg.Auth.Token)
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var issues []struct {
		Number    int       `json:"number"`
		Title     string    `json:"title"`
		State     string    `json:"state"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Labels    []struct {
			Name string `json:"name"`
		} `json:"labels"`
		Assignees []struct {
			Login string `json:"login"`
		} `json:"assignees"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&issues); err != nil {
		return nil, err
	}

	result := make([]Issue, 0, len(issues))
	for _, i := range issues {
		labels := make([]string, 0, len(i.Labels))
		for _, l := range i.Labels {
			labels = append(labels, l.Name)
		}

		assignees := make([]string, 0, len(i.Assignees))
		for _, a := range i.Assignees {
			assignees = append(assignees, a.Login)
		}

		result = append(result, Issue{
			Number:    i.Number,
			Title:     i.Title,
			State:     i.State,
			CreatedAt: i.CreatedAt,
			UpdatedAt: i.UpdatedAt,
			Labels:    labels,
			Assignees: assignees,
		})
	}

	return result, nil
}

// updateMilestone updates a milestone on GitHub
func (m *Manager) updateMilestone(ctx context.Context, repoName string, milestone *Milestone) error {
	// Implementation for updating milestone
	// This would involve making PATCH requests to GitHub API
	// Left as an exercise for actual implementation
	return nil
}

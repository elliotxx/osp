package task

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/elliotxx/osp/internal/config"
)

// Manager handles task management
type Manager struct {
	cfg    *config.Config
	client *http.Client
}

// Task represents a GitHub issue task
type Task struct {
	Number      int       `json:"number"`
	Title       string    `json:"title"`
	Body        string    `json:"body"`
	State       string    `json:"state"`
	Labels      []string  `json:"labels"`
	Assignees   []string  `json:"assignees"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Difficulty  string    `json:"difficulty"`
	Type        string    `json:"type"`
}

// NewManager creates a new task manager
func NewManager(cfg *config.Config) *Manager {
	return &Manager{
		cfg:    cfg,
		client: http.DefaultClient,
	}
}

// List returns tasks based on filters
func (m *Manager) List(ctx context.Context, repoName string, taskType string) ([]Task, error) {
	// Convert task type to labels
	labels := m.getLabelsForType(taskType)
	
	// Build URL with query parameters
	url := fmt.Sprintf("https://api.github.com/repos/%s/issues?state=open", repoName)
	if len(labels) > 0 {
		url += "&labels=" + strings.Join(labels, ",")
	}

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
		Body      string    `json:"body"`
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

	tasks := make([]Task, 0, len(issues))
	for _, i := range issues {
		labels := make([]string, 0, len(i.Labels))
		for _, l := range i.Labels {
			labels = append(labels, l.Name)
		}

		assignees := make([]string, 0, len(i.Assignees))
		for _, a := range i.Assignees {
			assignees = append(assignees, a.Login)
		}

		task := Task{
			Number:    i.Number,
			Title:     i.Title,
			Body:      i.Body,
			State:     i.State,
			Labels:    labels,
			Assignees: assignees,
			CreatedAt: i.CreatedAt,
			UpdatedAt: i.UpdatedAt,
			Type:      m.getTypeFromLabels(labels),
			Difficulty: m.getDifficultyFromLabels(labels),
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// Generate generates new tasks based on repository analysis
func (m *Manager) Generate(ctx context.Context, repoName string, taskType string) error {
	// This is a placeholder for task generation logic
	// In a real implementation, you would:
	// 1. Analyze the repository (code, issues, PRs)
	// 2. Identify areas that need work
	// 3. Generate appropriate tasks
	// 4. Create issues on GitHub
	
	// Example task generation for documentation
	if taskType == "good-first-issue" {
		// Find files without documentation
		// Create issues for adding documentation
	}

	return nil
}

// getLabelsForType converts task type to GitHub labels
func (m *Manager) getLabelsForType(taskType string) []string {
	switch taskType {
	case "good-first-issue":
		return []string{"good first issue", "help wanted"}
	case "help-wanted":
		return []string{"help wanted"}
	case "bug":
		return []string{"bug"}
	case "enhancement":
		return []string{"enhancement"}
	default:
		return nil
	}
}

// getTypeFromLabels determines task type from labels
func (m *Manager) getTypeFromLabels(labels []string) string {
	for _, label := range labels {
		switch label {
		case "good first issue":
			return "good-first-issue"
		case "help wanted":
			return "help-wanted"
		case "bug":
			return "bug"
		case "enhancement":
			return "enhancement"
		}
	}
	return "other"
}

// getDifficultyFromLabels determines task difficulty from labels
func (m *Manager) getDifficultyFromLabels(labels []string) string {
	for _, label := range labels {
		switch label {
		case "difficulty/easy":
			return "easy"
		case "difficulty/medium":
			return "medium"
		case "difficulty/hard":
			return "hard"
		}
	}
	return "unknown"
}

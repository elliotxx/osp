package plan

import (
	"context"
	"fmt"
	"time"

	"github.com/elliotxx/osp/internal/config"
)

// Manager handles project planning
type Manager struct {
	cfg *config.Config
}

// Milestone represents a project milestone
type Milestone struct {
	Title     string  `json:"title"`
	DueDate   string  `json:"due_date"`
	Progress  int     `json:"progress"`
	Issues    []Issue `json:"issues,omitempty"`
}

// Issue represents a project issue
type Issue struct {
	Title    string `json:"title"`
	Number   int    `json:"number"`
	Status   string `json:"status"`
	Assignee string `json:"assignee,omitempty"`
}

// Plan represents a project plan
type Plan struct {
	Name      string      `json:"name"`
	CreatedAt string      `json:"created_at"`
	Status    string      `json:"status"`
	Progress  int         `json:"progress"`
	Milestones []Milestone `json:"milestones"`
}

// NewManager creates a new plan manager
func NewManager(cfg *config.Config) *Manager {
	return &Manager{cfg: cfg}
}

// Generate generates a project plan
func (m *Manager) Generate(ctx context.Context, repoName string, includeIssues bool) (*Plan, error) {
	// TODO: Implement actual plan generation logic
	// This is just a mock implementation
	plan := &Plan{
		Name:      fmt.Sprintf("%s Plan", repoName),
		CreatedAt: time.Now().Format(time.RFC3339),
		Status:    "Active",
		Progress:  50,
		Milestones: []Milestone{
			{
				Title:    "First Milestone",
				DueDate:  time.Now().AddDate(0, 1, 0).Format("2006-01-02"),
				Progress: 75,
				Issues: []Issue{
					{
						Title:    "Implement feature X",
						Number:   1,
						Status:   "Open",
						Assignee: "user1",
					},
				},
			},
		},
	}

	if !includeIssues {
		for i := range plan.Milestones {
			plan.Milestones[i].Issues = nil
		}
	}

	return plan, nil
}

// List returns all project plans
func (m *Manager) List(ctx context.Context, repoName string) ([]Plan, error) {
	// TODO: Implement actual plan listing logic
	// This is just a mock implementation
	return []Plan{
		{
			Name:      fmt.Sprintf("%s Plan", repoName),
			CreatedAt: time.Now().Format(time.RFC3339),
			Status:    "Active",
			Progress:  50,
		},
	}, nil
}

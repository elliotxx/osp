package stats

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/elliotxx/osp/pkg/config"
)

// Manager handles repository statistics
type Manager struct {
	cfg    *config.Config
	client *http.Client
}

// Stats represents repository statistics
type Stats struct {
	Stars      int    `json:"stars"`
	Forks      int    `json:"forks"`
	OpenIssues int    `json:"open_issues"`
	LastUpdated string `json:"last_updated"`
}

// StarHistory represents star count at a specific date
type StarHistory struct {
	Date  time.Time `json:"date"`
	Stars int       `json:"stars"`
}

// NewManager creates a new stats manager
func NewManager(cfg *config.Config) *Manager {
	return &Manager{
		cfg:    cfg,
		client: http.DefaultClient,
	}
}

// Get returns repository statistics
func (m *Manager) Get(ctx context.Context, repoName string) (*Stats, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s", repoName)
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

	var repo struct {
		StargazersCount int       `json:"stargazers_count"`
		ForksCount     int       `json:"forks_count"`
		OpenIssues     int       `json:"open_issues_count"`
		UpdatedAt      time.Time `json:"updated_at"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&repo); err != nil {
		return nil, err
	}

	return &Stats{
		Stars:      repo.StargazersCount,
		Forks:      repo.ForksCount,
		OpenIssues: repo.OpenIssues,
		LastUpdated: repo.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// GetStarHistory returns star history for the specified number of days
func (m *Manager) GetStarHistory(ctx context.Context, repoName string, days int) ([]StarHistory, error) {
	// Calculate time range
	now := time.Now()
	from := now.AddDate(0, 0, -days)

	// Get current stars
	stats, err := m.Get(ctx, repoName)
	if err != nil {
		return nil, err
	}

	// For now, we'll generate mock data since GitHub API doesn't provide historical data
	// In a real implementation, you would need to use GitHub Archive or similar service
	history := make([]StarHistory, days+1)
	currentStars := stats.Stars
	starsPerDay := currentStars / (days + 1)

	for i := 0; i <= days; i++ {
		date := from.AddDate(0, 0, i)
		stars := starsPerDay * i
		if i == days {
			stars = currentStars // Make sure the last day matches current stars
		}

		history[i] = StarHistory{
			Date:  date,
			Stars: stars,
		}
	}

	return history, nil
}

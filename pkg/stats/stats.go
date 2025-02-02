package stats

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/elliotxx/osp/pkg/auth"
	"github.com/elliotxx/osp/pkg/config"
)

// Manager manages repository statistics
type Manager struct {
	state  *config.State
	client *http.Client
}

// Stats represents repository statistics
type Stats struct {
	Stars       int    `json:"stars"`
	Forks       int    `json:"forks"`
	OpenIssues  int    `json:"open_issues"`
	LastUpdated string `json:"last_updated"`
}

// StarHistory represents star count at a specific date
type StarHistory struct {
	Date  time.Time `json:"date"`
	Stars int       `json:"stars"`
}

// StarEvent represents a GitHub star event
type StarEvent struct {
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

// NewManager creates a new stats manager
func NewManager() (*Manager, error) {
	state, err := config.LoadState()
	if err != nil {
		return nil, fmt.Errorf("failed to load state: %w", err)
	}

	return &Manager{
		state:  state,
		client: &http.Client{},
	}, nil
}

// Get returns repository statistics
func (m *Manager) Get(ctx context.Context, repoName string) (*Stats, error) {
	// Get token
	token, err := auth.GetToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	// Create request
	url := fmt.Sprintf("https://api.github.com/repos/%s", repoName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))

	// Send request
	resp, err := m.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse response
	var data struct {
		Stars      int    `json:"stargazers_count"`
		Forks      int    `json:"forks_count"`
		OpenIssues int    `json:"open_issues_count"`
		UpdatedAt  string `json:"updated_at"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &Stats{
		Stars:       data.Stars,
		Forks:       data.Forks,
		OpenIssues:  data.OpenIssues,
		LastUpdated: data.UpdatedAt,
	}, nil
}

// GetStarHistory returns star history for the specified number of days
func (m *Manager) GetStarHistory(ctx context.Context, repoName string, days int) ([]StarHistory, error) {
	// Get token
	token, err := auth.GetToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	// Calculate time range
	now := time.Now()
	from := now.AddDate(0, 0, -days)

	// Get current stars
	stats, err := m.Get(ctx, repoName)
	if err != nil {
		return nil, err
	}
	currentStars := stats.Stars

	// Get star events
	events, err := m.getStarEvents(ctx, repoName, token, from)
	if err != nil {
		return nil, err
	}

	// Create daily star counts
	history := make([]StarHistory, days+1)
	starsByDate := make(map[string]int)

	// Initialize with current stars
	for i := 0; i <= days; i++ {
		date := from.AddDate(0, 0, i)
		dateStr := date.Format("2006-01-02")
		starsByDate[dateStr] = currentStars
	}

	// Process star events backwards
	for _, event := range events {
		if event.Type == "WatchEvent" { // WatchEvent is GitHub's term for starring
			// Decrease star count for all dates before this event
			for i := 0; i <= days; i++ {
				d := from.AddDate(0, 0, i)
				if d.Before(event.CreatedAt) {
					dateStr := d.Format("2006-01-02")
					starsByDate[dateStr]--
				}
			}
		}
	}

	// Convert map to slice
	for i := 0; i <= days; i++ {
		date := from.AddDate(0, 0, i)
		dateStr := date.Format("2006-01-02")
		history[i] = StarHistory{
			Date:  date,
			Stars: starsByDate[dateStr],
		}
	}

	return history, nil
}

// getStarEvents returns star events for a repository
func (m *Manager) getStarEvents(ctx context.Context, repoName, token string, from time.Time) ([]StarEvent, error) {
	var events []StarEvent
	page := 1
	perPage := 100

	owner, repo, _ := strings.Cut(repoName, "/")
	baseURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/events", owner, repo)

	for {
		// Create request
		url := fmt.Sprintf("%s?page=%d&per_page=%d", baseURL, page, perPage)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Add headers
		req.Header.Set("Accept", "application/vnd.github.v3+json")
		req.Header.Set("Authorization", fmt.Sprintf("token %s", token))

		// Send request
		resp, err := m.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to send request: %w", err)
		}
		defer resp.Body.Close()

		// Check status code
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		// Parse response
		var pageEvents []StarEvent
		if err := json.NewDecoder(resp.Body).Decode(&pageEvents); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}

		// Check if we've reached events before our cutoff date
		reachedEnd := false
		for _, event := range pageEvents {
			if event.CreatedAt.Before(from) {
				reachedEnd = true
				break
			}
			if event.Type == "WatchEvent" {
				events = append(events, event)
			}
		}

		if reachedEnd || len(pageEvents) < perPage {
			break
		}

		page++
	}

	return events, nil
}

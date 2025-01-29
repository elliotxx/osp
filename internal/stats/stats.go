package stats

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/elliotxx/osp/internal/config"
)

// Manager handles repository statistics
type Manager struct {
	cfg    *config.Config
	client *http.Client
}

// Stats represents repository statistics
type Stats struct {
	Stars       int       `json:"stars"`
	Forks      int       `json:"forks"`
	Issues     int       `json:"issues"`
	PRs        int       `json:"prs"`
	LastUpdate time.Time `json:"last_update"`
}

// StarHistory represents star history data
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

// GetStats returns repository statistics
func (m *Manager) GetStats(ctx context.Context, repoName string) (*Stats, error) {
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

	var data struct {
		Stars       int       `json:"stargazers_count"`
		Forks      int       `json:"forks_count"`
		Issues     int       `json:"open_issues_count"`
		LastUpdate time.Time `json:"updated_at"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	// Get PR count
	prs, err := m.getPRCount(ctx, repoName)
	if err != nil {
		return nil, err
	}

	return &Stats{
		Stars:       data.Stars,
		Forks:      data.Forks,
		Issues:     data.Issues,
		PRs:        prs,
		LastUpdate: data.LastUpdate,
	}, nil
}

// GetStarHistory returns repository star history
func (m *Manager) GetStarHistory(ctx context.Context, repoName string, from, to time.Time) ([]StarHistory, error) {
	// Note: GitHub's API doesn't provide direct star history
	// This is a simplified implementation that uses traffic data
	// In a real implementation, you might want to use GitHub Archive or store historical data
	
	url := fmt.Sprintf("https://api.github.com/repos/%s/traffic/views", repoName)
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

	var data struct {
		Views []struct {
			Timestamp time.Time `json:"timestamp"`
			Count     int      `json:"count"`
		} `json:"views"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	// Convert traffic data to star history
	// This is just an example, in reality you'd want to store actual star events
	history := make([]StarHistory, 0, len(data.Views))
	currentStars := 0
	
	for _, view := range data.Views {
		if view.Timestamp.Before(from) || view.Timestamp.After(to) {
			continue
		}
		
		// Simulate star changes based on traffic
		// This is not accurate and should be replaced with real data
		starChange := view.Count / 100
		currentStars += starChange
		
		history = append(history, StarHistory{
			Date:  view.Timestamp,
			Stars: currentStars,
		})
	}

	return history, nil
}

// getPRCount returns the count of open pull requests
func (m *Manager) getPRCount(ctx context.Context, repoName string) (int, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/pulls?state=open&per_page=1", repoName)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, err
	}

	if m.cfg.Auth.Token != "" {
		req.Header.Set("Authorization", "token "+m.cfg.Auth.Token)
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	// Get total count from Link header
	if link := resp.Header.Get("Link"); link != "" {
		if last := resp.Header.Get("X-Total-Count"); last != "" {
			var count int
			fmt.Sscanf(last, "%d", &count)
			return count, nil
		}
	}

	return 0, nil
}

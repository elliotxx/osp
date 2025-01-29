package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/elliotxx/osp/internal/config"
)

// Manager handles repository management
type Manager struct {
	cfg    *config.Config
	client *http.Client
}

// Repository represents a GitHub repository
type Repository struct {
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Stars       int    `json:"stargazers_count"`
	Forks       int    `json:"forks_count"`
	Private     bool   `json:"private"`
}

// NewManager creates a new repository manager
func NewManager(cfg *config.Config) *Manager {
	return &Manager{
		cfg:    cfg,
		client: http.DefaultClient,
	}
}

// Add adds a repository to manage
func (m *Manager) Add(ctx context.Context, repoName string) error {
	// Validate repository name
	parts := strings.Split(repoName, "/")
	if len(parts) != 2 {
		return fmt.Errorf("invalid repository name format, expected 'owner/repo'")
	}

	// Check if repository exists on GitHub
	repo, err := m.getRepoInfo(ctx, repoName)
	if err != nil {
		return fmt.Errorf("failed to get repository info: %w", err)
	}

	// Check if repository is already managed
	for _, r := range m.cfg.Repos {
		if r.Name == repoName {
			return fmt.Errorf("repository %s is already managed", repoName)
		}
	}

	// Add repository to config
	m.cfg.Repos = append(m.cfg.Repos, config.RepoConfig{
		Name:   repoName,
		Alias:  parts[1],
		Config: make(map[string]interface{}),
	})

	// Set as current if it's the first repository
	if len(m.cfg.Repos) == 1 {
		m.cfg.Current = repoName
	}

	// Save config
	if err := m.cfg.Save(""); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// Remove removes a repository from management
func (m *Manager) Remove(repoName string) error {
	found := false
	newRepos := make([]config.RepoConfig, 0, len(m.cfg.Repos)-1)
	
	for _, r := range m.cfg.Repos {
		if r.Name == repoName {
			found = true
			continue
		}
		newRepos = append(newRepos, r)
	}

	if !found {
		return fmt.Errorf("repository %s is not managed", repoName)
	}

	m.cfg.Repos = newRepos

	// Update current repository if needed
	if m.cfg.Current == repoName {
		if len(m.cfg.Repos) > 0 {
			m.cfg.Current = m.cfg.Repos[0].Name
		} else {
			m.cfg.Current = ""
		}
	}

	// Save config
	if err := m.cfg.Save(""); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// List returns all managed repositories
func (m *Manager) List() []config.RepoConfig {
	return m.cfg.Repos
}

// Switch changes the current repository
func (m *Manager) Switch(repoName string) error {
	for _, r := range m.cfg.Repos {
		if r.Name == repoName {
			m.cfg.Current = repoName
			if err := m.cfg.Save(""); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}
			return nil
		}
	}
	return fmt.Errorf("repository %s is not managed", repoName)
}

// Current returns the current repository
func (m *Manager) Current() string {
	return m.cfg.Current
}

// getRepoInfo fetches repository information from GitHub
func (m *Manager) getRepoInfo(ctx context.Context, repoName string) (*Repository, error) {
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

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("repository %s not found", repoName)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var repo Repository
	if err := json.NewDecoder(resp.Body).Decode(&repo); err != nil {
		return nil, err
	}

	return &repo, nil
}

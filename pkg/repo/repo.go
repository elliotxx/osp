package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/elliotxx/osp/pkg/config"
)

// Manager handles repository operations
type Manager struct {
	cfg    *config.Config
	client *http.Client
}

// Repository represents a GitHub repository
type Repository struct {
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
	Fork        bool   `json:"fork"`
	Stars       int    `json:"stargazers_count"`
	Forks       int    `json:"forks_count"`
	Issues      int    `json:"open_issues_count"`
	UpdatedAt   string `json:"updated_at"`
}

// NewManager creates a new repository manager
func NewManager(cfg *config.Config) *Manager {
	return &Manager{
		cfg:    cfg,
		client: http.DefaultClient,
	}
}

// Add adds a repository to the config
func (m *Manager) Add(ctx context.Context, repoName string) error {
	// Verify repository exists
	repo, err := m.getRepository(ctx, repoName)
	if err != nil {
		return fmt.Errorf("failed to verify repository: %w", err)
	}

	// Add to config
	m.cfg.Repositories = append(m.cfg.Repositories, repo.FullName)
	if err := m.cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// Remove removes a repository from the config
func (m *Manager) Remove(repoName string) error {
	// Find and remove repository
	for i, r := range m.cfg.Repositories {
		if r == repoName {
			m.cfg.Repositories = append(m.cfg.Repositories[:i], m.cfg.Repositories[i+1:]...)
			break
		}
	}

	// Update current if needed
	if m.cfg.Current == repoName {
		m.cfg.Current = ""
	}

	// Save config
	if err := m.cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// List returns all repositories in the config
func (m *Manager) List() []string {
	return m.cfg.Repositories
}

// Switch sets the current repository
func (m *Manager) Switch(repoName string) error {
	// Verify repository is in config
	found := false
	for _, r := range m.cfg.Repositories {
		if r == repoName {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("repository %s not found in config", repoName)
	}

	// Update current
	m.cfg.Current = repoName
	if err := m.cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// Current returns the current repository
func (m *Manager) Current() string {
	return m.cfg.Current
}

// getRepository fetches repository information from GitHub
func (m *Manager) getRepository(ctx context.Context, repoName string) (*Repository, error) {
	// Split owner/repo
	parts := strings.Split(repoName, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid repository name: %s", repoName)
	}

	// Make request
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", parts[0], parts[1])
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// Add auth header if token exists
	if token := m.cfg.Auth.Token; token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var repo Repository
	if err := json.NewDecoder(resp.Body).Decode(&repo); err != nil {
		return nil, err
	}

	return &repo, nil
}

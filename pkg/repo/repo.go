package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
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

	// Automatically select the newly added repository
	m.cfg.Current = repo.FullName

	if err := m.cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// Remove removes a repository from the config
func (m *Manager) Remove(repoName string) error {
	// Check if trying to remove current git repository
	if currentGitRepo, err := getCurrentGitRepo(); err == nil && repoName == currentGitRepo {
		return fmt.Errorf("cannot remove current git repository")
	}

	// Find and remove repository
	found := false
	newRepos := make([]string, 0, len(m.cfg.Repositories))
	for _, repo := range m.cfg.Repositories {
		if repo == repoName {
			found = true
			// If removing current repository, we'll need to select a new one
			if repo == m.cfg.Current {
				m.cfg.Current = ""
			}
			continue
		}
		newRepos = append(newRepos, repo)
	}

	if !found {
		return fmt.Errorf("repository %s not found", repoName)
	}

	m.cfg.Repositories = newRepos

	// If we removed the current repository, select a new one
	if m.cfg.Current == "" {
		// Try to select current git repository first
		if currentGitRepo, err := getCurrentGitRepo(); err == nil {
			m.cfg.Current = currentGitRepo
		} else if len(newRepos) > 0 {
			// Otherwise select the first repository in the list
			m.cfg.Current = newRepos[0]
		}
	}

	if err := m.cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// getCurrentGitRepo returns the current git repository in owner/repo format
func getCurrentGitRepo() (string, error) {
	// Run git remote -v
	cmd := exec.Command("git", "remote", "-v")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git remote: %w", err)
	}

	// Parse output
	lines := strings.Split(string(output), "\n")
	if len(lines) == 0 {
		return "", fmt.Errorf("no git remote found")
	}

	// Get first remote URL
	parts := strings.Fields(lines[0])
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid git remote format")
	}

	// Extract owner/repo from remote URL
	url := parts[1]
	if strings.HasPrefix(url, "git@github.com:") {
		// SSH format: git@github.com:owner/repo.git
		repoPath := strings.TrimPrefix(url, "git@github.com:")
		repoPath = strings.TrimSuffix(repoPath, ".git")
		return repoPath, nil
	} else if strings.HasPrefix(url, "https://github.com/") {
		// HTTPS format: https://github.com/owner/repo.git
		repoPath := strings.TrimPrefix(url, "https://github.com/")
		repoPath = strings.TrimSuffix(repoPath, ".git")
		return repoPath, nil
	}

	return "", fmt.Errorf("unsupported git remote URL format")
}

// List returns all repositories in the config and the current git repository
func (m *Manager) List() []string {
	repos := make([]string, 0, len(m.cfg.Repositories)+1)

	// Get current git repository
	if currentRepo, err := getCurrentGitRepo(); err == nil {
		// Add current repo if it's not already in the list
		found := false
		for _, repo := range m.cfg.Repositories {
			if repo == currentRepo {
				found = true
				break
			}
		}
		if !found {
			repos = append(repos, currentRepo)
		}
	}

	// Add repositories from config
	repos = append(repos, m.cfg.Repositories...)

	return repos
}

// Current returns the current repository or the current git repository if none is set
func (m *Manager) Current() string {
	if m.cfg.Current != "" {
		return m.cfg.Current
	}

	// Try to get current git repository
	if currentRepo, err := getCurrentGitRepo(); err == nil {
		return currentRepo
	}

	return ""
}

// Switch sets the current repository
func (m *Manager) Switch(repoName string) error {
	found := false

	// Verify repository is current git repository
	if currentGitRepo, err := getCurrentGitRepo(); err == nil {
		if repoName == currentGitRepo {
			found = true
		}
	}

	// Verify repository is in config
	for _, r := range m.cfg.Repositories {
		if r == repoName {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("repository %s not found in config or current git repository", repoName)
	}

	// Update current
	m.cfg.Current = repoName
	if err := m.cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
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

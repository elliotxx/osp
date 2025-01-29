package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/cli/oauth"
	"github.com/elliotxx/osp/internal/config"
)

const (
	// OAuth endpoints
	oauthHost     = "github.com"
	oauthTokenURL = "https://github.com/login/oauth/access_token"
)

// Manager handles GitHub authentication
type Manager struct {
	cfg *config.Config
}

// NewManager creates a new auth manager
func NewManager(cfg *config.Config) *Manager {
	return &Manager{
		cfg: cfg,
	}
}

// Login performs GitHub OAuth login
func (m *Manager) Login(ctx context.Context) error {
	flow := &oauth.Flow{
		Host:     oauthHost,
		ClientID: os.Getenv("GITHUB_CLIENT_ID"),
		Scopes:   []string{"repo", "read:org"},
	}

	token, err := flow.DetectFlow(ctx)
	if err != nil {
		return fmt.Errorf("failed to perform OAuth flow: %w", err)
	}

	// Verify token
	if err := m.verifyToken(token); err != nil {
		return fmt.Errorf("failed to verify token: %w", err)
	}

	// Save token
	m.cfg.Auth.Token = token
	if err := m.cfg.Save(""); err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}

	return nil
}

// Logout removes stored credentials
func (m *Manager) Logout() error {
	m.cfg.Auth.Token = ""
	if err := m.cfg.Save(""); err != nil {
		return fmt.Errorf("failed to remove token: %w", err)
	}
	return nil
}

// verifyToken verifies the GitHub token
func (m *Manager) verifyToken(token string) error {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "token "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var response struct {
		Login string `json:"login"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}

	return nil
}

// GetToken returns the stored GitHub token
func (m *Manager) GetToken() string {
	return m.cfg.Auth.Token
}

// HasToken checks if a token is stored
func (m *Manager) HasToken() bool {
	return m.cfg.Auth.Token != ""
}

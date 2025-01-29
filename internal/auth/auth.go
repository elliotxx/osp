package auth

import (
	"context"
	"fmt"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/cli/go-gh/v2/pkg/auth"
	"github.com/elliotxx/osp/internal/config"
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
	// Get token from GitHub CLI
	token, source := auth.TokenForHost("github.com")
	if token == "" {
		return fmt.Errorf("no authentication token found, please run 'gh auth login' first")
	}

	// Verify token
	if err := m.verifyToken(token); err != nil {
		return fmt.Errorf("failed to verify token: %w", err)
	}

	// Save token
	m.cfg.Auth.Token = token
	if err := m.cfg.Save(); err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}

	fmt.Printf("✓ Authenticated via %s\n", source)
	return nil
}

// Logout removes stored credentials
func (m *Manager) Logout() error {
	m.cfg.Auth.Token = ""
	if err := m.cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}
	return nil
}

// verifyToken verifies the GitHub token
func (m *Manager) verifyToken(token string) error {
	client, err := api.NewRESTClient(api.ClientOptions{
		AuthToken: token,
		Host:      "github.com",
	})
	if err != nil {
		return err
	}

	var response struct {
		Login string `json:"login"`
	}
	if err := client.Get("user", &response); err != nil {
		return err
	}

	return nil
}

// GetToken returns the stored GitHub token
func (m *Manager) GetToken() string {
	token, _ := auth.TokenForHost("github.com")
	return token
}

// HasToken checks if a token is stored
func (m *Manager) HasToken() bool {
	token := m.GetToken()
	return token != ""
}

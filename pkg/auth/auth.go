package auth

import (
	"fmt"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/cli/go-gh/v2/pkg/auth"
)

// Login performs GitHub OAuth login
func Login() (string, error) {
	// Get token from GitHub CLI
	token, source := auth.TokenForHost("github.com")
	if token == "" {
		return "", fmt.Errorf("no authentication token found, please run 'gh auth login' first")
	}

	// Verify token
	if err := verifyToken(token); err != nil {
		return "", fmt.Errorf("failed to verify token: %w", err)
	}

	fmt.Printf(" Authenticated via %s\n", source)
	return token, nil
}

// Logout removes stored credentials
func Logout() error {
	// We don't need to do anything here since we're using gh auth
	return nil
}

// verifyToken verifies the GitHub token
func verifyToken(token string) error {
	opts := api.ClientOptions{
		AuthToken: token,
	}

	client, err := api.NewRESTClient(opts)
	if err != nil {
		return err
	}

	var response struct {
		Login string `json:"login"`
	}
	err = client.Get("user", &response)
	if err != nil {
		return err
	}

	if response.Login == "" {
		return fmt.Errorf("invalid token")
	}

	return nil
}

// GetToken returns the GitHub token
func GetToken() (string, error) {
	token, _ := auth.TokenForHost("github.com")
	if token == "" {
		return "", fmt.Errorf("no authentication token found, please run 'gh auth login' first")
	}
	return token, nil
}

// HasToken checks if a token is stored
func HasToken() bool {
	token, _ := auth.TokenForHost("github.com")
	return token != ""
}

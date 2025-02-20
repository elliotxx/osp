package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/cli/oauth/device"
	"github.com/elliotxx/osp/pkg/config"
	"github.com/elliotxx/osp/pkg/log"
	"github.com/zalando/go-keyring"
)

const (
	// GitHub OAuth application credentials (same as GitHub CLI)
	clientID     = "178c6fc778ccc68e1d6a"
	clientSecret = "34ddeff2b558a23d38fba8a6de74f086ede1cc0b"

	// Token storage
	serviceName = "osp:github.com"

	// GitHub API endpoints
	githubAPI = "https://api.github.com"
)

// ErrNotAuthenticated is returned when user is not authenticated
var ErrNotAuthenticated = errors.New("not authenticated")

// Login performs GitHub OAuth device flow login
func Login() (string, error) {
	// 1. Start OAuth device flow
	code, err := device.RequestCode(
		http.DefaultClient,
		"https://github.com/login/device/code",
		clientID,
		[]string{"repo", "read:org"},
	)
	if err != nil {
		return "", fmt.Errorf("failed to initialize OAuth flow: %w", err)
	}

	// 2. Show device code to user
	log.Info("First copy your one-time code: %s", log.Bold(code.UserCode))
	log.N().Info("%s to open github.com in your browser... ", log.Bold("Press Enter"))
	fmt.Scanln() // Wait for Enter

	if err := openBrowser(code.VerificationURI); err != nil {
		log.Error("Failed to open browser: %v", err)
		log.Info("Please visit %s to authenticate", log.Bold(code.VerificationURI))
	}

	// 3. Wait for user to complete authentication
	accessToken, err := device.Wait(
		context.Background(),
		http.DefaultClient,
		"https://github.com/login/oauth/access_token",
		device.WaitOptions{
			ClientID:   clientID,
			DeviceCode: code,
		},
	)
	if err != nil {
		return "", fmt.Errorf("authentication failed: %w", err)
	}

	// 4. Get user info and store token
	username, err := getUserInfo(accessToken.Token)
	if err != nil {
		return "", fmt.Errorf("failed to get user info: %w", err)
	}

	// 5. Store token securely
	if err := SaveToken(username, accessToken.Token); err != nil {
		return "", fmt.Errorf("failed to store token: %w", err)
	}

	log.Success("Authentication complete.")
	log.Success("Logged in as %s", log.Bold(username))
	return accessToken.Token, nil
}

// Logout removes stored credentials
func Logout() error {
	if err := RemoveToken(); err != nil {
		return fmt.Errorf("failed to remove token: %w", err)
	}
	return nil
}

// GetToken returns the stored GitHub token
func GetToken() (string, error) {
	// Try to get token from environment variables first
	log.Debug("Checking environment variables for token...")
	if token := os.Getenv("GH_TOKEN"); token != "" {
		log.Debug("Found token in GH_TOKEN")
		return token, nil
	}
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		log.Debug("Found token in GITHUB_TOKEN")
		return token, nil
	}

	log.Debug("No token found in environment variables, checking stored credentials...")
	token, err := getStoredToken()
	if err != nil {
		return "", fmt.Errorf("failed to get stored token: %w", err)
	}
	log.Debug("Successfully retrieved token from keyring")
	return token, nil
}

// GetStatus returns the current authentication status
func GetStatus() ([]*Status, error) {
	// Check authentication
	if err := CheckAuth(); err != nil {
		return nil, err
	}

	log.Debug("Checking authentication status...")
	statuses := make([]*Status, 0, 3)

	// Check environment variables first
	log.Debug("Checking environment variables...")
	envTokens := map[string]string{
		"GITHUB_TOKEN": os.Getenv("GITHUB_TOKEN"),
		"GH_TOKEN":     os.Getenv("GH_TOKEN"),
	}

	for envName, token := range envTokens {
		if token == "" {
			log.Debug("No token found in %s", envName)
			continue
		}
		log.Debug("Found token in %s, validating...", envName)

		// Validate token
		if err := validateToken(token); err != nil {
			log.Warn("Failed to validate token from %s: %v", envName, err)
			continue // Skip invalid token
		}
		log.Debug("Token validated successfully")

		// Get token scopes
		log.Debug("Getting token scopes...")
		scopes, err := getTokenScopes(token)
		if err != nil {
			log.Warn("Failed to get token scopes: %v", err)
			scopes = []string{"unknown"}
		} else {
			log.Debug("Token scopes: %v", scopes)
		}

		// Get username (optional, don't fail if this fails)
		username := "unknown"
		if u, err := getUserInfo(token); err == nil {
			username = u
		}

		statuses = append(statuses, &Status{
			Username:     username,
			Token:        token,
			TokenDisplay: token[:3] + strings.Repeat("*", 37),
			StorageType:  envName,
			IsKeyring:    false,
			Scopes:       scopes,
			Active:       true,
		})
	}

	// Then check stored token
	log.Debug("Checking stored credentials...")
	token, err := getStoredToken()
	if err != nil {
		log.Warn("Failed to get stored token: %v", err)
		return statuses, nil
	}
	log.Debug("Found stored token")

	// Validate token
	if err := validateToken(token); err != nil {
		log.Warn("Failed to validate token from keyring: %v", err)
	} else {
		log.Debug("Token validated successfully")

		// Get token scopes
		log.Debug("Getting token scopes...")
		scopes, err := getTokenScopes(token)
		if err != nil {
			log.Warn("Failed to get token scopes: %v", err)
			scopes = []string{"unknown"}
		} else {
			log.Debug("Token scopes: %v", scopes)
		}

		username, err := config.GetUsername()
		if err != nil {
			return nil, fmt.Errorf("failed to get stored username: %w", err)
		}

		statuses = append(statuses, &Status{
			Username:     username,
			Token:        token,
			TokenDisplay: token[:3] + strings.Repeat("*", 37),
			StorageType:  "keyring",
			IsKeyring:    true,
			Scopes:       scopes,
			Active:       len(statuses) == 0, // Active only if no env token
		})
	}

	log.Debug("Found %d authentication methods", len(statuses))
	return statuses, nil
}

// Status represents the current authentication status
type Status struct {
	Username     string
	Token        string
	TokenDisplay string
	StorageType  string
	IsKeyring    bool
	Scopes       []string
	Active       bool
}

// validateToken validates the token using the rate_limit API
// This is a minimal permission API that should work for any valid token
func validateToken(token string) error {
	req, err := http.NewRequest(http.MethodGet, githubAPI+"/rate_limit", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("GitHub API error: %s", string(body))
	}

	return nil
}

// getUserInfo gets the GitHub user information using the token
func getUserInfo(token string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, githubAPI+"/user", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("GitHub API error: %s", string(body))
	}

	var response struct {
		Login string `json:"login"`
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	if response.Login == "" {
		return "", fmt.Errorf("invalid token")
	}

	return response.Login, nil
}

// getTokenScopes gets the scopes of the token
func getTokenScopes(token string) ([]string, error) {
	req, err := http.NewRequest(http.MethodGet, githubAPI+"/user", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return []string{"unknown"}, nil
		}
		return nil, fmt.Errorf("failed to get token scopes")
	}

	// Get scopes from response header
	scopesHeader := resp.Header.Get("X-OAuth-Scopes")
	if scopesHeader == "" {
		return []string{}, nil
	}

	// Parse scopes
	scopes := strings.Split(scopesHeader, ", ")
	return scopes, nil
}

// getStoredToken gets the stored token from the keyring
func getStoredToken() (string, error) {
	username, err := config.GetUsername()
	if err != nil {
		return "", fmt.Errorf("failed to get username: %w", err)
	}

	token, err := keyring.Get(serviceName, username)
	if err != nil {
		return "", fmt.Errorf("failed to get token from system keyring: %w", err)
	}

	return token, nil
}

// SaveToken saves the token to keyring
func SaveToken(username, token string) error {
	// Save username to state
	if err := config.SaveUsername(username); err != nil {
		return fmt.Errorf("failed to save username: %w", err)
	}

	// Save token to keyring
	if err := keyring.Set(serviceName, username, token); err != nil {
		return fmt.Errorf("failed to save token to keyring: %w", err)
	}

	return nil
}

// RemoveToken removes the token from keyring
func RemoveToken() error {
	username, err := config.GetUsername()
	if err != nil {
		//nolint:nilerr
		return nil // If no username found, nothing to remove
	}

	// Remove token from keyring
	if err := keyring.Delete(serviceName, username); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove token from keyring: %w", err)
	}

	// Remove username from state
	if err := config.RemoveUsername(); err != nil {
		return fmt.Errorf("failed to remove username: %w", err)
	}

	return nil
}

// openBrowser opens the specified URL in the default browser
func openBrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("cmd", "/c", "start", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return err
}

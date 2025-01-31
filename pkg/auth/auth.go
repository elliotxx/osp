package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"os/exec"

	"github.com/cli/oauth/device"
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
		log.Error("Failed to open browser: %v\n", err)
		log.Info("Please visit %s to authenticate\n", log.Bold(code.VerificationURI))
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
	if err := storeToken(username, accessToken.Token); err != nil {
		return "", fmt.Errorf("failed to store token: %w", err)
	}

	log.Success("Authentication complete.")
	log.Success("Logged in as %s", log.Bold(username))
	return accessToken.Token, nil
}

// Logout removes stored credentials
func Logout() error {
	username, err := getStoredUsername()
	if err != nil {
		return fmt.Errorf("failed to get stored username: %w", err)
	}

	if err := keyring.Delete(serviceName, username); err != nil {
		// If keyring is unavailable, attempt to delete the config file
		configDir, err := getConfigDir()
		if err != nil {
			return err
		}
		tokenFile := filepath.Join(configDir, "token")
		if err := os.Remove(tokenFile); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove token file: %w", err)
		}
	}
	return nil
}

// GetToken returns the stored GitHub token
func GetToken() (string, error) {
	username, err := getStoredUsername()
	if err != nil {
		return "", fmt.Errorf("failed to get stored username: %w", err)
	}

	// 1. Try to get token from keyring
	token, err := keyring.Get(serviceName, username)
	if err == nil {
		return token, nil
	}

	// 2. If keyring is not available, try config file
	configDir, err := getConfigDir()
	if err != nil {
		return "", err
	}

	tokenFile := filepath.Join(configDir, "token")
	data, err := os.ReadFile(tokenFile)
	if err != nil {
		return "", fmt.Errorf("no authentication token found, please run 'osp auth login' first")
	}

	return strings.TrimSpace(string(data)), nil
}

// HasToken checks if a token is stored
func HasToken() bool {
	_, err := GetToken()
	return err == nil
}

// GetStatus returns the current authentication status
func GetStatus() (*Status, error) {
	// Get token
	token, err := GetToken()
	if err != nil {
		return nil, err
	}

	// Get username
	username, err := getStoredUsername()
	if err != nil {
		return nil, err
	}

	// Get token scopes
	scopes, err := getTokenScopes(token)
	if err != nil {
		return nil, err
	}

	// Check if using keyring
	isKeyring := true
	_, err = keyring.Get(serviceName, username)
	if err != nil {
		isKeyring = false
	}

	// Format token display
	tokenDisplay := "none"
	if token != "" {
		tokenDisplay = token[:3] + strings.Repeat("*", 37)
	}

	// Format storage type
	storageType := "file"
	if isKeyring {
		storageType = "keyring"
	}

	return &Status{
		Username:     username,
		Token:        token,
		TokenDisplay: tokenDisplay,
		StorageType:  storageType,
		IsKeyring:    isKeyring,
		Scopes:       scopes,
	}, nil
}

// Status represents the current authentication status
type Status struct {
	Username     string
	Token        string
	TokenDisplay string // Token with most characters masked
	StorageType  string // "keyring" or "file"
	IsKeyring    bool
	Scopes       []string
}

// getUserInfo gets the GitHub user information using the token
func getUserInfo(token string) (string, error) {
	req, err := http.NewRequest("GET", githubAPI+"/user", nil)
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
	req, err := http.NewRequest("GET", githubAPI+"/user", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
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

// storeToken stores the token securely
func storeToken(username, token string) error {
	// 1. Create config directory
	configDir, err := getConfigDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// 2. Store username
	usernameFile := filepath.Join(configDir, "username")
	if err := os.WriteFile(usernameFile, []byte(username), 0600); err != nil {
		return fmt.Errorf("failed to store username: %w", err)
	}

	// 3. Try to store token in keyring
	err = keyring.Set(serviceName, username, token)
	if err == nil {
		return nil
	}

	// 4. If keyring is not available, store token in config file
	tokenFile := filepath.Join(configDir, "token")
	return os.WriteFile(tokenFile, []byte(token), 0600)
}

// getStoredUsername gets the username from the config file
func getStoredUsername() (string, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return "", err
	}

	usernameFile := filepath.Join(configDir, "username")
	data, err := os.ReadFile(usernameFile)
	if err != nil {
		return "", fmt.Errorf("no stored username found, please run 'osp auth login' first")
	}

	username := strings.TrimSpace(string(data))
	if username == "" {
		return "", fmt.Errorf("invalid stored username")
	}

	return username, nil
}

// getConfigDir returns the configuration directory
func getConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, ".config", "osp"), nil
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

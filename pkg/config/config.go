package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/elliotxx/osp/pkg/log"
	"gopkg.in/yaml.v3"
)

const (
	// AppName is the application name used in paths
	AppName = "osp"

	// ConfigFileName is the name of the config file
	ConfigFileName = "config.yaml"

	// StateFileName is the name of the state file
	StateFileName = "state.yaml"

	// DefaultDirMode is the default mode for directories
	DefaultDirMode = 0o700

	// DefaultFileMode is the default mode for files
	DefaultFileMode = 0o600
)

// Config represents the application configuration
type Config struct{}

// State represents the application state
type State struct {
	// Username for authentication
	Username string `yaml:"username,omitempty"`

	// Current repository
	Current string `yaml:"current,omitempty"`

	// List of repositories
	Repositories []string `yaml:"repositories,omitempty"`
}

// GetConfigHome returns XDG_CONFIG_HOME
func GetConfigHome() string {
	return xdg.ConfigHome
}

// GetStateHome returns XDG_STATE_HOME
func GetStateHome() string {
	return xdg.StateHome
}

// GetDataHome returns XDG_DATA_HOME
func GetDataHome() string {
	return xdg.DataHome
}

// GetCacheHome returns XDG_CACHE_HOME
func GetCacheHome() string {
	return xdg.CacheHome
}

// GetConfigDir returns the configuration directory path
func GetConfigDir() string {
	configDir := filepath.Join(xdg.ConfigHome, AppName)
	log.Debug("Config directory: %s", configDir)

	// Check if config directory exists
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		log.Debug("Config directory does not exist: %s", configDir)
		// Create config directory with proper permissions
		if err := os.MkdirAll(configDir, DefaultDirMode); err != nil {
			log.Debug("Failed to create config directory: %v", err)
			return "."
		}
		log.Debug("Created config directory: %s", configDir)
	}

	return configDir
}

// GetStateDir returns OSP state directory for storing program state
func GetStateDir() string {
	stateDir := filepath.Join(xdg.StateHome, AppName)
	log.Debug("State directory: %s", stateDir)

	// Create state directory if it doesn't exist
	if _, err := os.Stat(stateDir); os.IsNotExist(err) {
		if err := os.MkdirAll(stateDir, DefaultDirMode); err != nil {
			log.Debug("Failed to create state directory: %v", err)
			return "."
		}
	}

	return stateDir
}

// GetStateFile returns the path to the state file
func GetStateFile() string {
	return filepath.Join(GetStateDir(), StateFileName)
}

// GetConfigFile returns the path to the config file
func GetConfigFile() string {
	// Get the config file path according to XDG specification
	configPath := filepath.Join(GetConfigDir(), ConfigFileName)
	log.Debug("Config file path: %s", configPath)
	return configPath
}

// GetUsername gets the username from the state file
func GetUsername() (string, error) {
	state, err := LoadState()
	if err != nil {
		return "", err
	}
	if state.Username == "" {
		return "", fmt.Errorf("username not found")
	}
	return state.Username, nil
}

// SaveUsername saves the username to state file
func SaveUsername(username string) error {
	state, err := LoadState()
	if err != nil {
		state = &State{}
	}
	state.Username = username
	return SaveState(state)
}

// RemoveUsername removes the username from state
func RemoveUsername() error {
	state, err := LoadState()
	if err != nil {
		//nolint:nilerr
		return nil // If state doesn't exist, nothing to remove
	}
	state.Username = ""
	return SaveState(state)
}

// GetCurrentRepo gets the current repository from state
func GetCurrentRepo() (string, error) {
	state, err := LoadState()
	if err != nil {
		return "", err
	}
	return state.Current, nil
}

// SaveCurrentRepo saves the current repository to state
func SaveCurrentRepo(current string) error {
	state, err := LoadState()
	if err != nil {
		state = &State{}
	}
	state.Current = current
	return SaveState(state)
}

// GetRepositories gets the list of repositories from state
func GetRepositories() ([]string, error) {
	state, err := LoadState()
	if err != nil {
		return nil, err
	}
	return state.Repositories, nil
}

// SaveRepositories saves the list of repositories to state
func SaveRepositories(repos []string) error {
	state, err := LoadState()
	if err != nil {
		state = &State{}
	}
	state.Repositories = repos
	return SaveState(state)
}

// LoadState loads the application state
func LoadState() (*State, error) {
	statePath := GetStateFile()
	log.Debug("Loading state from: %s", statePath)

	// Return empty state if file doesn't exist
	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		return &State{}, nil
	}

	data, err := os.ReadFile(statePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	state := &State{}
	if err := yaml.Unmarshal(data, state); err != nil {
		return nil, fmt.Errorf("failed to parse state file: %w", err)
	}

	return state, nil
}

// SaveState saves the application state
func SaveState(state *State) error {
	statePath := GetStateFile()
	log.Debug("Saving state to: %s", statePath)

	// Create backup if file exists
	if _, err := os.Stat(statePath); err == nil {
		backupPath := statePath + ".bak"
		log.Debug("Creating backup: %s", backupPath)
		if err := os.Rename(statePath, backupPath); err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
	}

	// Marshal state to YAML
	data, err := yaml.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	// Write state file
	if err := os.WriteFile(statePath, data, DefaultFileMode); err != nil {
		// Try to restore backup if write failed
		if _, err := os.Stat(statePath + ".bak"); err == nil {
			log.Debug("Write failed, attempting to restore backup")
			if restoreErr := os.Rename(statePath+".bak", statePath); restoreErr != nil {
				log.Debug("Failed to restore backup: %v", restoreErr)
			}
		}
		return fmt.Errorf("failed to write state file: %w", err)
	}

	// Remove backup after successful write
	if _, err := os.Stat(statePath + ".bak"); err == nil {
		log.Debug("Removing backup file")
		if err := os.Remove(statePath + ".bak"); err != nil {
			log.Debug("Failed to remove backup: %v", err)
		}
	}

	return nil
}

// Load loads the configuration from file
func Load(path string) (*Config, error) {
	if path == "" {
		path = GetConfigFile()
	}
	log.Debug("Loading config from: %s", path)

	// Create default config if file doesn't exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		cfg := &Config{}
		if err := cfg.Save(); err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return cfg, nil
}

// Save saves the configuration to file
func (c *Config) Save() error {
	path := GetConfigFile()
	log.Debug("Saving config to: %s", path)

	// Backup existing config if it exists
	if _, err := os.Stat(path); err == nil {
		backupPath := path + ".bak"
		log.Debug("Creating backup: %s", backupPath)
		if err := os.Rename(path, backupPath); err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
	}

	// Marshal config to YAML
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write config file
	if err := os.WriteFile(path, data, DefaultFileMode); err != nil {
		// Try to restore backup if write failed
		if _, err := os.Stat(path + ".bak"); err == nil {
			log.Debug("Write failed, attempting to restore backup")
			if restoreErr := os.Rename(path+".bak", path); restoreErr != nil {
				log.Debug("Failed to restore backup: %v", restoreErr)
			}
		}
		return fmt.Errorf("failed to write config file: %w", err)
	}

	// Remove backup after successful write
	if _, err := os.Stat(path + ".bak"); err == nil {
		log.Debug("Removing backup file")
		if err := os.Remove(path + ".bak"); err != nil {
			log.Debug("Failed to remove backup: %v", err)
		}
	}

	return nil
}

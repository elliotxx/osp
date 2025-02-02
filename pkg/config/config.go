package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"
	"github.com/elliotxx/osp/pkg/log"
	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	// Authentication settings
	Auth struct {
		Token string `yaml:"token"`
	} `yaml:"auth"`

	// Repository settings
	Current      string   `yaml:"current"`
	Repositories []string `yaml:"repositories"`
}

// GetConfigDir returns the configuration directory path
func GetConfigDir() string {
	configDir := filepath.Join(xdg.ConfigHome, "osp")
	log.Debug("Config directory: %s", configDir)

	// Check if config directory exists
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		log.Debug("Config directory does not exist: %s", configDir)
		// Create config directory with proper permissions
		if err := os.MkdirAll(configDir, 0700); err != nil {
			log.Debug("Failed to create config directory: %v", err)
			return "."
		}
		log.Debug("Created config directory: %s", configDir)
	}

	return configDir
}

// getConfigPath returns the path to the config file
func getConfigPath() (string, error) {
	// Get the config file path according to XDG specification
	configPath := filepath.Join(GetConfigDir(), "config.yml")
	log.Debug("Config file path: %s", configPath)
	return configPath, nil
}

// Load loads the configuration from file
func Load(path string) (*Config, error) {
	if path == "" {
		var err error
		path, err = getConfigPath()
		if err != nil {
			return nil, fmt.Errorf("failed to get config path: %w", err)
		}
	}
	log.Debug("Loading config from: %s", path)

	// Initialize empty config
	cfg := &Config{}

	// Load config file if it exists
	if _, err := os.Stat(path); err == nil {
		log.Debug("Config file exists, reading...")
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		if err := yaml.Unmarshal(data, cfg); err != nil {
			// Backup the corrupted config file
			backupPath := path + fmt.Sprintf(".bak.%d", time.Now().Unix())
			log.Debug("Config file is corrupted, creating backup: %s", backupPath)
			if err := os.Rename(path, backupPath); err != nil {
				return nil, fmt.Errorf("failed to backup corrupted config: %w", err)
			}
			return nil, fmt.Errorf("failed to parse config (backup created at %s): %w", backupPath, err)
		}
		log.Debug("Config loaded successfully")
	} else {
		log.Debug("Config file does not exist, using empty config")
	}

	return cfg, nil
}

// Save saves the configuration to file
func (c *Config) Save() error {
	path, err := getConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get config path: %w", err)
	}
	log.Debug("Saving config to: %s", path)

	// Backup existing config if it exists
	backupPath := path + ".bak"
	if _, err := os.Stat(path); err == nil {
		log.Debug("Creating backup of existing config: %s", backupPath)
		if err := os.Rename(path, backupPath); err != nil {
			return fmt.Errorf("failed to backup config: %w", err)
		}
	}

	// Marshal config to YAML
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write config file
	if err := os.WriteFile(path, data, 0600); err != nil {
		// Try to restore backup if write failed
		if _, err := os.Stat(backupPath); err == nil {
			log.Debug("Write failed, attempting to restore backup")
			if restoreErr := os.Rename(backupPath, path); restoreErr != nil {
				return fmt.Errorf("failed to write config and restore backup: %v (restore error: %v)", err, restoreErr)
			}
		}
		return fmt.Errorf("failed to write config: %w", err)
	}

	// Remove backup file if write succeeded
	if _, err := os.Stat(backupPath); err == nil {
		log.Debug("Write succeeded, removing backup")
		if err := os.Remove(backupPath); err != nil {
			log.Debug("Failed to remove backup file: %v", err)
		}
	}

	log.Debug("Config saved successfully")
	return nil
}

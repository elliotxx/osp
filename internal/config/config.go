package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Auth struct {
		Token string `yaml:"token"`
	} `yaml:"auth"`
	Current      string   `yaml:"current"`
	Repositories []string `yaml:"repositories"`
}

// Load loads the configuration from file
func Load(path string) (*Config, error) {
	if path == "" {
		configDir, err := os.UserConfigDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get config directory: %w", err)
		}
		path = filepath.Join(configDir, "osp", "config.yml")
	}

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(path)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	// Load config file
	cfg := &Config{}
	if _, err := os.Stat(path); err == nil {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	}

	return cfg, nil
}

// Save saves the configuration to file
func (c *Config) Save(path string) error {
	if path == "" {
		configDir, err := os.UserConfigDir()
		if err != nil {
			return fmt.Errorf("failed to get config directory: %w", err)
		}
		path = filepath.Join(configDir, "osp", "config.yml")
	}

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(path)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal config to YAML
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write config file
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

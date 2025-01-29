package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the main configuration structure
type Config struct {
	Version string       `yaml:"version"`
	Debug   bool         `yaml:"debug"`
	Quiet   bool         `yaml:"quiet"`
	Auth    AuthConfig   `yaml:"auth"`
	Repos   []RepoConfig `yaml:"repos"`
	Current string       `yaml:"current_repo"`
	Default Defaults     `yaml:"defaults"`
	Custom  Custom       `yaml:"custom"`
}

// AuthConfig holds authentication related configuration
type AuthConfig struct {
	Token   string `yaml:"token"`
	Host    string `yaml:"host"`
	Keyring bool   `yaml:"keyring"`
}

// RepoConfig holds repository specific configuration
type RepoConfig struct {
	Name   string                 `yaml:"name"`
	Alias  string                 `yaml:"alias"`
	Path   string                 `yaml:"path"`
	Config map[string]interface{} `yaml:"config"`
}

// Defaults holds default configuration values
type Defaults struct {
	Period    string `yaml:"period"`
	Format    string `yaml:"format"`
	AutoSync  bool   `yaml:"auto_sync"`
	CacheTTL  string `yaml:"cache_ttl"`
}

// Custom holds custom configuration values
type Custom struct {
	TemplatesDir string `yaml:"templates_dir"`
	ReportsDir   string `yaml:"reports_dir"`
}

// Load loads configuration from file
func Load(path string) (*Config, error) {
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		path = filepath.Join(home, ".config", "osp", "config.yml")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// Save saves configuration to file
func (c *Config) Save(path string) error {
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		path = filepath.Join(home, ".config", "osp", "config.yml")
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		Version: "1.0",
		Debug:   false,
		Quiet:   false,
		Auth: AuthConfig{
			Host:    "github.com",
			Keyring: true,
		},
		Default: Defaults{
			Period:    "7d",
			Format:    "markdown",
			AutoSync:  true,
			CacheTTL:  "24h",
		},
		Custom: Custom{
			TemplatesDir: "${HOME}/.config/osp/templates",
			ReportsDir:   "${HOME}/.config/osp/reports",
		},
	}
}

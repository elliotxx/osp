package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/elliotxx/osp/pkg/config"
	"github.com/elliotxx/osp/pkg/log"
	"github.com/elliotxx/osp/pkg/util/prompt"
	"github.com/spf13/cobra"
)

var (
	configCmd = &cobra.Command{
		Use:   "config",
		Short: "Manage configuration files and data",
		Long: `Manage configuration files and application data.

This command helps you manage OSP's configuration files and data directories.
It provides subcommands to list configuration locations, edit configuration files,
and clean up configuration files.`,
		RunE: runConfigList,
	}

	configListCmd = &cobra.Command{
		Use:     "list",
		Short:   "Show config locations",
		Aliases: []string{"ls", "locations"},
		Long:    "Show configuration file locations and application data directories.",
		RunE:    runConfigList,
	}

	configEditCmd = &cobra.Command{
		Use:   "edit",
		Short: "Edit configuration file with default editor",
		RunE:  runConfigEdit,
	}

	// Options for clean command
	cleanFlags struct {
		force bool
	}

	configCleanCmd = &cobra.Command{
		Use:   "clean",
		Short: "Clean up all configuration files and data",
		Long: `Clean up all configuration files and data.

This will remove:
- Configuration files (~/.config/osp/*)
- State data (~/.local/state/osp/*)

Note: This action cannot be undone!`,
		RunE: runConfigClean,
	}
)

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configEditCmd)
	configCmd.AddCommand(configCleanCmd)

	// Add flags for clean command
	configCleanCmd.Flags().BoolVarP(&cleanFlags.force, "force", "f", false, "Skip confirmation prompt")
}

func runConfigList(cmd *cobra.Command, args []string) error {
	// Get XDG paths from config package
	configHome := config.GetConfigHome()
	stateHome := config.GetStateHome()

	// Get OSP paths
	configDir := config.GetConfigDir()
	configFile := config.GetConfigFile()
	stateDir := config.GetStateDir()

	// Print XDG environment variables
	log.B().Log("XDG Base Directories:")
	log.L(1).Info("%-16s = %s", "XDG_CONFIG_HOME", configHome)
	log.L(1).Info("%-16s = %s", "XDG_STATE_HOME", stateHome)

	// Print OSP locations
	log.B().Log("\nOSP Locations:")
	log.L(1).Info("Config Directory:")
	log.L(2).Info("%-12s %s", "Path:", configDir)
	if fileExists(configDir) {
		log.L(2).Success("%-12s %v", "Exists:", true)
	} else {
		log.L(2).Error("%-12s %v", "Exists:", false)
	}

	log.L(1).Info("Config File:")
	log.L(2).Info("%-12s %s", "Path:", configFile)
	if fileExists(configFile) {
		log.L(2).Success("%-12s %v", "Exists:", true)
	} else {
		log.L(2).Error("%-12s %v", "Exists:", false)
	}

	log.L(1).Info("State Directory:")
	log.L(2).Info("%-12s %s", "Path:", stateDir)
	if fileExists(stateDir) {
		log.L(2).Success("%-12s %v", "Exists:", true)
	} else {
		log.L(2).Error("%-12s %v", "Exists:", false)
	}

	return nil
}

func runConfigEdit(cmd *cobra.Command, args []string) error {
	// Get config file path
	configFile := config.GetConfigFile()

	// Get editor from environment or use default
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim" // Default to vim
	}

	// Create config file if it doesn't exist
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		log.Debug("Config file does not exist, creating empty file")
		if err := os.MkdirAll(filepath.Dir(configFile), 0o700); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}
		if err := os.WriteFile(configFile, []byte(""), 0o600); err != nil {
			return fmt.Errorf("failed to create config file: %w", err)
		}
		log.Debug("Created empty config file: %s", configFile)
	}

	log.Debug("Opening config file with editor: %s %s", editor, configFile)

	// Open editor
	cmd2 := exec.Command(editor, configFile)
	cmd2.Stdin = os.Stdin
	cmd2.Stdout = os.Stdout
	cmd2.Stderr = os.Stderr

	return cmd2.Run()
}

func runConfigClean(cmd *cobra.Command, args []string) error {
	// Get config directory
	configDir := config.GetConfigDir()
	configFile := config.GetConfigFile()

	// Get data directory
	stateDir := config.GetStateDir()

	// Print locations
	log.Info("The following files and directories will be removed:")
	log.L(1).Info("Config directory: %s", configDir)
	log.L(2).Info("Config file: %s", configFile)
	log.L(1).Info("State directory: %s", stateDir)

	// Check which directories exist
	var existingDirs []string
	if _, err := os.Stat(configDir); err == nil {
		existingDirs = append(existingDirs, configDir)
	}
	if _, err := os.Stat(stateDir); err == nil {
		existingDirs = append(existingDirs, stateDir)
	}

	// Skip if nothing to clean
	if len(existingDirs) == 0 {
		log.Info("Nothing to clean")
		return nil
	}

	// Show plan and get confirmation
	if !cleanFlags.force {
		confirmed, err := prompt.AskForConfirmation("Do you want to continue?")
		if err != nil {
			return err
		}
		if !confirmed {
			log.Info("Operation cancelled")
			return nil
		}
	}

	// Execute plan
	for _, dir := range existingDirs {
		log.Debug("Removing directory: %s", dir)
		if err := os.RemoveAll(dir); err != nil {
			return fmt.Errorf("failed to remove directory: %w", err)
		}
		log.Info("Removed directory: %s", dir)
	}

	log.Info("All configuration files and data have been removed")
	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

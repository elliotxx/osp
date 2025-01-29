package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "osp",
	Short: "Open Source Project Management Tool",
	Long: `OSP is a command-line tool for managing open source projects.
It helps you manage issues, milestones, planning, and more.`,
}

// Repository represents a GitHub repository
type Repository struct {
	Owner string
	Name  string
}

// Execute executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(
		newAuthCmd(),
		newPlanCmd(),
	)
}

// getCurrentRepository returns the current repository from git remote
func getCurrentRepository() (*Repository, error) {
	// Get git remote URL
	output, err := runGitCommand("remote", "get-url", "origin")
	if err != nil {
		return nil, fmt.Errorf("failed to get git remote URL: %w", err)
	}

	// Parse URL to get owner and repo
	url := strings.TrimSpace(output)
	url = strings.TrimSuffix(url, ".git")

	// Handle different URL formats:
	// - https://github.com/owner/repo.git
	// - git@github.com:owner/repo.git
	var ownerRepo string
	if strings.HasPrefix(url, "https://") {
		parts := strings.Split(url, "/")
		ownerRepo = strings.Join(parts[len(parts)-2:], "/")
	} else {
		parts := strings.Split(url, ":")
		ownerRepo = parts[len(parts)-1]
	}

	parts := strings.Split(ownerRepo, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid repository URL format: %s", url)
	}

	return &Repository{
		Owner: parts[0],
		Name:  parts[1],
	}, nil
}

// runGitCommand runs a git command and returns its output
func runGitCommand(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

package cmd

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"github.com/elliotxx/osp/pkg/config"
	"github.com/elliotxx/osp/pkg/log"
	"github.com/elliotxx/osp/pkg/repo"
)

// selectRepository prompts user to select a repository
func selectRepository(repos []string, current string) (string, error) {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "→ {{ . | cyan }}{{ if eq . \"" + current + "\" }} (current){{ end }}",
		Inactive: "  {{ . }}{{ if eq . \"" + current + "\" }} (current){{ end }}",
		Selected: "✓ {{ . | green }}",
	}

	prompt := promptui.Select{
		Label:     "Select a repository",
		Items:     repos,
		Templates: templates,
		Size:      10,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return repos[i], nil
}

var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Manage repositories",
	Long: `Manage GitHub repositories.

This command allows you to manage the GitHub repositories you want to work with.
You can add, remove, list, and switch between repositories.

Examples:
  # List all repositories
  osp repo list

  # Add a new repository
  osp repo add owner/repo

  # Remove a repository
  osp repo remove owner/repo

  # Switch to a repository
  osp repo switch owner/repo

  # Show current repository
  osp repo current`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load config
		cfg, err := config.Load("")
		if err != nil {
			return err
		}

		// Create repository manager
		repoManager, err := repo.NewManager(cfg)
		if err != nil {
			return err
		}

		repos := repoManager.List()
		current := repoManager.Current()

		if len(repos) == 0 {
			log.Info("No repositories found.")
			return nil
		}

		// Select repository
		selected, err := selectRepository(repos, current)
		if err != nil {
			return fmt.Errorf("failed to select repository: %w", err)
		}

		// Update config if changed
		if selected != current {
			if err := repoManager.Switch(selected); err != nil {
				return err
			}
		}
		return nil
	},
}

var repoAddCmd = &cobra.Command{
	Use:   "add [owner/repo]",
	Short: "Add a repository to manage",
	Long:  `Add a GitHub repository to manage.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load("")
		if err != nil {
			return err
		}

		// Create repository manager
		repoManager, err := repo.NewManager(cfg)
		if err != nil {
			return err
		}

		if err := repoManager.Add(cmd.Context(), args[0]); err != nil {
			return err
		}

		log.Success("Successfully added repository %s", args[0])
		return nil
	},
}

var repoRemoveCmd = &cobra.Command{
	Use:   "remove [owner/repo]",
	Short: "Remove a repository from management",
	Long:  `Remove a GitHub repository from management.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load("")
		if err != nil {
			return err
		}

		// Create repository manager
		repoManager, err := repo.NewManager(cfg)
		if err != nil {
			return err
		}

		if err := repoManager.Remove(args[0]); err != nil {
			return err
		}

		log.Success("Successfully removed repository %s", args[0])
		return nil
	},
}

var repoListCmd = &cobra.Command{
	Use:   "list",
	Short: "List managed repositories",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load("")
		if err != nil {
			return err
		}

		repoManager, err := repo.NewManager(cfg)
		if err != nil {
			return err
		}

		repos := repoManager.List()
		current := repoManager.Current()

		if len(repos) == 0 {
			log.Info("No repositories found.")
			return nil
		}

		fmt.Println("Managed repositories:")
		for _, r := range repos {
			if r == current {
				fmt.Printf("* %s\n", r)
			} else {
				fmt.Printf("  %s\n", r)
			}
		}

		return nil
	},
}

var repoSwitchCmd = &cobra.Command{
	Use:   "switch [owner/repo]",
	Short: "Switch current repository",
	Long:  `Switch the current repository being managed.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load config
		cfg, err := config.Load("")
		if err != nil {
			return err
		}

		repoManager, err := repo.NewManager(cfg)
		if err != nil {
			return err
		}
		if err := repoManager.Switch(args[0]); err != nil {
			return err
		}

		log.Success("Successfully switched to repository %s", args[0])
		return nil
	},
}

var repoCurrentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show current repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load config
		cfg, err := config.Load("")
		if err != nil {
			return err
		}

		repoManager, err := repo.NewManager(cfg)
		if err != nil {
			return err
		}
		current := repoManager.Current()

		if current == "" {
			log.Info("No repository selected.")
			return nil
		}

		fmt.Printf("Current repository: %s\n", current)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(repoCmd)
	repoCmd.AddCommand(repoAddCmd)
	repoCmd.AddCommand(repoRemoveCmd)
	repoCmd.AddCommand(repoListCmd)
	repoCmd.AddCommand(repoSwitchCmd)
	repoCmd.AddCommand(repoCurrentCmd)
}

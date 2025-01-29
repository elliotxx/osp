package cmd

import (
	"fmt"

	"github.com/elliotxx/osp/internal/config"
	"github.com/elliotxx/osp/internal/repo"
	"github.com/spf13/cobra"
)

var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Manage repositories",
	Long:  `Manage GitHub repositories, including adding, removing, listing, and switching between them.`,
}

var repoAddCmd = &cobra.Command{
	Use:   "add [owner/repo]",
	Short: "Add a repository to manage",
	Long:  `Add a GitHub repository to manage.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load("")
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		repoManager := repo.NewManager(cfg)
		if err := repoManager.Add(cmd.Context(), args[0]); err != nil {
			return err
		}

		fmt.Printf("✓ Added repository %s\n", args[0])
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
			return fmt.Errorf("failed to load config: %w", err)
		}

		repoManager := repo.NewManager(cfg)
		if err := repoManager.Remove(args[0]); err != nil {
			return err
		}

		fmt.Printf("✓ Removed repository %s\n", args[0])
		return nil
	},
}

var repoListCmd = &cobra.Command{
	Use:   "list",
	Short: "List managed repositories",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load("")
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		repoManager := repo.NewManager(cfg)
		repos := repoManager.List()

		if len(repos) == 0 {
			fmt.Println("No repositories found.")
			return nil
		}

		fmt.Println("Managed repositories:")
		for _, r := range repos {
			if r == cfg.Current {
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
		cfg, err := config.Load("")
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		repoManager := repo.NewManager(cfg)
		if err := repoManager.Switch(args[0]); err != nil {
			return err
		}

		fmt.Printf("✓ Switched to repository %s\n", args[0])
		return nil
	},
}

var repoCurrentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show current repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load("")
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		repoManager := repo.NewManager(cfg)
		current := repoManager.Current()

		if current == "" {
			fmt.Println("No repository selected.")
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

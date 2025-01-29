package cmd

import (
	"fmt"

	"github.com/elliotxx/osp/internal/config"
	"github.com/elliotxx/osp/internal/repo"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <owner/repo>",
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
			return fmt.Errorf("failed to add repository: %w", err)
		}

		fmt.Printf("✓ Successfully added repository %s\n", args[0])
		return nil
	},
}

var removeCmd = &cobra.Command{
	Use:   "remove <owner/repo>",
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
			return fmt.Errorf("failed to remove repository: %w", err)
		}

		fmt.Printf("✓ Successfully removed repository %s\n", args[0])
		return nil
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List managed repositories",
	Long:  `List all GitHub repositories being managed.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load("")
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		repoManager := repo.NewManager(cfg)
		repos := repoManager.List()
		current := repoManager.Current()

		if len(repos) == 0 {
			fmt.Println("No repositories are being managed.")
			return nil
		}

		fmt.Println("Managed repositories:")
		for _, r := range repos {
			prefix := "  "
			if r.Name == current {
				prefix = "* "
			}
			fmt.Printf("%s%s", prefix, r.Name)
			if r.Alias != "" {
				fmt.Printf(" (%s)", r.Alias)
			}
			fmt.Println()
		}
		return nil
	},
}

var switchCmd = &cobra.Command{
	Use:   "switch <owner/repo>",
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
			return fmt.Errorf("failed to switch repository: %w", err)
		}

		fmt.Printf("✓ Switched to repository %s\n", args[0])
		return nil
	},
}

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show current repository",
	Long:  `Display the current repository being managed.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load("")
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		repoManager := repo.NewManager(cfg)
		current := repoManager.Current()

		if current == "" {
			fmt.Println("No repository is currently selected.")
			return nil
		}

		fmt.Printf("Current repository: %s\n", current)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(switchCmd)
	rootCmd.AddCommand(currentCmd)
}

package cmd

import (
	"fmt"

	"github.com/elliotxx/osp/internal/auth"
	"github.com/elliotxx/osp/internal/config"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage GitHub authentication",
	Long:  `Manage GitHub authentication, including login and logout.`,
}

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to GitHub",
	Long:  `Login to GitHub using OAuth authentication.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load("")
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		authManager := auth.NewManager(cfg)
		if err := authManager.Login(cmd.Context()); err != nil {
			return fmt.Errorf("failed to login: %w", err)
		}

		fmt.Println("✓ Successfully logged in")
		return nil
	},
}

var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from GitHub",
	Long:  `Logout from GitHub and remove stored credentials.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load("")
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		authManager := auth.NewManager(cfg)
		if err := authManager.Logout(); err != nil {
			return fmt.Errorf("failed to logout: %w", err)
		}

		fmt.Println("✓ Successfully logged out")
		return nil
	},
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show authentication status",
	Long:  `Show current GitHub authentication status.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load("")
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		authManager := auth.NewManager(cfg)
		if authManager.HasToken() {
			fmt.Println("✓ Logged in to GitHub")
		} else {
			fmt.Println("✗ Not logged in to GitHub")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(authLoginCmd)
	authCmd.AddCommand(authLogoutCmd)
	authCmd.AddCommand(authStatusCmd)
}

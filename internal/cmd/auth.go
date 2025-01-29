package cmd

import (
	"fmt"

	"github.com/elliotxx/osp/internal/auth"
	"github.com/elliotxx/osp/internal/config"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth <command>",
	Short: "Manage GitHub authentication",
	Long:  `Manage GitHub authentication and credentials.`,
}

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to GitHub",
	Long:  `Log in to GitHub using OAuth or a personal access token.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load("")
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		authManager := auth.NewManager(cfg)
		if err := authManager.Login(cmd.Context()); err != nil {
			return fmt.Errorf("login failed: %w", err)
		}

		fmt.Println("✓ Successfully logged in to GitHub")
		return nil
	},
}

var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out of GitHub",
	Long:  `Log out of GitHub and remove stored credentials.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load("")
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		authManager := auth.NewManager(cfg)
		if err := authManager.Logout(); err != nil {
			return fmt.Errorf("logout failed: %w", err)
		}

		fmt.Println("✓ Successfully logged out of GitHub")
		return nil
	},
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show authentication status",
	Long:  `Display the current authentication status and token information.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load("")
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		authManager := auth.NewManager(cfg)
		if !authManager.HasToken() {
			fmt.Println("× Not logged in to GitHub")
			return nil
		}

		fmt.Println("✓ Logged in to GitHub")
		return nil
	},
}

func init() {
	authCmd.AddCommand(authLoginCmd)
	authCmd.AddCommand(authLogoutCmd)
	authCmd.AddCommand(authStatusCmd)
}

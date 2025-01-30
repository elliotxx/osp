package cmd

import (
	"fmt"

	"github.com/elliotxx/osp/pkg/auth"
	"github.com/spf13/cobra"
)

func newAuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Authenticate with GitHub",
		Long: `Authenticate with GitHub using GitHub CLI's authentication.
This command will help you set up authentication for OSP.`,
	}

	cmd.AddCommand(newAuthLoginCmd())
	cmd.AddCommand(newAuthStatusCmd())
	cmd.AddCommand(newAuthLogoutCmd())

	return cmd
}

func newAuthLoginCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Login to GitHub",
		Long:  "Login to GitHub using GitHub CLI's authentication.",
		RunE: func(cmd *cobra.Command, args []string) error {
			token, err := auth.GetToken()
			if err != nil {
				return err
			}
			fmt.Printf("✓ Logged in with token: %s\n", token)
			return nil
		},
	}
}

func newAuthStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show authentication status",
		Long:  "Show current GitHub authentication status.",
		RunE: func(cmd *cobra.Command, args []string) error {
			token, err := auth.GetToken()
			if err != nil {
				return err
			}
			if token == "" {
				fmt.Println("× Not logged in")
				return nil
			}
			fmt.Printf("✓ Logged in with token: %s\n", token)
			return nil
		},
	}
}

func newAuthLogoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Logout from GitHub",
		Long:  "Logout from GitHub and remove stored credentials.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := auth.Logout(); err != nil {
				return err
			}
			fmt.Println("✓ Logged out")
			return nil
		},
	}
}

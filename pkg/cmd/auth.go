package cmd

import (
	"strings"

	"github.com/elliotxx/osp/pkg/auth"
	"github.com/elliotxx/osp/pkg/log"
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
			_, err := auth.Login()
			return err
		},
	}
}

func newAuthStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show authentication status",
		Long:  "Show current GitHub authentication status.",
		RunE: func(cmd *cobra.Command, args []string) error {
			status, err := auth.GetStatus()
			if err != nil {
				return err
			}

			// Format token for display
			tokenDisplay := "none"
			if status.Token != "" {
				tokenDisplay = status.Token[:3] + strings.Repeat("*", 37)
			}

			// Format storage type
			storageType := "file"
			if status.IsKeyring {
				storageType = "keyring"
			}

			// Print status
			log.B().Log("github.com")
			log.L(1).Success("Logged in to github.com account %s (%s)\n", log.Bold(status.Username), storageType)
			log.L(1).Info("Active account: %s", log.Bold("true"))
			log.L(1).Info("Token: %s", log.Bold(tokenDisplay))
			if len(status.Scopes) > 0 {
				log.L(1).Info("Token scopes: '%s'", log.Bold(strings.Join(status.Scopes, "', '")))
			}
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
			log.Success("Successfully logged out")
			return nil
		},
	}
}

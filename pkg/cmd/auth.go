package cmd

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/elliotxx/osp/pkg/auth"
	"github.com/elliotxx/osp/pkg/log"
	"github.com/spf13/cobra"
)

func newAuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Authenticate with GitHub",
		Long: heredoc.Docf(`
			Authenticate with GitHub.

			The default authentication mode is a web-based browser flow using GitHub's OAuth device flow.
			After completion, an authentication token will be stored securely in the system credential store.
			If a credential store is not found, the token will be stored in a plain text file.

			You can also authenticate by setting the %[1]sGH_TOKEN%[1]s environment variable
			to a personal access token.
		`, "`"),
		Example: heredoc.Doc(`
			# Start interactive setup
			$ osp auth login

			# Check authentication status
			$ osp auth status
		`),
	}

	cmd.AddCommand(newAuthLoginCmd())
	cmd.AddCommand(newAuthStatusCmd())
	cmd.AddCommand(newAuthLogoutCmd())

	return cmd
}

func newAuthLoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Log in to a GitHub account",
		Long: heredoc.Doc(`
			Log in to a GitHub account.

			This command will help you authenticate with GitHub using a web-based browser flow.
			A one-time code will be displayed, which you can enter at the specified URL to complete
			the authentication process.
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			log.SetNoColor(true)
			defer log.SetNoColor(false)
			_, err := auth.Login()
			return err
		},
	}

	return cmd
}

func newAuthStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "View authentication status",
		Long: heredoc.Doc(`
			Verifies and displays information about your authentication state.

			This command will test your authentication state and report whether you are properly
			authenticated. It will also display information about the authenticated user and token.
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			log.SetNoColor(true)
			defer log.SetNoColor(false)

			statuses, err := auth.GetStatus()
			if err != nil {
				return err
			}

			// Print status
			log.B().Log("github.com")
			for _, status := range statuses {
				log.L(1).Success("Logged in to github.com account %s (%s)", log.Bold(status.Username), status.StorageType)
				log.L(2).Info("Active account: %s", log.Bold(fmt.Sprintf("%v", status.Active)))
				log.L(2).Info("Token: %s", log.Bold(status.TokenDisplay))
				if len(status.Scopes) > 0 {
					log.L(2).Info("Token scopes: '%s'", log.Bold(strings.Join(status.Scopes, "', '")))
				}
			}
			return nil
		},
	}

	return cmd
}

func newAuthLogoutCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Log out from a GitHub account",
		Long: heredoc.Doc(`
			Remove authentication for a GitHub account.

			This command removes the authentication token from your system.
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := auth.Logout(); err != nil {
				return err
			}
			log.Success("Successfully logged out")
			return nil
		},
	}

	return cmd
}

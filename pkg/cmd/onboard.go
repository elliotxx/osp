package cmd

import (
	"fmt"
	"strings"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/elliotxx/osp/pkg/config"
	"github.com/elliotxx/osp/pkg/log"
	"github.com/elliotxx/osp/pkg/onboard"
	"github.com/elliotxx/osp/pkg/repo"
	"github.com/spf13/cobra"
)

var onboardCmd = &cobra.Command{
	Use:   "onboard",
	Short: "Manage onboarding content for community contributors",
	Long: `Generate and maintain onboarding content to help community contributors get started.

This command will create or update an issue that summarizes all issues suitable for community contribution,
organized by difficulty level and category. The content is designed to help new contributors find issues
that match their interests and skill levels.

Examples:
  # Update onboarding content with default settings
  osp onboard

  # Use custom labels for finding beginner-friendly issues
  osp onboard --onboard-labels="good first issue,help wanted"

  # Use custom difficulty levels
  osp onboard --difficulty-labels="difficulty/easy,difficulty/medium,difficulty/hard"

  # Use custom categories within each difficulty level
  osp onboard --category-labels="bug,feature,documentation"

  # Preview changes without updating any issues
  osp onboard --dry-run

  # Update automatically without confirmation
  osp onboard --yes

  # Specify a custom label for the target issue
  osp onboard --target-label="getting-started"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load("")
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Get repository name
		repoManager := repo.NewManager(cfg)
		repoName := repoManager.Current()
		if repoName == "" {
			return fmt.Errorf("no repository selected, use 'osp repo switch' to select a repository first")
		}
		log.Debug("Generating onboarding issues for %s", repoName)

		// Get flags
		onboardLabels, err := cmd.Flags().GetStringSlice("onboard-labels")
		if err != nil {
			return err
		}
		difficultyLabels, err := cmd.Flags().GetStringSlice("difficulty-labels")
		if err != nil {
			return err
		}
		categoryLabels, err := cmd.Flags().GetStringSlice("category-labels")
		if err != nil {
			return err
		}
		dryRun, err := cmd.Flags().GetBool("dry-run")
		if err != nil {
			return err
		}
		autoConfirm, err := cmd.Flags().GetBool("yes")
		if err != nil {
			return err
		}
		targetLabel, err := cmd.Flags().GetString("target-label")
		if err != nil {
			return err
		}
		log.Debug("Onboard labels: [%s]", strings.Join(onboardLabels, ", "))
		log.Debug("Difficulty labels: [%s]", strings.Join(difficultyLabels, ", "))
		log.Debug("Category labels: [%s]", strings.Join(categoryLabels, ", "))

		// Create GitHub client
		client, err := api.DefaultRESTClient()
		if err != nil {
			return fmt.Errorf("failed to create GitHub client: %w", err)
		}

		// Create options
		opts := onboard.Options{
			OnboardLabels:    onboardLabels,
			DifficultyLabels: difficultyLabels,
			CategoryLabels:   categoryLabels,
		}
		onboardOpts := onboard.OnboardOptions{
			TargetLabel: targetLabel,
			DryRun:      dryRun,
			AutoConfirm: autoConfirm,
		}

		// Create onboard manager
		onboardManager := onboard.NewManager(cfg, client)

		// Update onboarding issue
		err = onboardManager.Update(cmd.Context(), repoName, opts, onboardOpts)
		if err != nil {
			return fmt.Errorf("failed to update onboarding issue: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(onboardCmd)

	// Add flags
	onboardCmd.Flags().StringSliceP("onboard-labels", "l", onboard.DefaultOptions().OnboardLabels, "Labels used to find issues suitable for community contribution (e.g., 'good first issue', 'help wanted')")
	onboardCmd.Flags().StringSliceP("difficulty-labels", "d", onboard.DefaultOptions().DifficultyLabels, "Labels used to indicate issue difficulty, ordered from easy to hard (e.g., 'difficulty/easy', 'difficulty/medium')")
	onboardCmd.Flags().StringSliceP("category-labels", "c", onboard.DefaultOptions().CategoryLabels, "Labels used to classify issues by type within each difficulty level (e.g., 'bug', 'feature')")
	onboardCmd.Flags().BoolP("dry-run", "n", false, "Preview the changes without modifying any issues")
	onboardCmd.Flags().BoolP("yes", "y", false, "Automatically apply changes without confirmation")
	onboardCmd.Flags().String("target-label", onboard.DefaultOnboardOptions().TargetLabel, "Label used to locate the issue where onboarding content will be updated")
}

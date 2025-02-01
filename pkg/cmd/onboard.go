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
	Short: "Generate community onboarding issues",
	Long: `Generate community onboarding issues based on issue labels.
This command will generate a list of issues suitable for new contributors,
based on issues with help wanted or good first issue labels.

By default, it uses "help wanted" and "good first issue" as help labels.
You can customize the help labels using the -l, --help-labels flag.

You can also specify difficulty labels using the -d, --difficulty-labels flag.
For example: -d "Easy,Medium,Hard"

You can also specify categories to group issues by using the -c, --categories flag.
For example: -c "Category1,Category2,Category3"`,
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
		helpLabels, err := cmd.Flags().GetStringSlice("help-labels")
		if err != nil {
			return err
		}
		difficultyLabels, err := cmd.Flags().GetStringSlice("difficulty-labels")
		if err != nil {
			return err
		}
		categories, err := cmd.Flags().GetStringSlice("categories")
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
		onboardingLabel, err := cmd.Flags().GetString("onboarding-label")
		if err != nil {
			return err
		}
		log.Debug("Help labels: [%s]", strings.Join(helpLabels, ", "))
		log.Debug("Difficulty labels: [%s]", strings.Join(difficultyLabels, ", "))
		log.Debug("Categories: [%s]", strings.Join(categories, ", "))

		// Create GitHub client
		client, err := api.DefaultRESTClient()
		if err != nil {
			return fmt.Errorf("failed to create GitHub client: %w", err)
		}

		// Create options
		opts := onboard.Options{
			HelpLabels:       helpLabels,
			DifficultyLabels: difficultyLabels,
			Categories:       categories,
		}
		onboardOpts := onboard.OnboardOptions{
			OnboardingLabel: onboardingLabel,
			DryRun:          dryRun,
			AutoConfirm:     autoConfirm,
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
	onboardCmd.Flags().StringSliceP("help-labels", "l", []string{"good first issue", "help wanted"}, "Help labels")
	onboardCmd.Flags().StringSliceP("difficulty-labels", "d", []string{"good first issue", "help wanted"}, "Difficulty labels, from easy to hard")
	onboardCmd.Flags().StringSliceP("categories", "c", []string{"bug", "documentation", "enhancement"}, "Categories to group issues by")
	onboardCmd.Flags().BoolP("dry-run", "n", false, "Only show what would be done without making changes")
	onboardCmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompt")
	onboardCmd.Flags().String("onboarding-label", "onboarding", "Label used to identify onboarding issues")
}

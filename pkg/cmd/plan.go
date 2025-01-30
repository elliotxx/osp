package cmd

import (
	"fmt"
	"strings"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/elliotxx/osp/pkg/config"
	"github.com/elliotxx/osp/pkg/log"
	"github.com/elliotxx/osp/pkg/planning"
	"github.com/elliotxx/osp/pkg/repo"
	"github.com/spf13/cobra"
)

func newPlanCmd() *cobra.Command {
	var (
		planningLabel string
		categories    []string
		excludePR     bool
		dryRun        bool
		autoConfirm   bool
	)

	cmd := &cobra.Command{
		Use:   "plan [milestone-number]",
		Short: "Generate and update community planning",
		Long: `Generate and update community planning based on milestone issues.
This command will create or update a planning issue that summarizes all issues
in the specified milestone or all open milestones, categorized by their labels.

If no milestone number is provided, it will scan all open milestones.

By default, it will show the preview of the planning content and ask for confirmation
before creating or updating the planning issue.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get GitHub client
			client, err := api.DefaultRESTClient()
			if err != nil {
				return fmt.Errorf("failed to create GitHub client: %w", err)
			}

			// Get config
			cfg, err := config.Load("")
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// Get current repository
			repoManager := repo.NewManager(cfg)
			currentRepo := repoManager.Current()
			if currentRepo == "" {
				return fmt.Errorf("no repository selected, please use 'osp repo current' to select one")
			}

			// Parse owner and repo from current repository
			parts := strings.Split(currentRepo, "/")
			if len(parts) != 2 {
				return fmt.Errorf("invalid repository format: %s", currentRepo)
			}
			owner, repoName := parts[0], parts[1]

			// Create plan manager
			manager := planning.NewManager(client)

			// Create options
			opts := planning.Options{
				PlanningLabel: planningLabel,
				Categories:    categories,
				ExcludePR:     excludePR,
				DryRun:        dryRun,
				AutoConfirm:   autoConfirm,
			}

			// If milestone number is provided, update that specific milestone
			if len(args) > 0 {
				var milestoneNumber int
				_, err := fmt.Sscanf(args[0], "%d", &milestoneNumber)
				if err != nil {
					return fmt.Errorf("invalid milestone number: %w", err)
				}

				return manager.Update(cmd.Context(), owner, repoName, milestoneNumber, opts)
			}

			// Otherwise, get all open milestones and update their planning
			milestones, err := manager.ListOpenMilestones(cmd.Context(), owner, repoName)
			if err != nil {
				return fmt.Errorf("failed to list open milestones: %w", err)
			}

			if len(milestones) == 0 {
				log.Info("No open milestones found")
				return nil
			}

			log.Info("Found %d open milestones", len(milestones))
			for _, m := range milestones {
				if err := manager.Update(cmd.Context(), owner, repoName, m.Number, opts); err != nil {
					log.Error("Failed to update planning for milestone %d: %v", m.Number, err)
					continue
				}
			}

			return nil
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&planningLabel, "label", "l", planning.DefaultOptions().PlanningLabel, "Label to use for planning issues")
	cmd.Flags().StringSliceVarP(&categories, "categories", "c", planning.DefaultOptions().Categories, "Categories to group issues by")
	cmd.Flags().BoolVarP(&excludePR, "exclude-pr", "e", planning.DefaultOptions().ExcludePR, "Exclude pull requests from planning")
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "d", planning.DefaultOptions().DryRun, "Only show what would be done without making actual changes")
	cmd.Flags().BoolVarP(&autoConfirm, "yes", "y", planning.DefaultOptions().AutoConfirm, "Skip confirmation and update automatically")

	return cmd
}

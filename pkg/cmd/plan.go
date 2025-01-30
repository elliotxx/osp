package cmd

import (
	"fmt"
	"strings"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/elliotxx/osp/pkg/config"
	"github.com/elliotxx/osp/pkg/planning"
	"github.com/elliotxx/osp/pkg/repo"
	"github.com/spf13/cobra"
)

func newPlanCmd() *cobra.Command {
	var (
		planningLabel string
		categories    []string
		excludePR     bool
	)

	cmd := &cobra.Command{
		Use:   "plan [milestone-number]",
		Short: "Generate and update community planning",
		Long: `Generate and update community planning based on milestone issues.
This command will create or update a planning issue that summarizes all issues
in the specified milestone, categorized by their labels.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Parse milestone number
			var milestoneNumber int
			_, err := fmt.Sscanf(args[0], "%d", &milestoneNumber)
			if err != nil {
				return fmt.Errorf("invalid milestone number: %w", err)
			}

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
			}

			// Update planning
			err = manager.Update(cmd.Context(), owner, repoName, milestoneNumber, opts)
			if err != nil {
				return err
			}

			fmt.Println("✓ Successfully updated planning")
			return nil
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&planningLabel, "label", "l", "planning", "Label for planning issues")
	cmd.Flags().StringSliceVarP(&categories, "categories", "c", []string{"bug", "documentation", "enhancement"}, "Categories for issues")
	cmd.Flags().BoolVarP(&excludePR, "exclude-pr", "e", true, "Exclude pull requests from the summary")

	return cmd
}

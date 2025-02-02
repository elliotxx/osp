package cmd

import (
	"fmt"
	"strings"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/elliotxx/osp/pkg/auth"
	"github.com/elliotxx/osp/pkg/config"
	"github.com/elliotxx/osp/pkg/log"
	"github.com/elliotxx/osp/pkg/planning"
	"github.com/elliotxx/osp/pkg/repo"
	"github.com/spf13/cobra"
)

var (
	planningLabel string
	targetTitle   string
	categories    []string
	priorities    []string
	excludePR     bool
	dryRun        bool
	autoConfirm   bool
)

func newPlanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plan [milestone-number]",
		Short: "Generate and update community planning",
		Long: `Generate and maintain planning content to track milestone progress.

This command will create or update an issue that summarizes all issues in a milestone,
organized by priority and category. The content is designed to help track milestone
progress and highlight high-priority tasks.

If no milestone number is provided, it will scan all open milestones.

Available fields in title template:
  .Title       - Milestone title (e.g., "v1.0.0")
  .Description - Milestone description
  .Number      - Milestone number (e.g., 1)
  .State       - Milestone state (e.g., "open" or "closed")
  .DueOn       - Milestone due date (e.g., "2025-12-31T23:59:59Z")
  .HTMLURL     - Milestone URL on GitHub

Examples:
  # Update planning content for all open milestones
  osp plan

  # Update planning content for milestone #1
  osp plan 1

  # Use custom category labels
  osp plan --category-labels="bug,feature,documentation"

  # Use custom priority labels
  osp plan --priority-labels="priority/high,priority/medium,priority/low"

  # Preview changes without updating any issues
  osp plan --dry-run

  # Update automatically without confirmation
  osp plan --yes

  # Specify a custom label for the target issue
  osp plan --target-label="milestone-plan"

  # Specify a custom title template for the target issue
  osp plan --target-title="Planning: {{ .Title }}"

  # Use milestone fields in title template
  osp plan --target-title="Planning for {{ .Title }} (Due: {{ .DueOn.Format \"2006-01-02\" }})"

  # Exclude pull requests from planning content
  osp plan --exclude-pr`,
		Args: cobra.MaximumNArgs(1),
		RunE: runPlanUpdate,
	}

	// Add flags
	cmd.Flags().StringVarP(&planningLabel, "target-label", "t", planning.DefaultOptions().PlanningLabel, "Label used to locate the issue where planning content will be updated")
	cmd.Flags().StringVarP(&targetTitle, "target-title", "T", planning.DefaultOptions().TargetTitle, "Title template of the target issue where planning content will be updated")
	cmd.Flags().StringSliceVarP(&categories, "category-labels", "c", planning.DefaultOptions().Categories, "Labels used to classify issues by type (e.g., 'bug', 'feature')")
	cmd.Flags().StringSliceVarP(&priorities, "priority-labels", "p", planning.DefaultOptions().Priorities, "Labels used to indicate issue priority, ordered from high to low (e.g., 'priority/high', 'priority/medium')")
	cmd.Flags().BoolVarP(&excludePR, "exclude-pr", "e", planning.DefaultOptions().ExcludePR, "Exclude pull requests from planning content")
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "n", planning.DefaultOptions().DryRun, "Preview the changes without modifying any issues")
	cmd.Flags().BoolVarP(&autoConfirm, "yes", "y", planning.DefaultOptions().AutoConfirm, "Automatically apply changes without confirmation")

	return cmd
}

func runPlanUpdate(cmd *cobra.Command, args []string) error {
	// Check authentication
	if err := auth.CheckAuth(); err != nil {
		return err
	}

	// Load config
	cfg, err := config.Load("")
	if err != nil {
		return err
	}

	// Get GitHub client
	client, err := api.DefaultRESTClient()
	if err != nil {
		return fmt.Errorf("failed to create GitHub client: %w", err)
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
		TargetTitle:   targetTitle,
		Categories:    categories,
		Priorities:    priorities,
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
}

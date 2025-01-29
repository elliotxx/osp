package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elliotxx/osp/internal/config"
	"github.com/elliotxx/osp/internal/plan"
	"github.com/spf13/cobra"
)

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Project planning",
	Long:  `Manage project plans, including milestones and issues.`,
}

var planGenerateCmd = &cobra.Command{
	Use:   "generate [owner/repo]",
	Short: "Generate project plan",
	Long:  `Generate a project plan based on milestones and issues.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load("")
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Get repository name
		repoName := cfg.Current
		if len(args) > 0 {
			repoName = args[0]
		}
		if repoName == "" {
			return fmt.Errorf("no repository specified and no current repository set")
		}

		// Get flags
		format, _ := cmd.Flags().GetString("format")
		includeIssues, _ := cmd.Flags().GetBool("include-issues")

		planManager := plan.NewManager(cfg)
		plan, err := planManager.Generate(cmd.Context(), repoName, includeIssues)
		if err != nil {
			return fmt.Errorf("failed to generate plan: %w", err)
		}

		// Output plan
		switch strings.ToLower(format) {
		case "json":
			data, err := json.MarshalIndent(plan, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(data))

		default:
			fmt.Printf("Project Plan for %s:\n\n", repoName)
			for _, milestone := range plan.Milestones {
				fmt.Printf("Milestone: %s\n", milestone.Title)
				fmt.Printf("Due Date: %s\n", milestone.DueDate)
				fmt.Printf("Progress: %d%%\n", milestone.Progress)

				if includeIssues && len(milestone.Issues) > 0 {
					fmt.Println("\nIssues:")
					for _, issue := range milestone.Issues {
						fmt.Printf("- %s (#%d)\n", issue.Title, issue.Number)
						fmt.Printf("  Status: %s\n", issue.Status)
						if issue.Assignee != "" {
							fmt.Printf("  Assignee: %s\n", issue.Assignee)
						}
					}
				}
				fmt.Println()
			}
		}

		return nil
	},
}

var planListCmd = &cobra.Command{
	Use:   "list [owner/repo]",
	Short: "List project plans",
	Long:  `List all project plans.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load("")
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Get repository name
		repoName := cfg.Current
		if len(args) > 0 {
			repoName = args[0]
		}
		if repoName == "" {
			return fmt.Errorf("no repository specified and no current repository set")
		}

		// Get flags
		format, _ := cmd.Flags().GetString("format")

		planManager := plan.NewManager(cfg)
		plans, err := planManager.List(cmd.Context(), repoName)
		if err != nil {
			return fmt.Errorf("failed to list plans: %w", err)
		}

		// Output plans
		switch strings.ToLower(format) {
		case "json":
			data, err := json.MarshalIndent(plans, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(data))

		default:
			if len(plans) == 0 {
				fmt.Println("No plans found.")
				return nil
			}

			fmt.Printf("Project Plans for %s:\n\n", repoName)
			for _, p := range plans {
				fmt.Printf("Plan: %s\n", p.Name)
				fmt.Printf("Created: %s\n", p.CreatedAt)
				fmt.Printf("Status: %s\n", p.Status)
				fmt.Printf("Progress: %d%%\n\n", p.Progress)
			}
		}

		return nil
	},
}

func init() {
	// Add plan commands
	rootCmd.AddCommand(planCmd)
	planCmd.AddCommand(planGenerateCmd)
	planCmd.AddCommand(planListCmd)

	// Add flags
	planGenerateCmd.Flags().String("format", "text", "Output format (text, json)")
	planGenerateCmd.Flags().Bool("include-issues", false, "Include issues in the plan")
	planListCmd.Flags().String("format", "text", "Output format (text, json)")
}

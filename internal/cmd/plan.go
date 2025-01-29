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
	Use:   "plan [owner/repo]",
	Short: "Manage project plans",
	Long:  `Generate and manage project plans based on milestones and issues.`,
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
		update, _ := cmd.Flags().GetBool("update")
		format, _ := cmd.Flags().GetString("format")
		auto, _ := cmd.Flags().GetBool("auto")

		planManager := plan.NewManager(cfg)

		if update {
			// TODO: Implement plan update logic
			return fmt.Errorf("plan update not implemented yet")
		}

		// Generate plan
		plan, err := planManager.Generate(cmd.Context(), repoName)
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

		case "markdown", "md":
			fmt.Printf("# %s\n\n", plan.Title)
			fmt.Printf("%s\n\n", plan.Description)
			fmt.Printf("Period: %s - %s\n\n", 
				plan.StartDate.Format("2006-01-02"),
				plan.EndDate.Format("2006-01-02"))

			fmt.Println("## Milestones\n")
			for _, m := range plan.Milestones {
				fmt.Printf("### %s\n", m.Title)
				fmt.Printf("- Due: %s\n", m.DueDate.Format("2006-01-02"))
				fmt.Printf("- Status: %s\n", m.State)
				if m.Description != "" {
					fmt.Printf("\n%s\n", m.Description)
				}
				
				if len(m.Issues) > 0 {
					fmt.Println("\nIssues:")
					for _, i := range m.Issues {
						status := "❌"
						if i.State == "closed" {
							status = "✓"
						}
						fmt.Printf("- [%s] #%d %s\n", status, i.Number, i.Title)
					}
				}
				fmt.Println()
			}

		default:
			fmt.Printf("Project Plan: %s\n\n", plan.Title)
			fmt.Printf("Description: %s\n", plan.Description)
			fmt.Printf("Period: %s - %s\n\n",
				plan.StartDate.Format("2006-01-02"),
				plan.EndDate.Format("2006-01-02"))

			fmt.Println("Milestones:")
			for _, m := range plan.Milestones {
				fmt.Printf("\n- %s\n", m.Title)
				fmt.Printf("  Due: %s\n", m.DueDate.Format("2006-01-02"))
				fmt.Printf("  Status: %s\n", m.State)
				
				if len(m.Issues) > 0 {
					fmt.Println("  Issues:")
					for _, i := range m.Issues {
						status := "[ ]"
						if i.State == "closed" {
							status = "[x]"
						}
						fmt.Printf("  - %s #%d %s\n", status, i.Number, i.Title)
					}
				}
			}
		}

		return nil
	},
}

func init() {
	planCmd.Flags().Bool("update", false, "Update existing plan")
	planCmd.Flags().String("format", "text", "Output format (text, json, markdown)")
	planCmd.Flags().Bool("auto", false, "Auto-generate suggestions")
}

package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elliotxx/osp/internal/config"
	"github.com/elliotxx/osp/internal/task"
	"github.com/spf13/cobra"
)

var taskCmd = &cobra.Command{
	Use:   "task [owner/repo]",
	Short: "Manage community tasks",
	Long:  `List and generate community tasks (issues) for the repository.`,
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
		taskType, _ := cmd.Flags().GetString("type")
		format, _ := cmd.Flags().GetString("format")
		auto, _ := cmd.Flags().GetBool("auto")

		taskManager := task.NewManager(cfg)

		if auto {
			if err := taskManager.Generate(cmd.Context(), repoName, taskType); err != nil {
				return fmt.Errorf("failed to generate tasks: %w", err)
			}
			fmt.Println("✓ Successfully generated new tasks")
			return nil
		}

		// List tasks
		tasks, err := taskManager.List(cmd.Context(), repoName, taskType)
		if err != nil {
			return fmt.Errorf("failed to list tasks: %w", err)
		}

		// Output tasks
		switch strings.ToLower(format) {
		case "json":
			data, err := json.MarshalIndent(tasks, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(data))

		case "markdown", "md":
			fmt.Printf("# Tasks for %s\n\n", repoName)
			
			if taskType != "" {
				fmt.Printf("Type: %s\n\n", taskType)
			}

			for _, t := range tasks {
				fmt.Printf("## #%d: %s\n\n", t.Number, t.Title)
				
				if t.Body != "" {
					fmt.Printf("%s\n\n", t.Body)
				}

				fmt.Printf("- State: %s\n", t.State)
				fmt.Printf("- Type: %s\n", t.Type)
				fmt.Printf("- Difficulty: %s\n", t.Difficulty)
				
				if len(t.Labels) > 0 {
					fmt.Printf("- Labels: %s\n", strings.Join(t.Labels, ", "))
				}
				
				if len(t.Assignees) > 0 {
					fmt.Printf("- Assignees: %s\n", strings.Join(t.Assignees, ", "))
				}
				
				fmt.Printf("- Created: %s\n", t.CreatedAt.Format("2006-01-02"))
				fmt.Printf("- Updated: %s\n", t.UpdatedAt.Format("2006-01-02"))
				fmt.Println()
			}

		default:
			fmt.Printf("Tasks for %s:\n\n", repoName)
			
			if len(tasks) == 0 {
				fmt.Println("No tasks found.")
				return nil
			}

			for _, t := range tasks {
				fmt.Printf("#%d: %s\n", t.Number, t.Title)
				fmt.Printf("  State: %s\n", t.State)
				fmt.Printf("  Type: %s\n", t.Type)
				fmt.Printf("  Difficulty: %s\n", t.Difficulty)
				
				if len(t.Labels) > 0 {
					fmt.Printf("  Labels: %s\n", strings.Join(t.Labels, ", "))
				}
				
				if len(t.Assignees) > 0 {
					fmt.Printf("  Assignees: %s\n", strings.Join(t.Assignees, ", "))
				}
				
				fmt.Printf("  Created: %s\n", t.CreatedAt.Format("2006-01-02"))
				fmt.Printf("  Updated: %s\n", t.UpdatedAt.Format("2006-01-02"))
				fmt.Println()
			}
		}

		return nil
	},
}

func init() {
	taskCmd.Flags().String("type", "", "Task type (good-first-issue, help-wanted, bug, enhancement)")
	taskCmd.Flags().String("format", "text", "Output format (text, json, markdown)")
	taskCmd.Flags().Bool("auto", false, "Auto-generate new tasks")
}

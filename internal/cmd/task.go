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
	Use:   "task",
	Short: "Task management",
	Long:  `Manage community tasks, including generating and listing tasks.`,
}

var taskGenerateCmd = &cobra.Command{
	Use:   "generate [owner/repo]",
	Short: "Generate tasks",
	Long:  `Generate tasks based on repository analysis.`,
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
		category, _ := cmd.Flags().GetString("category")

		taskManager := task.NewManager(cfg)
		tasks, err := taskManager.Generate(cmd.Context(), repoName, category)
		if err != nil {
			return fmt.Errorf("failed to generate tasks: %w", err)
		}

		// Output tasks
		switch strings.ToLower(format) {
		case "json":
			data, err := json.MarshalIndent(tasks, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(data))

		default:
			fmt.Printf("Generated Tasks for %s:\n\n", repoName)
			for _, t := range tasks {
				fmt.Printf("Title: %s\n", t.Title)
				fmt.Printf("Category: %s\n", t.Category)
				fmt.Printf("Priority: %s\n", t.Priority)
				fmt.Printf("Estimated Time: %s\n", t.EstimatedTime)
				if t.Description != "" {
					fmt.Printf("\nDescription:\n%s\n", t.Description)
				}
				fmt.Println()
			}
		}

		return nil
	},
}

var taskListCmd = &cobra.Command{
	Use:   "list [owner/repo]",
	Short: "List tasks",
	Long:  `List all tasks for a repository.`,
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
		status, _ := cmd.Flags().GetString("status")
		category, _ := cmd.Flags().GetString("category")

		taskManager := task.NewManager(cfg)
		tasks, err := taskManager.List(cmd.Context(), repoName, status, category)
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

		default:
			if len(tasks) == 0 {
				fmt.Println("No tasks found.")
				return nil
			}

			fmt.Printf("Tasks for %s:\n\n", repoName)
			for _, t := range tasks {
				fmt.Printf("Title: %s\n", t.Title)
				fmt.Printf("Status: %s\n", t.Status)
				fmt.Printf("Category: %s\n", t.Category)
				fmt.Printf("Priority: %s\n", t.Priority)
				if t.Assignee != "" {
					fmt.Printf("Assignee: %s\n", t.Assignee)
				}
				fmt.Println()
			}
		}

		return nil
	},
}

func init() {
	// Add task commands
	rootCmd.AddCommand(taskCmd)
	taskCmd.AddCommand(taskGenerateCmd)
	taskCmd.AddCommand(taskListCmd)

	// Add flags
	taskGenerateCmd.Flags().String("format", "text", "Output format (text, json)")
	taskGenerateCmd.Flags().String("category", "", "Task category filter")
	
	taskListCmd.Flags().String("format", "text", "Output format (text, json)")
	taskListCmd.Flags().String("status", "", "Task status filter")
	taskListCmd.Flags().String("category", "", "Task category filter")
}

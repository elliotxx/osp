package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elliotxx/osp/pkg/config"
	"github.com/elliotxx/osp/pkg/stats"
	"github.com/spf13/cobra"
)

const (
	outputFormatJSON = "json"
)

var statsCmd = &cobra.Command{
	Use:   "stats [repository]",
	Short: "Show repository statistics",
	Long:  "Show repository statistics such as stars, forks, and open issues",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get repository name from args or current
		var repoName string
		if len(args) > 0 {
			repoName = args[0]
		} else {
			state, err := config.LoadState()
			if err != nil {
				return fmt.Errorf("failed to load state: %w", err)
			}
			repoName = state.Current
		}

		// Get format
		format, _ := cmd.Flags().GetString("format")

		// Create stats manager
		manager, err := stats.NewManager()
		if err != nil {
			return err
		}

		// Get stats
		stats, err := manager.Get(context.Background(), repoName)
		if err != nil {
			return err
		}

		// Output stats
		switch strings.ToLower(format) {
		case outputFormatJSON:
			data, err := json.MarshalIndent(stats, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(data))

		default:
			fmt.Printf("Repository: %s\n", repoName)
			fmt.Printf("Stars: %d\n", stats.Stars)
			fmt.Printf("Forks: %d\n", stats.Forks)
			fmt.Printf("Open Issues: %d\n", stats.OpenIssues)
			fmt.Printf("Last Updated: %s\n", stats.LastUpdated)
		}

		return nil
	},
}

var starCmd = &cobra.Command{
	Use:   "star",
	Short: "Star related commands",
	Long:  `Commands related to repository stars, including history and analysis.`,
}

var starHistoryCmd = &cobra.Command{
	Use:   "history [owner/repo]",
	Short: "Show star history",
	Long:  `Show the history of stars for a repository over time.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get repository name from args or current
		var repoName string
		if len(args) > 0 {
			repoName = args[0]
		} else {
			state, err := config.LoadState()
			if err != nil {
				return fmt.Errorf("failed to load state: %w", err)
			}
			repoName = state.Current
		}

		// Get flags
		days, _ := cmd.Flags().GetInt("days")
		format, _ := cmd.Flags().GetString("format")

		// Create stats manager
		manager, err := stats.NewManager()
		if err != nil {
			return err
		}

		// Get star history
		history, err := manager.GetStarHistory(context.Background(), repoName, days)
		if err != nil {
			return err
		}

		// Output history
		switch strings.ToLower(format) {
		case outputFormatJSON:
			data, err := json.MarshalIndent(history, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(data))

		default:
			fmt.Printf("Star history for %s (last %d days):\n\n", repoName, days)
			for _, h := range history {
				fmt.Printf("%s: %d stars\n", h.Date.Format("2006-01-02"), h.Stars)
			}
		}

		return nil
	},
}

func init() {
	// Add stats commands
	rootCmd.AddCommand(statsCmd)

	// Add star commands
	rootCmd.AddCommand(starCmd)
	starCmd.AddCommand(starHistoryCmd)

	// Add flags
	statsCmd.Flags().String("format", "text", "Output format (text, json)")
	starHistoryCmd.Flags().Int("days", 30, "Number of days to show history for")
	starHistoryCmd.Flags().String("format", "text", "Output format (text, json)")
}

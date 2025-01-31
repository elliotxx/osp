package cmd

import (
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
	Use:   "stats",
	Short: "Repository statistics",
	Long:  `View and analyze repository statistics.`,
}

var statsShowCmd = &cobra.Command{
	Use:   "show [owner/repo]",
	Short: "Show repository statistics",
	Long:  `Show basic statistics for a repository, including stars, forks, issues, etc.`,
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

		// Get format
		format, _ := cmd.Flags().GetString("format")

		statsManager := stats.NewManager(cfg)
		stats, err := statsManager.Get(cmd.Context(), repoName)
		if err != nil {
			return fmt.Errorf("failed to get statistics: %w", err)
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
			fmt.Printf("Statistics for %s:\n\n", repoName)
			fmt.Printf("Stars:        %d\n", stats.Stars)
			fmt.Printf("Forks:        %d\n", stats.Forks)
			fmt.Printf("Open Issues:  %d\n", stats.OpenIssues)
			fmt.Printf("Last Update:  %s\n", stats.LastUpdated)
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
		days, _ := cmd.Flags().GetInt("days")
		format, _ := cmd.Flags().GetString("format")

		statsManager := stats.NewManager(cfg)
		history, err := statsManager.GetStarHistory(cmd.Context(), repoName, days)
		if err != nil {
			return fmt.Errorf("failed to get star history: %w", err)
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
	statsCmd.AddCommand(statsShowCmd)

	// Add star commands
	rootCmd.AddCommand(starCmd)
	starCmd.AddCommand(starHistoryCmd)

	// Add flags
	statsShowCmd.Flags().String("format", "text", "Output format (text, json)")
	starHistoryCmd.Flags().Int("days", 30, "Number of days to show history for")
	starHistoryCmd.Flags().String("format", "text", "Output format (text, json)")
}

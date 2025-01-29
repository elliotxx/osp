package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/elliotxx/osp/internal/config"
	"github.com/elliotxx/osp/internal/stats"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats [owner/repo]",
	Short: "Show repository statistics",
	Long:  `Display statistics for the specified repository or current repository.`,
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
		stats, err := statsManager.GetStats(cmd.Context(), repoName)
		if err != nil {
			return fmt.Errorf("failed to get statistics: %w", err)
		}

		switch strings.ToLower(format) {
		case "json":
			data, err := json.MarshalIndent(stats, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(data))
		default:
			fmt.Printf("Statistics for %s:\n", repoName)
			fmt.Printf("Stars:        %d\n", stats.Stars)
			fmt.Printf("Forks:        %d\n", stats.Forks)
			fmt.Printf("Issues:       %d\n", stats.Issues)
			fmt.Printf("Pull Requests: %d\n", stats.PRs)
			fmt.Printf("Last Updated: %s\n", stats.LastUpdate.Format(time.RFC3339))
		}

		return nil
	},
}

var starCmd = &cobra.Command{
	Use:   "star",
	Short: "Star related commands",
	Long:  `Commands related to repository stars.`,
}

var starHistoryCmd = &cobra.Command{
	Use:   "history [owner/repo]",
	Short: "Show star history",
	Long:  `Display star history for the specified repository or current repository.`,
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

		// Get time range
		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")
		format, _ := cmd.Flags().GetString("format")

		fromTime, err := parseTime(from)
		if err != nil {
			return fmt.Errorf("invalid from date: %w", err)
		}

		toTime, err := parseTime(to)
		if err != nil {
			return fmt.Errorf("invalid to date: %w", err)
		}

		statsManager := stats.NewManager(cfg)
		history, err := statsManager.GetStarHistory(cmd.Context(), repoName, fromTime, toTime)
		if err != nil {
			return fmt.Errorf("failed to get star history: %w", err)
		}

		switch strings.ToLower(format) {
		case "json":
			data, err := json.MarshalIndent(history, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(data))
		default:
			fmt.Printf("Star history for %s:\n", repoName)
			for _, h := range history {
				fmt.Printf("%s: %d stars\n", h.Date.Format("2006-01-02"), h.Stars)
			}
		}

		return nil
	},
}

func init() {
	// Stats command flags
	statsCmd.Flags().String("format", "text", "Output format (text, json)")

	// Star history command flags
	starHistoryCmd.Flags().String("from", "", "Start date (YYYY-MM-DD)")
	starHistoryCmd.Flags().String("to", "", "End date (YYYY-MM-DD)")
	starHistoryCmd.Flags().String("format", "text", "Output format (text, json)")

	// Add star subcommands
	starCmd.AddCommand(starHistoryCmd)
}

func parseTime(date string) (time.Time, error) {
	if date == "" {
		return time.Now(), nil
	}
	return time.Parse("2006-01-02", date)
}

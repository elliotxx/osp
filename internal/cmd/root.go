package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "osp",
	Short: "OSP - Open Source Pilot",
	Long: `OSP (Open Source Pilot) is an automated open source software management tool.
It helps maintainers manage projects, track progress, and generate reports more efficiently.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().Bool("debug", false, "enable debug mode")
	rootCmd.PersistentFlags().Bool("quiet", false, "enable quiet mode")
	rootCmd.PersistentFlags().String("config", "", "config file (default is $HOME/.config/osp/config.yml)")

	// Add commands
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(statsCmd)
	rootCmd.AddCommand(planCmd)
	rootCmd.AddCommand(taskCmd)
	rootCmd.AddCommand(activityCmd)
	rootCmd.AddCommand(starCmd)
}

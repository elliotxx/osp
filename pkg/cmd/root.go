package cmd

import (
	"os"

	"github.com/elliotxx/osp/pkg/log"
	"github.com/spf13/cobra"
)

var (
	verbose bool
	noColor bool
	rootCmd = &cobra.Command{
		Use:   "osp",
		Short: "Open Source Project Management Tool",
		Long: `OSP is a command-line tool for managing open source projects.
It helps you manage issues, milestones, planning, and more.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			log.SetVerbose(verbose)
			log.SetNoColor(noColor)
		},
	}
)

func init() {
	// Add global flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable color output")

	rootCmd.AddCommand(
		newAuthCmd(),
		newPlanCmd(),
	)
}

// Execute executes the root command
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		log.Error("%v", err)
		os.Exit(1)
	}
	return nil
}

// GetRootCmd returns the root cobra command
func GetRootCmd() *cobra.Command {
	return rootCmd
}

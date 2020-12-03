package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "vergo",
		Short: "vergo [command]",
	}
)

func init() {
	rootCmd.PersistentFlags().String("log-level", "Info", "set log level")
	rootCmd.PersistentFlags().String("tag-prefix", "", "version prefix")
	rootCmd.PersistentFlags().String("repository-location", ".", "repository location")
	rootCmd.PersistentFlags().Bool("dry-run", false, "dry run")
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

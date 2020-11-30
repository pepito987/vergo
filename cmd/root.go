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
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version string
	Commit  string
	Date    string
	BuiltBy string
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of command",
	Long:  `All software has versions.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("version: %s\n", Version)
		fmt.Printf("commit : %s\n", Commit)
		fmt.Printf("date: %s\n", Date)
		fmt.Printf("builtBy: %s\n", BuiltBy)
	},
}

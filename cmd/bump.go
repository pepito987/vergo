package cmd

import (
	"github.com/spf13/cobra"
)

var bumpCmd = &cobra.Command{
	Use:       "bump (patch|minor|major)",
	Short:     "bumps the version numbers",
	Args:      cobra.ExactValidArgs(1),
	ValidArgs: []string{"patch", "minor", "major"},
	RunE: func(cmd *cobra.Command, args []string) error {
		prefix, err := cmd.Flags().GetString("tag-prefix")
		CheckIfError(err)
		repoLocation, err := cmd.Flags().GetString("repository-location")
		CheckIfError(err)
		Info(prefix, repoLocation)
		return nil
	},
}

func init() {
	bumpCmd.Flags().Bool("push-tag", false, "push the new tag")
	bumpCmd.Flags().String("tag-prefix", "", "version prefix")
	bumpCmd.Flags().String("repository-location", ".", "repository location")
	bumpCmd.Flags().String("public-key-location", ".ssh/", "public key location")
	rootCmd.AddCommand(bumpCmd)
}

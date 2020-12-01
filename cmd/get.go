package cmd

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	vergo "github.com/inanme/vergo/git"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "gets the version",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		setLogger(cmd)
		prefix, err := cmd.Flags().GetString("tag-prefix")
		prefix = sanitiseTagPrefix(prefix)
		CheckIfError(err)
		repoLocation, err := cmd.Flags().GetString("repository-location")
		CheckIfError(err)
		repo, err := git.PlainOpen(repoLocation)
		CheckIfError(err)
		version, err := vergo.LatestRef(repo, prefix)
		CheckIfError(err)
		fmt.Print(version.Version.String())
		return nil
	},
}

func init() {
	getCmd.Flags().String("repository-location", ".", "repository location")
	rootCmd.AddCommand(getCmd)
}

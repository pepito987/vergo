package cmd

import (
	"github.com/go-git/go-git/v5"
	vergo "github.com/inanme/vergo/git"
	"github.com/spf13/cobra"
)

var push = &cobra.Command{
	Use:   "push",
	Short: "push the latest tag to remote",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		setLogger(cmd)
		dryRun, err := cmd.Flags().GetBool("dry-run")
		CheckIfError(err)
		socket, err := checkAuthSocket(true)
		CheckIfError(err)
		prefix, err := cmd.Flags().GetString("tag-prefix")
		prefix = sanitiseTagPrefix(prefix)
		CheckIfError(err)
		repoLocation, err := cmd.Flags().GetString("repository-location")
		CheckIfError(err)
		repo, err := git.PlainOpen(repoLocation)
		CheckIfError(err)
		version, err := vergo.LatestRef(repo, prefix)
		CheckIfError(err)
		err = vergo.PushTag(repo, socket, version.Version, prefix, dryRun)
		CheckIfError(err)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(push)
}

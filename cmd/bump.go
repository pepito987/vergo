package cmd

import (
	"github.com/go-git/go-git/v5"
	"github.com/inanme/vergo/bump"
	vergo "github.com/inanme/vergo/git"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var bumpCmd = &cobra.Command{
	Use:       "bump (patch|minor|major)",
	Short:     "bumps the version numbers",
	Args:      cobra.ExactValidArgs(1),
	ValidArgs: []string{"patch", "minor", "major"},
	RunE: func(cmd *cobra.Command, args []string) error {
		setLogger(cmd)
		socket, err := checkAuthSocket(cmd)
		CheckIfError(err)
		prefix, err := cmd.Flags().GetString("tag-prefix")
		prefix = sanitiseTagPrefix(prefix)
		CheckIfError(err)
		repoLocation, err := cmd.Flags().GetString("repository-location")
		CheckIfError(err)
		repo, err := git.PlainOpen(repoLocation)
		CheckIfError(err)
		version, err := bump.Bump(repo, prefix, args[0])
		CheckIfError(err)
		pushTag, err := cmd.Flags().GetBool("push-tag")
		CheckIfError(err)
		if pushTag {
			err = vergo.PushTag(repo, socket, version, prefix)
			CheckIfError(err)
		} else {
			log.Trace("Push not enabled")
		}
		return nil
	},
}

func init() {
	bumpCmd.Flags().Bool("push-tag", false, "push the new tag")
	bumpCmd.Flags().String("repository-location", ".", "repository location")
	rootCmd.AddCommand(bumpCmd)
}

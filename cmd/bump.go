package cmd

import (
	"github.com/go-git/go-git/v5"
	"github.com/inanme/vergo/bump"
	vergo "github.com/inanme/vergo/git"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
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
		repo, err := git.PlainOpen(repoLocation)
		CheckIfError(err)
		version, err := bump.Bump(repo, prefix, args[0])
		CheckIfError(err)
		pushTag, err := cmd.Flags().GetBool("push-tag")
		CheckIfError(err)
		if pushTag {
			publicKeyLocation, err := cmd.Flags().GetString("public-key-location")
			CheckIfError(err)
			log.Debugf("Public key location %v", publicKeyLocation)
			err = vergo.PushTag(repo, publicKeyLocation, version.String())
			CheckIfError(err)
		} else {
			log.Trace("Push not enabled")
		}
		return nil
	},
}

func init() {
	homedir, err := homedir.Dir()
	if err != nil {
		log.WithError(err).Errorln("can not find homedir")
	}
	bumpCmd.Flags().Bool("push-tag", false, "push the new tag")
	bumpCmd.Flags().String("tag-prefix", "", "version prefix")
	bumpCmd.Flags().String("repository-location", ".", "repository location")
	bumpCmd.Flags().String("public-key-location", homedir+".ssh/", "public key location")
	rootCmd.AddCommand(bumpCmd)
}

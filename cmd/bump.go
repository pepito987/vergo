package cmd

import (
	"github.com/go-git/go-git/v5"
	"github.com/inanme/vergo/bump"
	vergo "github.com/inanme/vergo/git"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"golang.org/x/crypto/ssh"
	"io/ioutil"
)

var bumpCmd = &cobra.Command{
	Use:       "bump (patch|minor|major)",
	Short:     "bumps the version numbers",
	Args:      cobra.ExactValidArgs(1),
	ValidArgs: []string{"patch", "minor", "major"},
	RunE: func(cmd *cobra.Command, args []string) error {
		logLevelParam, err := cmd.Flags().GetString("log-level")
		CheckIfError(err)
		logLevel, err := log.ParseLevel(logLevelParam)
		if err != nil {
			log.WithError(err).Errorln("invalid log level, using INFO instead")
			log.SetLevel(log.InfoLevel)
		} else {
			log.SetLevel(logLevel)
		}

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
			keyLocation, err := cmd.Flags().GetString("key-location")
			CheckIfError(err)
			sshKey, err := ioutil.ReadFile(keyLocation)
			CheckIfError(err)
			passphrase, err := cmd.Flags().GetString("passphrase")
			CheckIfError(err)

			signer, err := ssh.ParsePrivateKeyWithPassphrase(sshKey, []byte(passphrase))
			CheckIfError(err)

			log.Debugf("Public key location %v", keyLocation)
			err = vergo.PushTag(repo, signer, version, prefix)
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
	bumpCmd.Flags().String("key-location", homedir+"/.ssh/id_rsa", `private key location, default: $homedir+"/.ssh/id_rsa"`)
	bumpCmd.Flags().String("passphrase", "", `private key passphrase`)
	rootCmd.AddCommand(bumpCmd)
}

package cmd

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	vergo "github.com/inanme/vergo/git"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "gets the version",
	Args:  cobra.NoArgs,
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
		version, err := vergo.LatestRef(repo, prefix)
		CheckIfError(err)
		fmt.Print(version.Version.String())
		return nil
	},
}

func init() {
	getCmd.Flags().String("tag-prefix", "", "version prefix")
	getCmd.Flags().String("repository-location", ".", "repository location")
	getCmd.Flags().String("public-key-location", ".ssh/", "public key location")
	rootCmd.AddCommand(getCmd)
}

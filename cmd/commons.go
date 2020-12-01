package cmd

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func sanitiseTagPrefix(tagPrefix string) string {
	switch strings.TrimSpace(tagPrefix) {
	case "":
		return "v"
	case "v":
		return "v"
	default:
		return tagPrefix + "-"
	}
}

func setLogger(cmd *cobra.Command) {
	logLevelParam, err := cmd.Flags().GetString("log-level")
	CheckIfError(err)
	logLevel, err := log.ParseLevel(logLevelParam)
	if err != nil {
		log.WithError(err).Errorln("invalid log level, using INFO instead")
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(logLevel)
	}
}

func checkAuthSocket(cmd *cobra.Command) (string, error) {
	pushTag, err := cmd.Flags().GetBool("push-tag")
	CheckIfError(err)
	if socket, found := os.LookupEnv("SSH_AUTH_SOCK"); pushTag && !found {
		return "", errors.New("SSH_AUTH_SOCK is not defined")
	} else {
		return socket, nil
	}
}

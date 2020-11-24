package cmd

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
)

const (
	appName    = "ygrph"
	appVersion = "0.1"

	logErrorCouldNotExecuteRootCommand = "could not execute root command"
)

var (
	rootCmd = &cobra.Command{
		Use:   appName,
		Short: appName + ": generate graphs from YAML definition files",
		Long:  "Define data in YAML then generate graph representations to model relationships",
	}

	username, password, dbURL string
	logLevel                  int8
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&dbURL, "dbURL", "d", "bolt://localhost:7687", "URL of graph database")
	rootCmd.PersistentFlags().StringVarP(&username, "username", "u", "username", "username for graph database")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "password", "password for graph database")
	rootCmd.PersistentFlags().Int8VarP(&logLevel, "logLevel", "l", int8(zerolog.WarnLevel), "log level")
}

// Execute TODO
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error().Err(err).Msg(logErrorCouldNotExecuteRootCommand)
		os.Exit(1)
	}
}

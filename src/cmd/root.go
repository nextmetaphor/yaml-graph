package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use:   commandRootUse,
		Short: commandRootUseShort,
		Long:  commandRootUseLong,
	}

	fileExtension string
	username      string
	password      string
	dbURL         string
	logLevel      int8
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&fileExtension, flagFileExtension, flagFileExtensionShorthand, flagFileExtensionDefault, flagFileExtensionUsage)
	rootCmd.PersistentFlags().StringVarP(&dbURL, flagDBURLName, flagDBURLShorthand, flagDBURLDefault, flagDBURLUsage)
	rootCmd.PersistentFlags().StringVarP(&username, flagUsernameName, flagUsernameShorthand, flagUsernameDefault, flagUsernameUsage)
	rootCmd.PersistentFlags().StringVarP(&password, flagPasswordName, flagPasswordShorthand, flagPasswordDefault, flagPasswordDefault)
	rootCmd.PersistentFlags().Int8VarP(&logLevel, flagLogLevelName, flagLogLevelShorthand, flagLogLevelDefault, flagLogLevelUsage)
}

// Execute TODO
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(exitCodeRootCmdFailed)
	}
}

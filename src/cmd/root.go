/*
 * Copyright 2020 Paul Tatham <paul@nextmetaphor.io>
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

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

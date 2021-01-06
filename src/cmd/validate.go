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
	"fmt"
	"github.com/nextmetaphor/yaml-graph/parser"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"os"
)

const (
	outputValidationSuccess = "successfully validated definitions"
	outputValidationFailure = "failed to validate definitions"
)

var (
	validateCmd = &cobra.Command{
		Use:   commandValidateUse,
		Short: commandValidateUseShort,
		Run:   validate,
	}

	validateSourceDir string
)

func init() {
	rootCmd.AddCommand(validateCmd)

	validateCmd.PersistentFlags().StringVarP(&validateSourceDir, flagSourceName, flagSourceShorthand, flagSourceDefault, flagSourceUsage)
}

func validate(_ *cobra.Command, _ []string) {
	zerolog.SetGlobalLevel(zerolog.Level(logLevel))

	d := parser.LoadDictionary(validateSourceDir, fileExtension)
	if parser.ValidateDictionary(d) != nil {
		fmt.Println(outputValidationFailure)
		os.Exit(exitCodeValidateCmdFailed)
	} else {
		fmt.Println(outputValidationSuccess)
	}
}

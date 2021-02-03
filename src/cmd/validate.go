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
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

const (
	outputValidationSuccess = "successfully validated definitions"
	outputValidationFailure = "failed to validate definitions"

	logErrorCouldNotOpenDefinitionFormatConfiguration             = "could not open definition format configuration [%s]"
	logErrorCouldNotUnmarshalDefinitionFormatConfiguration        = "could not unmarshal definition format configuration [%s]"
	logDebugSuccessfullyUnmarshalledDefinitionFormatConfiguration = "successfully unmarshalled definition format configuration [%s]"
)

var (
	validateCmd = &cobra.Command{
		Use:   commandValidateUse,
		Short: commandValidateUseShort,
		Run:   validate,
	}
)

func init() {
	rootCmd.AddCommand(validateCmd)

	validateCmd.PersistentFlags().StringVarP(&validateSourceDir, flagSourceName, flagSourceShorthand, flagSourceDefault,
		flagSourceUsage)
	validateCmd.PersistentFlags().StringVarP(&definitionFormatFile, flagDefinitionFormatName, flagDefinitionFormatShorthand,
		"", flagDefinitionFormatUsage)

}

func loadDefinitionFormatConf(cfgPath string) (definitionFormat *parser.DefinitionFormat, err error) {
	yamlFile, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		log.Error().Err(err).Msgf(logErrorCouldNotOpenDefinitionFormatConfiguration, cfgPath)
		return nil, err
	}
	err = yaml.UnmarshalStrict(yamlFile, &definitionFormat)
	if err != nil {
		log.Error().Err(err).Msgf(logErrorCouldNotUnmarshalDefinitionFormatConfiguration, cfgPath)
		return nil, err
	}

	log.Debug().Msgf(logDebugSuccessfullyUnmarshalledDefinitionFormatConfiguration, cfgPath)
	return definitionFormat, nil
}

func validate(_ *cobra.Command, _ []string) {
	zerolog.SetGlobalLevel(zerolog.Level(logLevel))

	var definitionFormat *parser.DefinitionFormat
	if definitionFormatFile != "" {
		var err error
		definitionFormat, err = loadDefinitionFormatConf(definitionFormatFile)
		if err != nil {
			fmt.Println(outputValidationFailure)
			os.Exit(exitCodeValidateCmdFailed)
		}
	}

	d := parser.LoadDictionary(validateSourceDir, fileExtension)
	if parser.ValidateDictionary(d, definitionFormat) != nil {
		fmt.Println(outputValidationFailure)
		os.Exit(exitCodeValidateCmdFailed)
	} else {
		fmt.Println(outputValidationSuccess)
	}
}

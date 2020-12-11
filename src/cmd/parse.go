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
	"github.com/nextmetaphor/yaml-graph/definition"
	"github.com/nextmetaphor/yaml-graph/graph"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
)

const (
	logDebugAboutToParseFile              = "about to parse file [%s]"
	logDebugSuccessfullyParsedFile        = "successfully parsed file [%s]"
	logWarnSkippingFile                   = "skipping file [%s] due to error [%s]"
	logErrorGraphDatabaseConnectionFailed = "graph database connection failed"
)

var (
	parseCmd = &cobra.Command{
		Use:   commandParseUse,
		Short: commandParseUseShort,
		Run:   parse,
	}

	parseSourceDir string
)

func init() {
	rootCmd.AddCommand(parseCmd)

	parseCmd.PersistentFlags().StringVarP(&parseSourceDir, flagSourceName, flagSourceShorthand, "", flagSourceUsage)
	parseCmd.MarkPersistentFlagRequired(flagSourceName)
}

func parse(_ *cobra.Command, _ []string) {
	zerolog.SetGlobalLevel(zerolog.Level(logLevel))

	driver, session, err := graph.Init(dbURL, username, password)
	if err != nil {
		log.Error().Err(err).Msg(logErrorGraphDatabaseConnectionFailed)
		os.Exit(exitCodeParseCmdFailed)
	}

	defer driver.Close()
	defer session.Close()

	graph.DeleteAll(session)

	// First create the nodes...
	definition.ProcessFiles(parseSourceDir, fileExtension, func(filePath string, _ os.FileInfo) (err error) {
		log.Debug().Msg(fmt.Sprintf(logDebugAboutToParseFile, filePath))

		spec, err := definition.LoadSpecificationFromFile(filePath)
		if (err == nil) && (spec != nil) {
			log.Debug().Msg(fmt.Sprintf(logDebugSuccessfullyParsedFile, filePath))
			graph.CreateSpecification(session, *spec)

		} else {
			log.Warn().Msgf(logWarnSkippingFile, filePath, err)
		}

		return nil
	})

	// ...then create the edges
	definition.ProcessFiles(parseSourceDir, fileExtension, func(filePath string, _ os.FileInfo) (err error) {
		log.Debug().Msg(fmt.Sprintf(logDebugAboutToParseFile, filePath))

		spec, err := definition.LoadSpecificationFromFile(filePath)
		if (err == nil) && (spec != nil) {
			log.Debug().Msg(fmt.Sprintf(logDebugSuccessfullyParsedFile, filePath))
			graph.CreateSpecificationEdge(session, *spec, nil)

		} else {
			log.Warn().Msgf(logWarnSkippingFile, filePath, err)
		}

		return nil
	})
}

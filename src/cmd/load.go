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
	logDebugAboutToLoadFile               = "about to load file [%s]"
	logDebugSuccessfullyLoadedFile        = "successfully loaded file [%s]"
	logWarnSkippingFile                   = "skipping file [%s] due to error [%s]"
	logErrorGraphDatabaseConnectionFailed = "graph database connection failed"
)

var (
	loadCmd = &cobra.Command{
		Use:   commandLoadUse,
		Short: commandLoadUseShort,
		Run:   load,
	}
)

func init() {
	rootCmd.AddCommand(loadCmd)

	loadCmd.PersistentFlags().StringSliceVarP(&sourceDir, flagSourceName, flagSourceShorthand, []string{flagSourceDefault}, flagSourceUsage)

}

func load(_ *cobra.Command, _ []string) {
	zerolog.SetGlobalLevel(zerolog.Level(logLevel))

	driver, session, err := graph.Init(dbURL, username, password)
	if err != nil {
		log.Error().Err(err).Msg(logErrorGraphDatabaseConnectionFailed)
		os.Exit(exitCodeLoadCmdFailed)
	}

	defer driver.Close()
	defer session.Close()

	graph.DeleteAll(session)

	// First create the nodes...
	for _, dir := range sourceDir {
		definition.ProcessFiles(dir, fileExtension, func(filePath string, _ os.FileInfo) (err error) {
			log.Debug().Msg(fmt.Sprintf(logDebugAboutToLoadFile, filePath))

			spec, err := definition.LoadSpecificationFromFile(filePath)
			if (err == nil) && (spec != nil) {
				log.Debug().Msg(fmt.Sprintf(logDebugSuccessfullyLoadedFile, filePath))
				graph.CreateSpecification(session, *spec)

			} else {
				log.Warn().Msgf(logWarnSkippingFile, filePath, err)
			}

			return nil
		})
	}

	// ...then create the edges
	for _, dir := range sourceDir {
		definition.ProcessFiles(dir, fileExtension, func(filePath string, _ os.FileInfo) (err error) {
			log.Debug().Msg(fmt.Sprintf(logDebugAboutToLoadFile, filePath))

			spec, err := definition.LoadSpecificationFromFile(filePath)
			if (err == nil) && (spec != nil) {
				log.Debug().Msg(fmt.Sprintf(logDebugSuccessfullyLoadedFile, filePath))
				graph.CreateSpecificationEdge(session, *spec, nil)

			} else {
				log.Warn().Msgf(logWarnSkippingFile, filePath, err)
			}

			return nil
		})
	}
}

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
	"os"
	"time"

	"github.com/nextmetaphor/yaml-graph/definition"
	"github.com/nextmetaphor/yaml-graph/graph"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
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

	loadCmd.Flags().StringSliceVarP(&sourceDir, flagSourceName, flagSourceShorthand, []string{flagSourceDefault}, flagSourceUsage)
	// default value provided so no need to mark flag as required
}

func load(_ *cobra.Command, _ []string) {
	zerolog.SetGlobalLevel(zerolog.Level(logLevel))

	start := time.Now()

	driver, session, err := graph.Init(dbURL, username, password)
	if err != nil {
		log.Error().Err(err).Msg(logErrorGraphDatabaseConnectionFailed)
		os.Exit(exitCodeLoadCmdFailed)
	}

	defer driver.Close()
	defer session.Close()

	graph.DeleteAll(session)
	batchConfig := graph.NewBatchConfig()

	fileCount := 0
	for _, dir := range sourceDir {
		definition.ProcessFiles(dir, fileExtension, func(filePath string, _ os.FileInfo) (err error) {
			fileCount++
			return nil
		})
	}

	fmt.Println("Loading definitions...")
	fileBar := progressbar.Default(int64(fileCount), "loading")

	for _, dir := range sourceDir {
		definition.ProcessFiles(dir, fileExtension, func(filePath string, _ os.FileInfo) (err error) {
			fileBar.Add(1)
			log.Debug().Msg(fmt.Sprintf(logDebugAboutToLoadFile, filePath))

			spec, err := definition.LoadSpecificationFromFile(filePath)
			if (err == nil) && (spec != nil) {
				log.Debug().Msg(fmt.Sprintf(logDebugSuccessfullyLoadedFile, filePath))
				batchConfig.AddSpecification(*spec, nil)

			} else {
				log.Warn().Msgf(logWarnSkippingFile, filePath, err)
			}

			return nil
		})
	}
	fileBar.Finish()
	fmt.Println()

	nodeClassCount := len(batchConfig.Nodes)
	fmt.Println("Creating definitions...")
	nodeBar := progressbar.Default(int64(nodeClassCount), "creating nodes")
	batchConfig.CreateNodes(session, func() { nodeBar.Add(1) })
	nodeBar.Finish()
	fmt.Println()

	edgeCount := 0
	for class := range batchConfig.Edges {
		edgeCount += len(batchConfig.Edges[class])
	}

	fmt.Println("Creating references")
	// We'll call progress for each edge in CreateEdges
	edgeBar := progressbar.Default(int64(edgeCount), "creating edges")
	batchConfig.CreateEdges(session, func() { edgeBar.Add(1) })
	edgeBar.Finish()
	fmt.Println()

	duration := time.Since(start)
	nodeCount := 0
	for class := range batchConfig.Nodes {
		nodeCount += len(batchConfig.Nodes[class])
	}

	fmt.Printf("Summary:\n")
	fmt.Printf("- Definition files loaded: %d\n", fileCount)
	fmt.Printf("- Definitions created:    %d\n", nodeCount)
	fmt.Printf("- References created:     %d\n", edgeCount)
	fmt.Printf("- Total time taken:       %v\n", duration)
}

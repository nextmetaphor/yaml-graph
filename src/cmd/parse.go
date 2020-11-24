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
		Use:   "parse",
		Short: "Parse definition files into graph",
		Run:   parse,
	}

	sourceDir string
)

func init() {
	rootCmd.AddCommand(parseCmd)
	parseCmd.Flags().StringVarP(&sourceDir, "source", "s", "definition/_test/CloudTaxonomy", "Source directory to read from")
}

func parse(cmd *cobra.Command, args []string) {
	zerolog.SetGlobalLevel(zerolog.Level(logLevel))

	driver, session, err := graph.Init(dbURL, username, password)
	if err != nil {
		log.Error().Err(err).Msg(logErrorGraphDatabaseConnectionFailed)
		os.Exit(2)
	}

	defer driver.Close()
	defer session.Close()

	graph.DeleteAll(session)

	// First create the nodes...
	definition.ProcessFiles(sourceDir, func(filePath string, _ os.FileInfo) (err error) {
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
	definition.ProcessFiles(sourceDir, func(filePath string, _ os.FileInfo) (err error) {
		log.Debug().Msg(fmt.Sprintf(logDebugAboutToParseFile, filePath))

		spec, err := definition.LoadSpecificationFromFile(filePath)
		if (err == nil) && (spec != nil) {
			log.Debug().Msg(fmt.Sprintf(logDebugSuccessfullyParsedFile, filePath))
			graph.CreateSpecificationEdge(session, *spec)

		} else {
			log.Warn().Msgf(logWarnSkippingFile, filePath, err)
		}

		return nil
	})

}

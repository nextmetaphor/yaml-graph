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
	"os"
)

const (
	outputTemplateFailure = "failed to generate template"
	logErrorReportFailed  = "report failed"
)

var (
	reportCmd = &cobra.Command{
		Use:   commandReportUse,
		Short: commandReportUseShort,
		Run:   doReport,
	}
)

func init() {
	rootCmd.AddCommand(reportCmd)

	reportCmd.Flags().StringVarP(&templateName, flagReportTemplateFileName, flagReportTemplateFileShorthand,
		"", flagReportTemplateFileUsage)
	if err := reportCmd.MarkFlagRequired(flagReportTemplateFileName); err != nil {
		log.Error().Err(err).Msg(logErrorReportFailed)
		os.Exit(exitCodeTemplateCmdFailed)
	}

	reportCmd.Flags().StringVarP(&templateFormat, flagReportFieldsFileName, flagReportFieldsFileShorthand,
		"", flagReportFieldsFileUsage)
	if err := reportCmd.MarkFlagRequired(flagReportFieldsFileName); err != nil {
		log.Error().Err(err).Msg(logErrorReportFailed)
		os.Exit(exitCodeTemplateCmdFailed)
	}

	reportCmd.Flags().BoolVarP(&loadDefinitions, flagLoadDefinitionsName, "", false, flagLoadDefinitionsUsage)

	reportCmd.Flags().StringSliceVarP(&sourceDir, flagSourceName, flagSourceShorthand, []string{flagSourceDefault}, flagSourceUsage)
	// default value provided so no need to mark flag as required
}

func doReport(c *cobra.Command, s []string) {
	zerolog.SetGlobalLevel(zerolog.Level(logLevel))

	if loadDefinitions {
		// TODO this is horrible - refactor
		load(c, s)
	}

	if err := parser.ParseTemplate(dbURL, username, password, templateFormat, templateName, os.Stdout); err != nil {
		fmt.Println(outputTemplateFailure)
		os.Exit(exitCodeTemplateCmdFailed)
	}
}

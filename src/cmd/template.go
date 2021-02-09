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
	outputTemplateFailure = "failed to generate template"
)

var (
	templateCmd = &cobra.Command{
		Use:   commandTemplateUse,
		Short: commandTemplateUseShort,
		Run:   doTemplate,
	}
)

func init() {
	rootCmd.AddCommand(templateCmd)

	templateCmd.PersistentFlags().StringVarP(&templateName, flagTemplateName, flagTemplateShorthand,
		"", flagTemplateUsage)
	templateCmd.MarkPersistentFlagRequired(flagTemplateName)

	templateCmd.PersistentFlags().StringVarP(&templateFormat, flagDefinitionFormatName, flagDefinitionFormatShorthand,
		"", flagDefinitionFormatUsage)
	templateCmd.MarkPersistentFlagRequired(flagDefinitionFormatName)

	templateCmd.PersistentFlags().BoolVarP(&loadDefinitions, flagLoadDefinitionsName, "", false, flagLoadDefinitionsUsage)

	templateCmd.PersistentFlags().StringVarP(&loadSourceDir, flagSourceName, flagSourceShorthand, flagSourceDefault, flagSourceUsage)

}

func doTemplate(c *cobra.Command, s []string) {
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

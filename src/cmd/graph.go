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
)

const (
	nodeString = "{\"id\": \"%s-%s\",\"class\": \"%s\",\"description\": \"%s\"},"
	linkString = "{\"source\": \"%s-%s\",\"target\": \"%s-%s\"},"
)

var (
	graphCmd = &cobra.Command{
		Use:   commandGraphUse,
		Short: commandGraphUseShort,
		Run:   graphFunc,
	}

	graphSourceDir string
)

func init() {
	rootCmd.AddCommand(graphCmd)

	graphCmd.PersistentFlags().StringVarP(&graphSourceDir, flagSourceName, flagSourceShorthand, "", flagSourceUsage)
	graphCmd.MarkPersistentFlagRequired(flagSourceName)
}

func graphFunc(cmd *cobra.Command, args []string) {
	zerolog.SetGlobalLevel(zerolog.Level(logLevel))

	d := parser.LoadDictionary(graphSourceDir, fileExtension)

	fmt.Print("{\"nodes\": [")
	for class := range d {
		for id := range d[class] {
			fmt.Print(fmt.Sprintf(nodeString, class, id, class, d[class][id].Fields["Name"]))
		}
	}
	fmt.Print("],\"links\": [")
	for class := range d {
		for id := range d[class] {
			for _, ref := range d[class][id].References {
				fmt.Print(fmt.Sprintf(linkString, class, id, ref.Class, ref.ID))
			}
		}
	}
	fmt.Print("]}")
}

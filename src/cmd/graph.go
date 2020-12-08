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
	nodeHeaderString      = "{\"nodes\": ["
	perNodeStringFirst    = "{\"id\": \"%s-%s\",\"class\": \"%s\",\"description\": \"%s\"}"
	perNodeStringNotFirst = ",{\"id\": \"%s-%s\",\"class\": \"%s\",\"description\": \"%s\"}"
	linkHeaderString      = "],\"links\": ["
	perLinkStringFirst    = "{\"source\": \"%s-%s\",\"target\": \"%s-%s\", \"relationship\": \"%s\"}"
	perLinkStringNotFirst = ",{\"source\": \"%s-%s\",\"target\": \"%s-%s\", \"relationship\": \"%s\"}"
	linkFooterString      = "]}"
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

	fmt.Print(nodeHeaderString)
	firstElement := true
	for class, definitions := range d {
		for id, definition := range definitions {
			if firstElement {
				fmt.Print(fmt.Sprintf(perNodeStringFirst, class, id, class, definition.Fields["Name"]))
			} else {
				fmt.Print(fmt.Sprintf(perNodeStringNotFirst, class, id, class, definition.Fields["Name"]))
			}
			firstElement = false
		}
	}
	fmt.Print(linkHeaderString)
	firstElement = true
	for class, definitions := range d {
		for id, definition := range definitions {
			for _, ref := range definition.References {
				if firstElement {
					fmt.Print(fmt.Sprintf(perLinkStringFirst, class, id, ref.Class, ref.ID, ref.Relationship))
				} else {
					fmt.Print(fmt.Sprintf(perLinkStringNotFirst, class, id, ref.Class, ref.ID, ref.Relationship))
				}
				firstElement = false
			}
		}
	}
	fmt.Print(linkFooterString)
}

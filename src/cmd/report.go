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
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/nextmetaphor/yaml-graph/graph"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

const (
	rootCypher  = "match (n:%s) return n"
	childCypher = "match (n:%s)-[:%s]-(p:%s {ID:\"%s\"}) return n"

	markdownSection     = "%s%s%s"       //prefix section suffix
	markdownDetailField = "%s%s%s%s%s%s" //field key (prefix value suffix) + field value (prefix value suffix)

	logErrorExecutingCypher                             = "error executing cypher"
	logErrorCouldNotOpenReportConfiguration             = "could not open report configuration [%s]"
	logErrorCouldNotUnmarshalReportConfiguration        = "could not unmarshal report configuration [%s]"
	logDebugSuccessfullyUnmarshalledReportConfiguration = "successfully unmarshalled report configuration [%s]"
)

type (
	// MarkdownSection TODO
	MarkdownSection struct {
		// Class indicates the class of defintion to use for the section
		Class string `yaml:"Class"`
		// SectionNameField indicates the field to use for the section name
		SectionNameField string `yaml:"SectionNameField"`
		// SectionNamePrefix is any prefix before the section name
		SectionNamePrefix string `yaml:"SectionNamePrefix"`
		// SectionNameSuffix  is any suffix after the section name
		SectionNameSuffix string `yaml:"SectionNameSuffix"`
		SectionHeader     string `yaml:"SectionHeader"`
		// DetailFields indicates which fields which will be shown in the section
		DetailFields           []string `yaml:"DetailFields"`
		DetailFieldKeyPrefix   string   `yaml:"DetailFieldKeyPrefix"`
		DetailFieldKeySuffix   string   `yaml:"DetailFieldKeySuffix"`
		DetailFieldValuePrefix string   `yaml:"DetailFieldValuePrefix"`
		DetailFieldValueSuffix string   `yaml:"DetailFieldValueSuffix"`

		ParentRelationship string            `yaml:"ParentRelationship"`
		ChildSection       []MarkdownSection `yaml:"ChildSection"`

		// SectionHeader is any footer text after the section details
		SectionFooter string `yaml:"SectionFooter"`
	}
)

var (
	reportCmd = &cobra.Command{
		Use:   commandReportUse,
		Short: commandReportUseShort,
		Run:   report,
	}

	reportDefinition string
)

func init() {
	rootCmd.AddCommand(reportCmd)

	rootCmd.PersistentFlags().StringVarP(&reportDefinition, flagReportDefinitionName, flagReportDefinitionShorthand,
		"", flagReportDefinitionUsage)
	reportCmd.MarkFlagRequired(flagReportDefinitionName)
}

func loadReportConf() (ms *MarkdownSection, err error) {
	yamlFile, err := ioutil.ReadFile(reportDefinition)
	if err != nil {
		log.Error().Err(err).Msgf(logErrorCouldNotOpenReportConfiguration, reportDefinition)
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &ms)
	if err != nil {
		log.Error().Err(err).Msgf(logErrorCouldNotUnmarshalReportConfiguration, reportDefinition)
	}

	log.Debug().Msgf(logDebugSuccessfullyUnmarshalledReportConfiguration, reportDefinition)
	return ms, nil
}

func report(cmd *cobra.Command, args []string) {
	zerolog.SetGlobalLevel(zerolog.Level(logLevel))

	// first load the report configuration
	markdownSection, err := loadReportConf()
	if err != nil {
		log.Error().Err(err).Msg(logErrorGraphDatabaseConnectionFailed)
		os.Exit(exitCodeReportCmdFailed)
	}

	// then connect to the graph database
	driver, session, err := graph.Init(dbURL, username, password)
	if err != nil {
		log.Error().Err(err).Msg(logErrorGraphDatabaseConnectionFailed)
		os.Exit(exitCodeReportCmdFailed)
	}

	defer driver.Close()
	defer session.Close()

	// now recurse through the sections
	recurseSection(session, *markdownSection, nil, nil)
}

func recurseSection(session neo4j.Session, section MarkdownSection, parentClass, parentID *string) error {
	var res neo4j.Result
	var err error

	if (parentClass == nil) || (parentID == nil) {
		res, err = graph.ExecuteCypher(session, fmt.Sprintf(rootCypher, section.Class), nil)
	} else {
		res, err = graph.ExecuteCypher(session, fmt.Sprintf(childCypher, section.Class, section.ParentRelationship, *parentClass,
			*parentID), nil)
	}

	if (err != nil) || (res.Err() != nil) {
		log.Error().Err(err).Msgf(logErrorExecutingCypher)
		return err
	}

	for res.Next() {
		record := res.Record()
		for _, kv := range record.Values() {
			node, isNode := kv.(neo4j.Node)
			if isNode {
				// write the section name plus prefixes & suffixes
				fmt.Print(fmt.Sprintf(markdownSection, section.SectionNamePrefix, node.Props()[section.SectionNameField],
					section.SectionNameSuffix))

				// write any section header
				fmt.Print(section.SectionHeader)

				// write the node details
				for _, key := range section.DetailFields {
					keyValue, keyOK := node.Props()[key].(string)
					if !keyOK {
						keyValue = "-"
					}

					fmt.Print(fmt.Sprintf(markdownDetailField,
						section.DetailFieldKeyPrefix, key, section.DetailFieldKeySuffix,
						section.DetailFieldValuePrefix, keyValue, section.DetailFieldValueSuffix))
				}

				// recurse through any child sections
				// TODO recursion, really?
				for _, childSection := range section.ChildSection {
					nodeID := node.Props()["ID"].(string)
					recurseSection(session, childSection, &(section.Class), &nodeID)
				}

				// write any section footer
				fmt.Print(section.SectionFooter)
			}
		}
	}

	return err
}

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
	"html/template"
	"io/ioutil"
	"os"
)

const (
	rootTemplateCypher  = "match (n:%s) return n order by n.%s"
	childTemplateCypher = "match (n:%s)-[:%s]-(p:%s {ID:\"%s\"}) return n order by n.%s"

	logErrorExecutingCypher                               = "error executing cypher"
	logErrorCouldNotOpenTemplateConfiguration             = "could not open template configuration [%s]"
	logErrorCouldNotUnmarshalTemplateConfiguration        = "could not unmarshal template configuration [%s]"
	logErrorParsingTemplateDefinitions                    = "error parsing template definitions"
	logErrorParsingTemplate                               = "error parsing template"
	logDebugSuccessfullyUnmarshalledTemplateConfiguration = "successfully unmarshalled template configuration [%s]"
)

type (
	// TemplateSection TODO
	TemplateSection struct {
		// Class indicates the class of definition to use for this section
		Class string `yaml:"Class"`
		// Fields indicates which fields which will be retrieved in the section
		Fields     []string `yaml:"Fields"`
		OrderField string   `yaml:"OrderField"`

		ParentRelationship string            `yaml:"ParentRelationship"`
		ChildSection       []TemplateSection `yaml:"ChildSection"`
	}

	templateDefinition struct {
		Class  string
		ID     string
		Fields map[string]string

		// ReferencedDefinitions is a map of definitions keyed by relationship
		ReferencedDefinitions map[string][]templateDefinition
	}
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

func loadTemplateConf(cfgPath string) (ms *TemplateSection, err error) {
	yamlFile, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		log.Error().Err(err).Msgf(logErrorCouldNotOpenTemplateConfiguration, cfgPath)
		return nil, err
	}
	err = yaml.UnmarshalStrict(yamlFile, &ms)
	if err != nil {
		log.Error().Err(err).Msgf(logErrorCouldNotUnmarshalTemplateConfiguration, cfgPath)
		return nil, err
	}

	log.Debug().Msgf(logDebugSuccessfullyUnmarshalledTemplateConfiguration, cfgPath)
	return ms, nil
}

func doTemplate(c *cobra.Command, s []string) {
	zerolog.SetGlobalLevel(zerolog.Level(logLevel))

	if loadDefinitions {
		// TODO this is horrible - refactor
		load(c, s)
	}

	// first load the template configuration
	templateSection, err := loadTemplateConf(templateFormat)
	if err != nil {
		log.Error().Err(err).Msg(logErrorGraphDatabaseConnectionFailed)
		os.Exit(exitCodeTemplateCmdFailed)
	}

	// then connect to the graph database
	driver, session, err := graph.Init(dbURL, username, password)
	if err != nil {
		log.Error().Err(err).Msg(logErrorGraphDatabaseConnectionFailed)
		os.Exit(exitCodeTemplateCmdFailed)
	}

	defer driver.Close()
	defer session.Close()

	// now recurse through the sections
	definitions, err := recurseTemplateSection(session, *templateSection, nil, nil)
	if err != nil {
		log.Error().Err(err).Msg(logErrorParsingTemplateDefinitions)
		os.Exit(exitCodeTemplateCmdFailed)
	}

	//fmt.Println(definitions)

	template := template.Must(template.ParseFiles(templateName))

	template.Execute(os.Stdout, definitions)
}

func recurseTemplateSection(session neo4j.Session, section TemplateSection, parentClass, parentID *string) ([]templateDefinition, error) {
	var res neo4j.Result
	var err error
	var definitions []templateDefinition

	if (parentClass == nil) || (parentID == nil) {
		res, err = graph.ExecuteCypher(session, fmt.Sprintf(rootTemplateCypher, section.Class, section.OrderField), nil)
	} else {
		res, err = graph.ExecuteCypher(session, fmt.Sprintf(childTemplateCypher, section.Class, section.ParentRelationship, *parentClass,
			*parentID, section.OrderField), nil)
	}

	if (err != nil) || (res.Err() != nil) {
		log.Error().Err(err).Msgf(logErrorExecutingCypher)
		return definitions, err
	}

	for res.Next() {
		record := res.Record()
		for _, kv := range record.Values() {
			node, isNode := kv.(neo4j.Node)
			if isNode {
				var definition templateDefinition = templateDefinition{
					Class:                 section.Class,
					Fields:                map[string]string{},
					ReferencedDefinitions: map[string][]templateDefinition{},
				}

				// TODO runtime checking needed
				definition.ID = node.Props()["ID"].(string)

				for _, key := range section.Fields {
					keyValue, keyOK := node.Props()[key].(string)
					if keyOK {
						definition.Fields[key] = keyValue
					}
				}

				// recurse through any child sections
				// TODO recursion, really?
				for _, childSection := range section.ChildSection {
					nodeID := node.Props()["ID"].(string)
					childDefinitions, err := recurseTemplateSection(session, childSection, &(section.Class), &nodeID)
					if err != nil {
						log.Err(err).Msg(logErrorParsingTemplateDefinitions)
						return definitions, err
					}
					definition.ReferencedDefinitions[childSection.ParentRelationship] = childDefinitions
				}

				definitions = append(definitions, definition)
			}
		}
	}

	return definitions, err
}

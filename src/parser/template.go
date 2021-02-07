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

package parser

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/nextmetaphor/yaml-graph/graph"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
	"html/template"
	"io/ioutil"
	"os"
	"strings"
)

const (
	orderClauseSingular = "%s.%s"
	orderClauseMultiple = "%s,%s.%s"
	rootTemplateCypher  = "match (%s:%s) return %s order by %s"
	childTemplateCypher = "match (%s:%s)-[:%s]-(%s:%s {ID:\"%s\"}) return %s order by %s"

	logErrorExecutingCypher                               = "error executing cypher"
	logErrorCouldNotOpenTemplateConfiguration             = "could not open template configuration [%s]"
	logErrorCouldNotUnmarshalTemplateConfiguration        = "could not unmarshal template configuration [%s]"
	logErrorParsingTemplateDefinitions                    = "error parsing template definitions"
	logErrorParsingTemplate                               = "error parsing template"
	logDebugSuccessfullyUnmarshalledTemplateConfiguration = "successfully unmarshalled template configuration [%s]"
	logErrorGraphDatabaseConnectionFailed                 = "graph database connection failed"
	logErrorNilDefinitionID                               = "definition ID is nil"
	logErrorNonStringDefinitionID                         = "definition ID is not a string"
)

type (
	// ClassFieldSelector TODO
	ClassFieldSelector struct {
		// Class indicates the class of definition to select
		Class string `yaml:"Class"`
		// Fields indicates which fields should be retrieved
		Fields []string `yaml:"Fields"`

		Relationship string `yaml:"Relationship,omitempty"`

		// OrderField indicates the fields to order the classes retrieved
		OrderFields []string `yaml:"OrderFields"`
	}

	// TemplateSection TODO
	TemplateSection struct {
		// SectionClass identifies the class for the section
		SectionClass ClassFieldSelector `yaml:"SectionClass"`

		// AggregateClasses identifies any aggregate classes referenced by the main class
		AggregateClasses []ClassFieldSelector `yaml:"AggregateClasses,omitempty"`

		// CompositeSections defines any sections for composite classes referenced by the section class
		CompositeSections []TemplateSection `yaml:"CompositeSections,omitempty"`
	}

	// ClassFieldIdentifier TODO
	ClassFieldIdentifier struct {
		Class string
		Field string
	}

	// SectionDefinition TODO
	SectionDefinition struct {
		Class  string
		ID     string
		Fields map[ClassFieldIdentifier]string

		// CompositeSectionDefinitions is a map of definitions keyed by relationship
		CompositeSectionDefinitions map[string][]SectionDefinition
	}
)

func getOrderClause(selector ClassFieldSelector) (orderClause string) {
	for i, field := range selector.OrderFields {
		if i == 0 {
			orderClause = fmt.Sprintf(orderClauseSingular, selector.Class, field)
		} else {
			orderClause = fmt.Sprintf(orderClauseMultiple, orderClause, selector.Class, field)
		}
	}

	return
}

func getCypherForSelector(parentClass string, parentID string, selector ClassFieldSelector) string {
	parentClass = strings.TrimSpace(parentClass)

	if parentClass == "" {
		return fmt.Sprintf(rootTemplateCypher, selector.Class, selector.Class, selector.Class,
			getOrderClause(selector))
	}
	return fmt.Sprintf(childTemplateCypher, selector.Class, selector.Class, selector.Relationship,
		parentClass, parentClass, strings.TrimSpace(parentID), selector.Class, getOrderClause(selector))
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

func parseTemplate(dbURL, username, password, templateConf, templatePath string) error {
	// first load the template configuration
	templateSection, err := loadTemplateConf(templateConf)
	if err != nil {
		log.Error().Err(err).Msg(logErrorGraphDatabaseConnectionFailed)
		return err
	}

	// then connect to the graph database
	driver, session, err := graph.Init(dbURL, username, password)
	if err != nil {
		log.Error().Err(err).Msg(logErrorGraphDatabaseConnectionFailed)
		return err
	}

	defer driver.Close()
	defer session.Close()

	// now recurse through the sections
	definitions, err := recurseTemplateSection(session, *templateSection, nil, nil)
	if err != nil {
		log.Error().Err(err).Msg(logErrorParsingTemplateDefinitions)
		return err
	}

	template := template.Must(template.ParseFiles(templatePath))
	return template.Execute(os.Stdout, definitions)
}

func recurseTemplateSection(session neo4j.Session, section TemplateSection, parentClass, parentID *string) ([]SectionDefinition, error) {
	var res neo4j.Result
	var err error
	var definitions []SectionDefinition
	var cypher string

	if parentClass == nil || parentID == nil {
		cypher = getCypherForSelector("", "", section.SectionClass)
	} else {
		cypher = getCypherForSelector(*parentClass, *parentID, section.SectionClass)
	}

	res, err = graph.ExecuteCypher(session,  cypher,nil)
	if (err != nil) || (res.Err() != nil) {
		log.Error().Err(err).Msgf(logErrorExecutingCypher)
		return definitions, err
	}

	for res.Next() {
		record := res.Record()
		for _, kv := range record.Values() {
			node, isNode := kv.(neo4j.Node)
			if isNode {
				definition := SectionDefinition{
					Class:                       section.SectionClass.Class,
					Fields:                      map[ClassFieldIdentifier]string{},
					CompositeSectionDefinitions: map[string][]SectionDefinition{},
				}

				if node.Props()["ID"] != nil {
					if definitionID, ok := node.Props()["ID"].(string); ok {
						definition.ID = definitionID
					} else {
						log.Warn().Msg(logErrorNonStringDefinitionID)
					}
				} else {
					log.Warn().Msg(logErrorNilDefinitionID)
				}

				for _, key := range section.SectionClass.Fields {
					keyValue, keyOK := node.Props()[key].(string)
					if keyOK {
						definition.Fields[ClassFieldIdentifier{
							Class: section.SectionClass.Class,
							Field: key,
						}] = keyValue
					}
				}

				// recurse through any child sections
				//// TODO recursion, really?
				//for _, childSection := range section.ChildSection {
				//	nodeID := node.Props()["ID"].(string)
				//	childDefinitions, err := recurseTemplateSection(session, childSection, &(section.Class), &nodeID)
				//	if err != nil {
				//		log.Err(err).Msg(logErrorParsingTemplateDefinitions)
				//		return definitions, err
				//	}
				//	definition.ReferencedDefinitions[childSection.ParentRelationship] = childDefinitions
				//}
				//

				definitions = append(definitions, definition)
			}
		}
	}

	return definitions, err
}

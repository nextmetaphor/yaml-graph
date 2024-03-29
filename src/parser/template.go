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
	"bytes"
	"fmt"
	htmltemplate "html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
	texttemplate "text/template"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/nextmetaphor/yaml-graph/graph"
	"github.com/rs/zerolog/log"
	"github.com/yuin/goldmark"
	"gopkg.in/yaml.v3"
)

const (
	funcMarkdown = "markdown"
	funcNilToStr = "nilToStr"

	orderClauseSingular        = "%s.%s"
	orderClauseMultiple        = "%s,%s.%s"
	baseTemplateCypher         = "match %s return %s order by %s"
	rootCypherMatchClause      = "(%s:%s)"
	compositeCypherMatchClause = "(%s:%s)%s-[:%s]-%s(%s:%s {ID:\"%s\"})"
	aggregateCypherMatchClause = " optional match (%s:%s)-[:%s]-(%s:%s)"
	aggregateCypherOrderClause = ",%s"

	classFieldIdentifier = "%s.%s"

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

		// ClassAlias indicates an optional alias for the class being selected
		ClassAlias string `yaml:"ClassAlias"`

		// Fields indicates which fields should be retrieved
		Fields []string `yaml:"Fields"`

		Relationship string `yaml:"Relationship,omitempty"`

		// RelationshipFrom indicates whether we want the relationship to be directed from this class (defaults to false)
		RelationshipFrom bool `yaml:"RelationshipFrom,omitempty"`

		// RelationshipFrom indicates whether we want the relationship to be directed to this class (defaults to false)
		RelationshipTo bool `yaml:"RelationshipTo,omitempty"`

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

	// SectionDefinition TODO
	SectionDefinition struct {
		Class string
		ID    string
		// Fields key must be in format Class.Field
		Fields map[string]interface{}

		// CompositeSectionDefinitions is a map of definitions keyed by relationship
		CompositeSectionDefinitions map[string][]SectionDefinition
	}
)

func getOrderClause(section TemplateSection) (orderClause string) {
	for i, field := range section.SectionClass.OrderFields {
		if i == 0 {
			orderClause = fmt.Sprintf(orderClauseSingular, section.SectionClass.Class, field)
		} else {
			orderClause = fmt.Sprintf(orderClauseMultiple, orderClause, section.SectionClass.Class, field)
		}
	}

	return
}

func getCypherForSection(parentClass string, parentID string, section TemplateSection) string {
	var matchClause, returnClause, orderClause string

	sectionClass := strings.TrimSpace(section.SectionClass.Class)
	sectionClassAlias := strings.TrimSpace(section.SectionClass.ClassAlias)

	// default to the class name if no alias is provided
	if sectionClassAlias == "" {
		sectionClassAlias = sectionClass
	}
	parentClass = strings.TrimSpace(parentClass)

	returnClause = sectionClassAlias
	orderClause = getOrderClause(section)

	relationshipFrom := ""
	if section.SectionClass.RelationshipFrom {
		relationshipFrom = "<"
	}

	relationshipTo := ""
	if section.SectionClass.RelationshipTo {
		relationshipTo = ">"
	}

	if parentClass == "" {
		matchClause = fmt.Sprintf(rootCypherMatchClause, sectionClass, sectionClass)
	} else {
		matchClause = fmt.Sprintf(compositeCypherMatchClause, sectionClassAlias, sectionClass,
			relationshipFrom, section.SectionClass.Relationship, relationshipTo, parentClass, parentClass, strings.TrimSpace(parentID))
	}

	for _, aggregateClass := range section.AggregateClasses {
		aggregateMatchClause := fmt.Sprintf(aggregateCypherMatchClause, sectionClass, sectionClass,
			aggregateClass.Relationship, aggregateClass.Class, aggregateClass.Class)
		matchClause = matchClause + aggregateMatchClause

		aggregateReturnClause := fmt.Sprintf(aggregateCypherOrderClause, aggregateClass.Class)
		returnClause = returnClause + aggregateReturnClause
	}

	return fmt.Sprintf(baseTemplateCypher, matchClause, returnClause, orderClause)
}

func loadTemplateConf(cfgPath string) (ms *TemplateSection, err error) {
	yamlFile, err := os.Open(cfgPath)
	if err != nil {
		log.Error().Err(err).Msgf(logErrorCouldNotOpenTemplateConfiguration, cfgPath)
		return nil, err
	}

	defer yamlFile.Close()
	d := yaml.NewDecoder(yamlFile)
	d.KnownFields(true)
	if err := d.Decode(&ms); err != nil {
		log.Error().Err(err).Msgf(logErrorCouldNotOpenTemplateConfiguration, cfgPath)
		return nil, err
	}

	log.Debug().Msgf(logDebugSuccessfullyUnmarshalledTemplateConfiguration, cfgPath)
	return ms, nil
}

func fromMarkdown(inputString string) htmltemplate.HTML {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(inputString), &buf); err != nil {
		panic(err)
	}

	return htmltemplate.HTML(buf.String())
}

func nilToStr(v interface{}) interface{} {
	if v == nil {
		return ""
	}
	return v
}

func getTemplateFuncs() texttemplate.FuncMap {
	return texttemplate.FuncMap{
		funcMarkdown: fromMarkdown,
		funcNilToStr: nilToStr,
	}
}

// ParseTemplate TODO
func ParseTemplate(dbURL, username, password, templateConf, templatePath string, writer io.Writer) error {
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

	template := texttemplate.Must(texttemplate.New(filepath.Base(templatePath)).Option("missingkey=zero").Funcs(getTemplateFuncs()).ParseFiles(templatePath))
	return template.Execute(writer, definitions)
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func recurseTemplateSection(session neo4j.Session, section TemplateSection, parentClass, parentID *string) ([]SectionDefinition, error) {
	var res neo4j.Result
	var err error
	var definitions []SectionDefinition
	var cypher string

	if parentClass == nil || parentID == nil {
		cypher = getCypherForSection("", "", section)
	} else {
		cypher = getCypherForSection(*parentClass, *parentID, section)
	}

	res, err = graph.ExecuteCypher(session, cypher, nil)
	if (err != nil) || (res.Err() != nil) {
		log.Error().Err(err).Msgf(logErrorExecutingCypher)
		return definitions, err
	}

	for res.Next() {
		record := res.Record()

		// create definition here: class+aggregate combinations are returned by the subsequent loop
		definition := SectionDefinition{
			Class:                       section.SectionClass.Class,
			Fields:                      map[string]interface{}{},
			CompositeSectionDefinitions: map[string][]SectionDefinition{},
		}
		for _, kv := range record.Values {
			node, isNode := kv.(neo4j.Node)
			if isNode {
				// TODO dangerous - refactor
				nodeClass := node.Labels[0]

				if nodeClass == section.SectionClass.Class {
					if node.Props["ID"] != nil {
						if definitionID, ok := node.Props["ID"].(string); ok {
							definition.ID = definitionID
						} else {
							log.Warn().Msg(logErrorNonStringDefinitionID)
						}
					} else {
						log.Warn().Msg(logErrorNilDefinitionID)
					}

					nodeClassAlias := strings.TrimSpace(section.SectionClass.ClassAlias)
					if nodeClassAlias == "" {
						nodeClassAlias = section.SectionClass.Class
					}

					for _, key := range section.SectionClass.Fields {
						if fieldTypeValid(node.Props[key]) {
							definition.Fields[fmt.Sprintf(classFieldIdentifier, nodeClassAlias, key)] = node.Props[key]
						}
					}
				} else {
					for _, a := range section.AggregateClasses {
						if nodeClass == a.Class {
							for _, key := range a.Fields {
								if fieldTypeValid(node.Props[key]) {
									// TODO need to use the alias for non-section class
									definition.Fields[fmt.Sprintf(classFieldIdentifier, nodeClass, key)] = node.Props[key]
								}
							}
						}
					}
				}

				//recurse through any child sections
				// TODO recursion, really?
				for _, childSection := range section.CompositeSections {

					childDefinitions, err := recurseTemplateSection(session, childSection, &(section.SectionClass.Class), &definition.ID)
					if err != nil {
						log.Err(err).Msg(logErrorParsingTemplateDefinitions)
						return definitions, err
					}
					definition.CompositeSectionDefinitions[childSection.SectionClass.Relationship] = childDefinitions
				}
			}
		}

		definitions = append(definitions, definition)
	}

	return definitions, err
}

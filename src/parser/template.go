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
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
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

	templateDefinition struct {
		Class  string
		ID     string
		Fields map[string]string

		// ReferencedDefinitions is a map of definitions keyed by relationship
		ReferencedDefinitions map[string][]templateDefinition
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

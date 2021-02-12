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
	"github.com/nextmetaphor/yaml-graph/definition"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
)

const (
	logDebugAboutToParseFile       = "about to parse file [%s]"
	logDebugSuccessfullyParsedFile = "successfully parsed file [%s]"
	logWarnSkippingFile            = "skipping file [%s] due to error [%s]"

	logWarnCannotFindClass          = "cannot find class [%s]"
	logWarnCannotFindDefinition     = "cannot find definition ID [%s] for class [%s]"
	logWarnMandatoryFieldMissing    = "mandatory field [%s] missing in definition ID [%s] for class [%s]"
	logWarnMandatoryFieldNotAString = "mandatory field [%s] is not a string in definition ID [%s] for class [%s]"
	logWarnAdditionalFieldFound     = "field [%s] is not a valid field in definition ID [%s] for class [%s]"

	errorDefinitionErrorsFound = "there were %d error(s) found in the definition files"
)

type (
	// DictionaryDefinition TODO
	DictionaryDefinition struct {
		Fields     definition.Fields
		References []definition.Reference
	}

	// Dictionary is a map of classes, keyed by class name; the value is a map of definitions keyed by
	// definition ID
	Dictionary map[string]map[string]*DictionaryDefinition
)

type (
	// DefinitionFormat TODO
	DefinitionFormat struct {
		// ClassFormat is a map of classes keyed by class name; the value is format of each class
		ClassFormat map[string]*ClassDefinitionFormat `yaml:"Class"`
	}

	// ClassField TODO
	ClassField struct {
		Description string `yaml:"Description,omitempty"`
	}

	// ClassDefinitionFormat TODO
	ClassDefinitionFormat struct {
		Description     string                `yaml:"Description,omitempty"`
		MandatoryFields map[string]ClassField `yaml:"MandatoryFields"`
		OptionalFields  map[string]ClassField `yaml:"OptionalFields"`
	}
)

func loadSpecification(s definition.Specification, d Dictionary, parentRef *definition.Reference) {
	if d == nil {
		return
	}

	if d[s.Class] == nil {
		// this is the first time we've encountered this class - create a map for it
		d[s.Class] = make(map[string]*DictionaryDefinition)
	}

	// iterate through the definitions in this specification and add to the dictionary
	for dfnID, dfn := range s.Definitions {
		d[s.Class][dfnID] = &DictionaryDefinition{
			Fields:     dfn.Fields,
			References: dfn.References,
		}
	}

	// add each reference from the specification to the individual definitions
	for _, ref := range s.References {
		for dfnID := range d[s.Class] {
			d[s.Class][dfnID].References = append(d[s.Class][dfnID].References, ref)
		}
	}

	// finally add any parent ref passed in
	if parentRef != nil {
		for dfnID := range s.Definitions {
			d[s.Class][dfnID].References = append(d[s.Class][dfnID].References, *parentRef)
		}
	}

	// now recurse through any subdefinitions, repeating the above process and passing the parent reference
	// TODO, recursion, really?
	for dfnID, dfn := range s.Definitions {
		for subDfnRelationship, subDfn := range dfn.SubDefinitions {
			loadSpecification(subDfn, d, &definition.Reference{
				Class:        s.Class,
				ID:           dfnID,
				Relationship: subDfnRelationship,
			})
		}
	}
}

// LoadDictionary TODO
func LoadDictionary(sourceDir, fileExtension string) Dictionary {
	d := make(Dictionary)
	definition.ProcessFiles(sourceDir, fileExtension, func(filePath string, _ os.FileInfo) (err error) {
		log.Debug().Msg(fmt.Sprintf(logDebugAboutToParseFile, filePath))

		spec, err := definition.LoadSpecificationFromFile(filePath)
		if (err == nil) && (spec != nil) {
			log.Debug().Msg(fmt.Sprintf(logDebugSuccessfullyParsedFile, filePath))
			loadSpecification(*spec, d, nil)

		} else {
			log.Warn().Msgf(logWarnSkippingFile, filePath, err)
		}

		return nil
	})

	return d
}

// ValidateDictionary TODO
func ValidateDictionary(d Dictionary, df *DefinitionFormat) error {
	errorsFound := 0

	if df != nil {
		for class, classFormat := range df.ClassFormat {
			for dID, definition := range d[class] {
				// first check that each field in the definition is either a mandatory or optional field...
				for defField := range definition.Fields {
					_, isMandatoryField := classFormat.MandatoryFields[defField]
					_, isOptionalField := classFormat.OptionalFields[defField]
					if !isMandatoryField && !isOptionalField {
						log.Warn().Msg(fmt.Sprintf(logWarnAdditionalFieldFound, defField, dID, class))
						errorsFound++
					}
				}

				// ...then validate each of the mandatory fields exists within the definition
				for f := range classFormat.MandatoryFields {
					if definition.Fields[f] == nil {
						log.Warn().Msg(fmt.Sprintf(logWarnMandatoryFieldMissing, f, dID, class))
						errorsFound++
					} else {
						s, ok := definition.Fields[f].(string)
						if ok {
							if strings.TrimSpace(s) == "" {
								log.Warn().Msg(fmt.Sprintf(logWarnMandatoryFieldMissing, f, dID, class))
								errorsFound++
							}
						} else {
							log.Warn().Msg(fmt.Sprintf(logWarnMandatoryFieldNotAString, f, dID, class))
							errorsFound++
						}
					}
				}
			}
		}
	}

	// for each definition in the dictionary, ensure that the references are valid
	for _, definitions := range d {
		for _, definition := range definitions {
			for _, ref := range definition.References {
				if d[ref.Class] == nil {
					log.Warn().Msg(fmt.Sprintf(logWarnCannotFindClass, ref.Class))
					errorsFound++
				} else if d[ref.Class][ref.ID] == nil {
					log.Warn().Msg(fmt.Sprintf(logWarnCannotFindDefinition, ref.ID, ref.Class))
					errorsFound++
				}
			}
		}
	}

	if errorsFound > 0 {
		log.Error().Msg(fmt.Sprintf(errorDefinitionErrorsFound, errorsFound))
		return fmt.Errorf(errorDefinitionErrorsFound, errorsFound)
	}

	return nil
}

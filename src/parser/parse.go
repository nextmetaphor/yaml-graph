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
	"github.com/nextmetaphor/yaml-graph/definition"
	"os"
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
		//log.Debug().Msg(fmt.Sprintf(logDebugAboutToParseFile, filePath))

		spec, err := definition.LoadSpecificationFromFile(filePath)
		if (err == nil) && (spec != nil) {
			loadSpecification(*spec, d, nil)

		} else {
			//log.Warn().Msgf(logWarnSkippingFile, filePath, err)
		}

		return nil
	})

	return d
}

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

package definition

import (
	"encoding/base64"
	"fmt"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	definitionFormat                     = ".%s"
	logDebugCannotLoadYAMLFile           = "cannot load YAML file [%s]"
	logDebugCannotParseYAMLFile          = "cannot parse YAML file [%s]"
	logDebugNoDefinitionsFoundInYAMLFile = "no definitions found in YAML file [%s]"
	logErrorCannotProcessFiles           = "cannot process files in root directory [%s]"
	logErrorCannotEncodeFile             = "cannot encode file [%s]"
	logWarnCannotProcessFile             = "cannot process files in directory [%s]"
	logDebugProcessingFile               = "processing file [%s] in directory [%s]"
	logDebugIgnoringFile                 = "ignoring file [%s] in directory [%s]"
)

type (
	// Reference TODO
	Reference struct {
		// Class TODO
		Class string `yaml:"Class"`

		// ID TODO
		ID string `yaml:"ID"`

		// Relationship TODO
		Relationship string `yaml:"Relationship"`

		// RelationshipFrom TODO
		RelationshipFrom bool `yaml:"RelationshipFrom"`

		// RelationshipTo TODO
		RelationshipTo bool `yaml:"RelationshipTo"`
	}

	// Fields TODO
	Fields map[string]interface{}

	// FileDefinition TODO
	FileDefinition struct {
		Path   string
		Prefix string
	}

	// FileFields TODO
	FileFields map[string]FileDefinition

	// Definition TODO
	Definition struct {
		Fields         Fields                   `yaml:"Fields"`
		FileFields     FileFields               `yaml:"FileFields"`
		References     []Reference              `yaml:"References"`
		SubDefinitions map[string]Specification `yaml:"SubDefinitions"`
	}

	// Specification TODO
	Specification struct {
		// Class allows the class for all of the definitions within the document to be specified.
		Class string `yaml:"Class,omitempty"`

		// References allows relationships to other classes for all of the definitions within the document to be
		// specified.
		References []Reference `yaml:"References,omitempty"`

		// Definitions TODO
		Definitions map[string]Definition `yaml:"Definitions,omitempty"`
	}

	processFileFuncType = func(filePath string, fileInfo os.FileInfo) (err error)
)

// simple function to base64 encode the contents of a file and return as a pointer to a string
func getFileFieldAsBase64(fileDefn FileDefinition) (*string, error) {
	dat, err := ioutil.ReadFile(fileDefn.Path)
	if err != nil {
		return nil, err
	}
	encoded := fileDefn.Prefix + base64.StdEncoding.EncodeToString(dat)

	return &encoded, nil
}

func getFileFields(dfn *Definition) {
	for fieldName, fileDefn := range dfn.FileFields {
		b64Str, err := getFileFieldAsBase64(fileDefn)
		if err != nil {
			log.Debug().Err(err).Msg(fmt.Sprintf(logErrorCannotEncodeFile, fileDefn))
		} else {
			dfn.Fields[fieldName] = *b64Str
		}
	}
}

// LoadSpecificationFromFile TODO
func LoadSpecificationFromFile(filename string) (*Specification, error) {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Debug().Err(err).Msg(fmt.Sprintf(logDebugCannotLoadYAMLFile, filename))
		return nil, err
	}

	spec := &Specification{}
	err = yaml.Unmarshal(yamlFile, spec)
	if err != nil {
		log.Debug().Err(err).Msg(fmt.Sprintf(logDebugCannotParseYAMLFile, filename))

		return nil, err
	}

	// if no definitions are found, return an error and a nil Specification
	if len(spec.Definitions) == 0 {
		log.Debug().Err(err).Msg(fmt.Sprintf(logDebugNoDefinitionsFoundInYAMLFile, filename))
		return nil, fmt.Errorf(logDebugNoDefinitionsFoundInYAMLFile, filename)
	}

	// TODO debug
	return spec, nil
}

// ProcessFiles TODO
func ProcessFiles(rootDir, fileExtension string, processFileFunc processFileFuncType) error {
	err := filepath.Walk(rootDir,
		func(filePath string, fileInfo os.FileInfo, err error) error {
			if err != nil {
				log.Warn().Err(err).Msgf(logWarnCannotProcessFile, filePath)
				return err
			}
			if !fileInfo.IsDir() {
				if strings.HasSuffix(fileInfo.Name(), fmt.Sprintf(definitionFormat, fileExtension)) {
					log.Debug().Msg(fmt.Sprintf(logDebugProcessingFile, fileInfo.Name(), filePath))
					return processFileFunc(filePath, fileInfo)
				}

				log.Debug().Msg(fmt.Sprintf(logDebugIgnoringFile, fileInfo.Name(), filePath))
			}
			return nil
		})

	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf(logErrorCannotProcessFiles, rootDir))
		return err
	}

	return nil
}
